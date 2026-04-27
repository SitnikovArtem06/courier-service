package order

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetOrder_RetryOn429_ThenSuccess(t *testing.T) {
	var calls int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)

		if !strings.HasPrefix(r.URL.Path, "/public/api/v1/order/") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if atomic.LoadInt32(&calls) <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	gw := NewHttpGateway(srv.URL, &http.Client{Timeout: 2 * time.Second})

	start := time.Now()
	_, err := gw.GetOrder(context.Background(), "123")
	end := time.Since(start)

	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}

	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Fatalf("expected 3 calls (429,429,200), got %d", got)
	}

	minExpected := 2*RetryDelay - 20*time.Millisecond
	if end < minExpected {
		t.Fatalf("expected elapsed >= %v, got %v", minExpected, end)
	}
}

func TestGetOrder_RetryOn429_Exhausted(t *testing.T) {
	var calls int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	gw := NewHttpGateway(srv.URL, &http.Client{Timeout: 2 * time.Second})

	start := time.Now()
	_, err := gw.GetOrder(context.Background(), "123")
	end := time.Since(start)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if got := atomic.LoadInt32(&calls); got != int32(RetryAttempts) {
		t.Fatalf("expected %d calls, got %d", RetryAttempts, got)
	}

	minExpected := time.Duration(RetryAttempts-1)*RetryDelay - 30*time.Millisecond
	if end < minExpected {
		t.Fatalf("expected elapsed >= %v, got %v", minExpected, end)
	}
}
