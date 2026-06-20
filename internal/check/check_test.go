package check

import "testing"

func TestVerdict(t *testing.T) {
	tests := []struct {
		name    string
		results []Result
		want    Status
	}{
		{
			name:    "empty is down",
			results: nil,
			want:    StatusFail,
		},
		{
			name: "http ok is up",
			results: []Result{
				{Kind: KindHTTP, Status: StatusOK},
			},
			want: StatusOK,
		},
		{
			name: "http fail is down",
			results: []Result{
				{Kind: KindHTTP, Status: StatusFail},
			},
			want: StatusFail,
		},
		{
			name: "ping ok without http is down",
			results: []Result{
				{Kind: KindPing, Status: StatusOK},
			},
			want: StatusFail,
		},
		{
			name: "ping fail with http ok is still up",
			results: []Result{
				{Kind: KindPing, Status: StatusFail},
				{Kind: KindHTTP, Status: StatusOK},
			},
			want: StatusOK,
		},
		{
			name: "any http ok wins over other http fail",
			results: []Result{
				{Kind: KindHTTP, Status: StatusFail},
				{Kind: KindHTTP, Status: StatusOK},
			},
			want: StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Verdict(test.results); got != test.want {
				t.Fatalf("expected %d, got %d", test.want, got)
			}
		})
	}
}
