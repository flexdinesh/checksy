package check

import (
	"context"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// ICMPAttempt tries to measure ICMP echo RTT to a host. It returns ok=false
// when ICMP is unavailable (e.g. unprivileged on macOS), signalling Ping to
// fall back to TCP.
type ICMPAttempt func(ctx context.Context, host string) (Result, bool)

// TCPDial measures a TCP-connect RTT. It is injectable so Ping's fallback path
// can be exercised against local fixtures without touching the public internet.
type TCPDial func(ctx context.Context, network, address string) (net.Conn, error)

// Ping measures reachability of a host: it tries ICMP echo and on any failure
// falls back to a TCP-connect RTT to port 443. Result.Method records the path.
// See docs/adr/0001-icmp-with-tcp-fallback.md.
func Ping(ctx context.Context, host string, attempt ICMPAttempt, dial TCPDial, timeout time.Duration) Result {
	if attempt != nil {
		if result, ok := attempt(ctx, host); ok {
			return result
		}
	}
	if dial == nil {
		dial = realTCPDial
	}
	return pingTCP(ctx, host, 443, dial, timeout)
}

func pingTCP(ctx context.Context, host string, port int, dial TCPDial, timeout time.Duration) Result {
	start := time.Now()
	conn, err := dial(ctx, "tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return failPing(host, err, timeout)
	}
	conn.Close()
	return Result{
		Kind:    KindPing,
		Label:   host,
		Status:  StatusOK,
		Latency: time.Since(start),
		Method:  "tcp",
	}
}

func realTCPDial(ctx context.Context, network, address string) (net.Conn, error) {
	dialer := net.Dialer{}
	return dialer.DialContext(ctx, network, address)
}

func failPing(host string, err error, timeout time.Duration) Result {
	return Result{Kind: KindPing, Label: host, Status: StatusFail, Method: "tcp", Detail: shortErr(err, timeout), Err: err}
}

// realICMP attempts an ICMP echo via raw sockets. On macOS this needs root and
// fails unprivileged; Ping then falls back to TCP. Returns ok=false on any
// error so the caller never has to special-case the permission failure.
func realICMP(ctx context.Context, host string) (Result, bool) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return Result{}, false
	}
	defer conn.Close()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(time.Second)
	}
	if err := conn.SetDeadline(deadline); err != nil {
		return Result{}, false
	}

	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{ID: os.Getpid() & 0xffff, Seq: 1, Data: []byte("checksy")},
	}
	packet, err := message.Marshal(nil)
	if err != nil {
		return Result{}, false
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return Result{}, false
	}

	start := time.Now()
	if _, err := conn.WriteTo(packet, &net.IPAddr{IP: ip}); err != nil {
		return Result{}, false
	}
	buffer := make([]byte, 1500)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		return Result{}, false
	}
	parsed, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), buffer[:n])
	if err != nil {
		return Result{}, false
	}
	if parsed.Type != ipv4.ICMPTypeEchoReply {
		return Result{}, false
	}
	return Result{
		Kind:    KindPing,
		Label:   host,
		Status:  StatusOK,
		Latency: time.Since(start),
		Method:  "icmp",
	}, true
}
