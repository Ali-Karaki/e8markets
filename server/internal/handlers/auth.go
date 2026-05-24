package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Ali-Karaki/e8markets/server/internal/config"
	"github.com/Ali-Karaki/e8markets/server/internal/httpx"
	"github.com/Ali-Karaki/e8markets/server/internal/middleware"
	"github.com/Ali-Karaki/e8markets/server/internal/store"
	"github.com/Ali-Karaki/e8markets/server/internal/tradelocker"
	"github.com/google/uuid"
)

type AuthHandler struct {
	cfg         config.Config
	sessions    *store.SessionStore
	tradelocker *tradelocker.Client
	logs        *store.APILogStore
	auth        *middleware.Auth
}

func NewAuthHandler(cfg config.Config, sessions *store.SessionStore, tl *tradelocker.Client, logs *store.APILogStore, auth *middleware.Auth) *AuthHandler {
	return &AuthHandler{
		cfg:         cfg,
		sessions:    sessions,
		tradelocker: tl,
		logs:        logs,
		auth:        auth,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Server   string `json:"server"`
}

type sessionResponse struct {
	Authenticated bool   `json:"authenticated"`
	Email         string `json:"email"`
	Server        string `json:"server"`
	ExpiresAt     string `json:"expiresAt"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Server = strings.TrimSpace(req.Server)
	if req.Server == "" {
		req.Server = h.cfg.TradeLockerServer
	}

	if req.Email == "" || req.Password == "" {
		httpx.Error(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	tokens, err := h.tradelocker.Login(r.Context(), nil, req.Email, req.Password, req.Server)
	if err != nil {
		h.logs.Log(r.Context(), nil, store.EventLoginAttempt, "POST", "/api/auth/login", http.StatusUnauthorized, err.Error())
		httpx.Error(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	session := store.Session{
		Email:        req.Email,
		Server:       req.Server,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
	}

	if err := h.sessions.Create(r.Context(), session); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	h.logs.Log(r.Context(), &session.ID, store.EventLoginAttempt, "POST", "/api/auth/login", http.StatusOK, "")
	h.auth.SetSessionCookie(w, session)

	httpx.JSON(w, http.StatusOK, sessionResponse{
		Authenticated: true,
		Email:         session.Email,
		Server:        session.Server,
		ExpiresAt:     session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	session, ok := h.auth.Optional(r)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	httpx.JSON(w, http.StatusOK, sessionResponse{
		Authenticated: true,
		Email:         session.Email,
		Server:        session.Server,
		ExpiresAt:     session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(h.cfg.SessionCookieName); err == nil {
		if id, err := uuid.Parse(cookie.Value); err == nil {
			_ = h.sessions.Delete(r.Context(), id)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}
