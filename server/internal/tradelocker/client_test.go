package tradelocker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

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

func TestDoJSON_logsDecodeError(t *testing.T) {
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
