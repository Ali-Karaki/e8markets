package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"

	"github.com/Ali-Karaki/e8markets/server/internal/apperr"
	"github.com/Ali-Karaki/e8markets/server/internal/config"
	"github.com/Ali-Karaki/e8markets/server/internal/handlers"
	"github.com/Ali-Karaki/e8markets/server/internal/middleware"
	"github.com/Ali-Karaki/e8markets/server/internal/server"
	"github.com/Ali-Karaki/e8markets/server/internal/store"
	"github.com/Ali-Karaki/e8markets/server/internal/tradelocker"
	"github.com/redis/go-redis/v9"
)

type noopAPILogger struct{}

func (noopAPILogger) Log(_ context.Context, _ *uuid.UUID, _, _, _ string, _ int, _ string) {}

type apiErrorBody struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func TestLogin_success(t *testing.T) {
	ts := newTestServer(t, mockTLHandler(t, mockTLConfig{loginStatus: http.StatusOK}))

	body := `{"email":"user@example.com","password":"secret","server":"E8"}`
	req := httptest.NewRequest(http.MethodPost, ts.appURL+"/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ts.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}

	var session map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &session); err != nil {
		t.Fatal(err)
	}
	if session["authenticated"] != true {
		t.Fatalf("authenticated = %v", session["authenticated"])
	}

	cookies := rec.Result().Cookies()
	if len(cookies) == 0 || cookies[0].Name != "session_id" || cookies[0].Value == "" {
		t.Fatalf("cookies = %+v", cookies)
	}
}

func TestLogin_invalidCredentials(t *testing.T) {
	ts := newTestServer(t, mockTLHandler(t, mockTLConfig{loginStatus: http.StatusUnauthorized}))

	body := `{"email":"user@example.com","password":"wrong","server":"E8"}`
	req := httptest.NewRequest(http.MethodPost, ts.appURL+"/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ts.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rec.Code)
	}

	errBody := decodeAPIError(t, rec.Body.Bytes())
	if errBody.Code != apperr.CodeInvalidCredentials {
		t.Fatalf("code = %q", errBody.Code)
	}
}

func TestInstruments_unauthorized(t *testing.T) {
	ts := newTestServer(t, mockTLHandler(t, mockTLConfig{}))

	req := httptest.NewRequest(http.MethodGet, ts.appURL+"/api/instruments?accountId=1&accNum=1", nil)
	rec := httptest.NewRecorder()
	ts.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rec.Code)
	}

	errBody := decodeAPIError(t, rec.Body.Bytes())
	if errBody.Code != apperr.CodeNotAuthenticated {
		t.Fatalf("code = %q", errBody.Code)
	}
}

func TestAccounts_success(t *testing.T) {
	ts := newTestServer(t, mockTLHandler(t, mockTLConfig{loginStatus: http.StatusOK}))

	cookie := loginAndGetCookie(t, ts)

	req := httptest.NewRequest(http.MethodGet, ts.appURL+"/api/accounts", nil)
	req.AddCookie(cookie)

	rec := httptest.NewRecorder()
	ts.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Accounts []map[string]any `json:"accounts"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Accounts) != 1 {
		t.Fatalf("accounts = %+v", resp.Accounts)
	}
}

func TestInstruments_success(t *testing.T) {
	ts := newTestServer(t, mockTLHandler(t, mockTLConfig{loginStatus: http.StatusOK}))

	cookie := loginAndGetCookie(t, ts)

	req := httptest.NewRequest(http.MethodGet, ts.appURL+"/api/instruments?accountId=2200766&accNum=1", nil)
	req.AddCookie(cookie)

	rec := httptest.NewRecorder()
	ts.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d body = %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Instruments []map[string]any `json:"instruments"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Instruments) != 1 || resp.Instruments[0]["name"] != "EURUSD" {
		t.Fatalf("instruments = %+v", resp.Instruments)
	}
}

