package server

import (
	"net/http"

	"github.com/Ali-Karaki/e8markets/server/internal/config"
	"github.com/Ali-Karaki/e8markets/server/internal/handlers"
	"github.com/Ali-Karaki/e8markets/server/internal/middleware"
)

type Deps struct {
	Cfg      config.Config
	Auth     *handlers.AuthHandler
	Accounts *handlers.AccountsHandler
	AuthMW   *middleware.Auth
}

func NewHandler(deps Deps) http.Handler {
	mux := http.NewServeMux()
	protected := func(h http.HandlerFunc) http.Handler {
		return deps.AuthMW.Require(h)
	}

	mux.HandleFunc("GET /health", handlers.Health)

	mux.HandleFunc("POST /api/auth/login", deps.Auth.Login)
	mux.HandleFunc("GET /api/auth/session", deps.Auth.Session)
	mux.HandleFunc("POST /api/auth/logout", deps.Auth.Logout)

	mux.Handle("GET /api/accounts", protected(deps.Accounts.List))
	mux.Handle("GET /api/accounts/state", protected(deps.Accounts.State))

	return middleware.CORS(deps.Cfg, middleware.RequestLogger(mux))
}
