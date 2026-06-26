package check

import (
	"context"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

const defaultRouteTarget = "1.1.1.1:443"

func defaultLocalIP(ctx context.Context, target string) string {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "udp", target)
	if err != nil {
		return ""
	}
	defer conn.Close()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok || addr.IP == nil {
		return ""
	}
	return addr.IP.String()
}

func defaultGateway(ctx context.Context) string {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.CommandContext(ctx, "route", "-n", "get", "default")
	case "linux":
		cmd = exec.CommandContext(ctx, "ip", "route", "show", "default")
	default:
		return ""
	}

	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return parseGateway(runtime.GOOS, string(out))
}

func parseGateway(goos, body string) string {
	switch goos {
	case "darwin":
		for _, line := range strings.Split(body, "\n") {
			line = strings.TrimSpace(line)
			gateway, ok := strings.CutPrefix(line, "gateway:")
			if ok {
				return cleanIP(gateway)
			}
		}
	case "linux":
		fields := strings.Fields(body)
		for index, field := range fields {
			if field == "via" && index+1 < len(fields) {
				return cleanIP(fields[index+1])
			}
		}
	}
	return ""
}

func cleanIP(value string) string {
	ip := net.ParseIP(strings.TrimSpace(value))
	if ip == nil {
		return ""
	}
	return ip.String()
}
