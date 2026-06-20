package check

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	result := HTTP(context.Background(), "test", server.URL, http.StatusNoContent, 5*time.Second)

	if result.Status != StatusOK {
		t.Fatalf("expected ok, got %d (detail %q)", result.Status, result.Detail)
	}
	if result.Detail != "204" {
		t.Fatalf("expected detail 204, got %q", result.Detail)
	}
	if result.Method != "http" {
		t.Fatalf("expected method http, got %q", result.Method)
	}
	if result.Latency <= 0 {
		t.Fatalf("expected positive latency, got %v", result.Latency)
	}
}

func TestHTTPWrongStatusIsFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := HTTP(context.Background(), "test", server.URL, http.StatusNoContent, 5*time.Second)

	if result.Status != StatusFail {
		t.Fatalf("expected fail for wrong status, got %d", result.Status)
	}
	if result.Detail != "200" {
		t.Fatalf("expected detail 200, got %q", result.Detail)
	}
	if result.Method != "http" {
		t.Fatalf("expected method http, got %q", result.Method)
	}
}

func TestHTTPUnreachableIsFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := HTTP(ctx, "test", "http://example.com", http.StatusNoContent, 5*time.Second)

	if result.Status != StatusFail {
		t.Fatalf("expected fail when unreachable, got %d", result.Status)
	}
}

func TestHTTPTimeoutIsFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w := r.Context().Done()
		<-w
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result := HTTP(ctx, "test", server.URL, http.StatusNoContent, 50*time.Millisecond)

	if result.Status != StatusFail {
		t.Fatalf("expected fail on timeout, got %d", result.Status)
	}
}

func TestHTTPTimeoutDetailShowsDuration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result := HTTP(ctx, "test", server.URL, http.StatusNoContent, 50*time.Millisecond)

	if result.Status != StatusFail {
		t.Fatalf("expected fail on timeout, got %d", result.Status)
	}
	if result.Detail != "timeout after 50ms" {
		t.Fatalf("expected detail %q, got %q", "timeout after 50ms", result.Detail)
	}
}
