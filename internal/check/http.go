package check

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"
)

var httpClient = &http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// Verdict target: Google's captive-check endpoint. A free internet returns 204;
// a captive portal intercepts it and returns a different status.
const (
	VerdictURL   = "https://connectivitycheck.gstatic.com/generate_204"
	VerdictLabel = "gstatic.com"
)

// HTTP runs a single HTTP check: GET url expecting wantStatus. The verdict
// (exit code) is driven only by whether an http-kind Result is OK.
func HTTP(ctx context.Context, label, url string, wantStatus int, timeout time.Duration) Result {
	start := time.Now()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return failHTTP(label, err, timeout)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return failHTTP(label, err, timeout)
	}
	defer response.Body.Close()

	latency := time.Since(start)
	status := StatusFail
	if response.StatusCode == wantStatus {
		status = StatusOK
	}
	return Result{
		Kind:    KindHTTP,
		Label:   label,
		Status:  status,
		Latency: latency,
		Method:  "http",
		Detail:  strconv.Itoa(response.StatusCode),
	}
}

func failHTTP(label string, err error, timeout time.Duration) Result {
	return Result{
		Kind:   KindHTTP,
		Label:  label,
		Status: StatusFail,
		Method: "http",
		Detail: shortErr(err, timeout),
		Err:    err,
	}
}

func shortErr(err error, timeout time.Duration) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout after " + timeout.String()
	}
	return err.Error()
}