type testServer struct {
	handler http.Handler
	appURL  string
}

type mockTLConfig struct {
	loginStatus int
}

func newTestServer(t *testing.T, tlHandler http.Handler) *testServer {
	t.Helper()

	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	tlServer := httptest.NewServer(tlHandler)
	t.Cleanup(tlServer.Close)

	cfg := config.Config{
		Port:               "8080",
		TradeLockerBaseURL: tlServer.URL,
		TradeLockerServer:  "E8",
		CORSOrigin:         "http://localhost:5173",
		SessionCookieName:  "session_id",
	}

	sessions := store.NewSessionStore(rdb)
	var logs noopAPILogger
	tl := tradelocker.NewClient(cfg.TradeLockerBaseURL, logs)
	authMW := middleware.NewAuth(cfg, sessions, tl)

	authHandler := handlers.NewAuthHandler(cfg, sessions, tl, logs, authMW)
	accountsHandler := handlers.NewAccountsHandler(tl)
	instrumentsHandler := handlers.NewInstrumentsHandler(tl)
	positionsHandler := handlers.NewPositionsHandler(tl, store.NewPositionStore(nil))

	handler := server.NewHandler(server.Deps{
		Cfg:         cfg,
		Auth:        authHandler,
		Accounts:    accountsHandler,
		Instruments: instrumentsHandler,
		Positions:   positionsHandler,
		AuthMW:      authMW,
	})

	return &testServer{
		handler: handler,
		appURL:  "http://localhost:5173",
	}
}

func mockTLHandler(t *testing.T, cfg mockTLConfig) http.Handler {
	t.Helper()

	if cfg.loginStatus == 0 {
		cfg.loginStatus = http.StatusOK
	}

	tokenBody := `{"accessToken":"access-token","refreshToken":"refresh-token","expireDate":"` +
		time.Now().Add(time.Hour).UTC().Format(time.RFC3339) + `"}`

	accountsBody := `{"accounts":[{"id":"2200766","name":"E8#demo","currency":"USD","accNum":"1","accountBalance":"100000.00","status":"ACTIVE"}]}`

	instrumentsBody := `{"d":{"instruments":[{"id":1,"name":"EURUSD","tradableInstrumentId":1001,"type":"FOREX","tradingExchange":"FX","marketDataExchange":"FX","description":"Euro vs US Dollar","routes":[{"id":42,"type":"TRADE"},{"id":43,"type":"INFO"}]}]},"s":"ok"}`

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/auth/jwt/token":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(cfg.loginStatus)
			if cfg.loginStatus == http.StatusOK {
				_, _ = w.Write([]byte(tokenBody))
			} else {
				_, _ = w.Write([]byte(`{"error":"invalid"}`))
			}
		case r.Method == http.MethodGet && r.URL.Path == "/auth/jwt/all-accounts":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(accountsBody))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/trade/accounts/") && strings.HasSuffix(r.URL.Path, "/instruments"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(instrumentsBody))
		default:
			t.Fatalf("unexpected TL request: %s %s", r.Method, r.URL.Path)
		}
	})
}

func loginAndGetCookie(t *testing.T, ts *testServer) *http.Cookie {
	t.Helper()

	body := `{"email":"user@example.com","password":"secret","server":"E8"}`
	req := httptest.NewRequest(http.MethodPost, ts.appURL+"/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ts.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("login status = %d body = %s", rec.Code, rec.Body.String())
	}

	for _, c := range rec.Result().Cookies() {
		if c.Name == "session_id" {
			return c
		}
	}
	t.Fatal("session cookie not set")
	return nil
}

func decodeAPIError(t *testing.T, body []byte) apiErrorBody {
	t.Helper()

	var errBody apiErrorBody
	if err := json.Unmarshal(body, &errBody); err != nil {
		t.Fatalf("decode error body: %v raw=%s", err, string(body))
	}
	return errBody
}
