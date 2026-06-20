package check

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestPingFallsBackToTCPSuccessAgainstLocalListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	attempt := func(context.Context, string) (Result, bool) { return Result{}, false }
	dial := func(ctx context.Context, network, _ string) (net.Conn, error) {
		dialer := net.Dialer{}
		return dialer.DialContext(ctx, network, listener.Addr().String())
	}

	result := Ping(context.Background(), "127.0.0.1", attempt, dial, 5*time.Second)

	if result.Kind != KindPing {
		t.Fatalf("expected kind ping, got %s", result.Kind)
	}
	if result.Status != StatusOK {
		t.Fatalf("expected ok, got %d (detail %q)", result.Status, result.Detail)
	}
	if result.Method != "tcp" {
		t.Fatalf("expected method tcp, got %s", result.Method)
	}
	if result.Latency <= 0 {
		t.Fatalf("expected positive latency, got %v", result.Latency)
	}
}

func TestPingTCPFailsOnClosedPort(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	result := pingTCP(ctx, "127.0.0.1", 1, realTCPDial, 200*time.Millisecond)

	if result.Status != StatusFail {
		t.Fatalf("expected fail on closed port, got %d", result.Status)
	}
}

func TestPingUsesICMPWhenAttemptSucceeds(t *testing.T) {
	attempt := func(context.Context, string) (Result, bool) {
		return Result{Kind: KindPing, Status: StatusOK, Method: "icmp", Latency: 1 * time.Millisecond}, true
	}

	result := Ping(context.Background(), "1.1.1.1", attempt, nil, 5*time.Second)

	if result.Method != "icmp" {
		t.Fatalf("expected method icmp, got %s", result.Method)
	}
	if result.Status != StatusOK {
		t.Fatalf("expected ok, got %d", result.Status)
	}
}

func TestPingFallsBackToTCPWhenICMPUnavailable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	attempt := func(context.Context, string) (Result, bool) { return Result{}, false }

	result := Ping(ctx, "127.0.0.1", attempt, realTCPDial, 200*time.Millisecond)

	if result.Method != "tcp" {
		t.Fatalf("expected fallback to tcp, got %s", result.Method)
	}
}
