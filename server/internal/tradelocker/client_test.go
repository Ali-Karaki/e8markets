package tradelocker

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ali-Karaki/e8markets/server/internal/apperr"
	"github.com/google/uuid"
)

func TestDoJSON_classifiesUpstreamErrors(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		wantCode string
	}{
		{
			name:     "rate limited",
			status:   http.StatusTooManyRequests,
			body:     `{"error":"slow down"}`,
			wantCode: apperr.CodeUpstreamRateLimited,
		},
		{
			name:     "upstream error",
			status:   http.StatusBadGateway,
			body:     `{"error":"bad gateway"}`,
			wantCode: apperr.CodeUpstreamError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewClient(server.URL, nil)
			_, err := client.doJSON(context.Background(), nil, http.MethodGet, "/test", "", nil, "", nil)
			if err == nil {
				t.Fatal("expected error")
			}

			var tlErr *Error
			if !errors.As(err, &tlErr) {
				t.Fatalf("expected *tradelocker.Error, got %T", err)
			}
			if tlErr.Code != tt.wantCode {
				t.Fatalf("code = %q, want %q", tlErr.Code, tt.wantCode)
			}
		})
	}
}

func TestDoJSON_logsUpstreamErrorOn4xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	logger := &mockLogger{}
	client := NewClient(server.URL, logger)

	_, err := client.doJSON(context.Background(), nil, http.MethodGet, "/test", "", nil, "", nil)
	if err == nil {
		t.Fatal("expected error")
	}

	if len(logger.calls) != 1 {
		t.Fatalf("log calls = %d", len(logger.calls))
	}

	call := logger.calls[0]
	if call.statusCode != http.StatusBadRequest {
		t.Fatalf("status = %d", call.statusCode)
	}
	if call.message != `{"error":"bad request"}` {
		t.Fatalf("message = %q", call.message)
	}
}

func TestDoJSON_classifiesMalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"accounts":[{"accountBalance":"not-a-number"}]}`))
	}))
	defer server.Close()

	logger := &mockLogger{}
	client := NewClient(server.URL, logger)

	var resp AccountsResponse
	_, err := client.doJSON(context.Background(), nil, http.MethodGet, "/auth/jwt/all-accounts", "", nil, "", &resp)
	if err == nil {
		t.Fatal("expected decode error")
	}

	var tlErr *Error
	if !errors.As(err, &tlErr) {
		t.Fatalf("expected *tradelocker.Error, got %T", err)
	}
	if tlErr.Code != apperr.CodeUpstreamMalformed {
		t.Fatalf("code = %q", tlErr.Code)
	}

	if len(logger.calls) != 1 {
		t.Fatalf("log calls = %d", len(logger.calls))
	}
	call := logger.calls[0]
	if call.statusCode != 200 {
		t.Fatalf("status = %d", call.statusCode)
	}
	if len(call.message) < 7 || call.message[:7] != "decode:" {
		t.Fatalf("message = %q", call.message)
	}
}

type mockLogger struct {
	calls []logCall
}

type logCall struct {
	path       string
	statusCode int
	message    string
}

func (m *mockLogger) Log(_ context.Context, _ *uuid.UUID, _, method, path string, statusCode int, message string) {
	m.calls = append(m.calls, logCall{path: method + " " + path, statusCode: statusCode, message: message})
}
