package check

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDNSRecordsResolvedIP(t *testing.T) {
	lookup := func(context.Context, string) ([]string, error) {
		return []string{"1.1.1.1", "1.0.0.1"}, nil
	}

	result := DNS(context.Background(), "one.one.one.one", lookup, 5*time.Second)

	if result.Kind != KindDNS {
		t.Fatalf("expected kind dns, got %s", result.Kind)
	}
	if result.Status != StatusOK {
		t.Fatalf("expected ok, got %d", result.Status)
	}
	if result.Detail != "1.1.1.1" {
		t.Fatalf("expected detail 1.1.1.1, got %q", result.Detail)
	}
	if result.Method != "system" {
		t.Fatalf("expected method system, got %q", result.Method)
	}
	if result.Latency <= 0 {
		t.Fatalf("expected positive latency, got %v", result.Latency)
	}
}

func TestDNSFailsOnLookupError(t *testing.T) {
	lookup := func(context.Context, string) ([]string, error) {
		return nil, errors.New("no such host")
	}

	result := DNS(context.Background(), "one.one.one.one", lookup, 5*time.Second)

	if result.Status != StatusFail {
		t.Fatalf("expected fail, got %d", result.Status)
	}
}

func TestDNSRespectsContextTimeout(t *testing.T) {
	lookup := func(ctx context.Context, _ string) ([]string, error) {
		select {
		case <-time.After(200 * time.Millisecond):
			return []string{"1.1.1.1"}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	result := DNS(ctx, "one.one.one.one", lookup, 20*time.Millisecond)

	if result.Status != StatusFail {
		t.Fatalf("expected fail on timeout, got %d", result.Status)
	}
	if result.Detail != "timeout after 20ms" {
		t.Fatalf("expected detail %q, got %q", "timeout after 20ms", result.Detail)
	}
}
