package check

import "testing"

func TestParseGatewayReadsDefaultRouteOutput(t *testing.T) {
	tests := []struct {
		name string
		goos string
		body string
		want string
	}{
		{name: "darwin", goos: "darwin", body: "route to: default\ndestination: default\ngateway: 192.168.1.1\ninterface: en0\n", want: "192.168.1.1"},
		{name: "linux", goos: "linux", body: "default via 172.20.0.1 dev eth0 proto dhcp src 172.20.0.10 metric 100\n", want: "172.20.0.1"},
		{name: "no gateway", goos: "linux", body: "default dev eth0 scope link\n", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := parseGateway(test.goos, test.body); got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}
