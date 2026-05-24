package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Ali-Karaki/e8markets/server/internal/config"
	"github.com/Ali-Karaki/e8markets/server/internal/httpx"
	"github.com/Ali-Karaki/e8markets/server/internal/store"
	"github.com/Ali-Karaki/e8markets/server/internal/tradelocker"
	"github.com/google/uuid"
)

type contextKey string

const SessionContextKey contextKey = "session"

type Auth struct {
	cfg           config.Config
	sessions      *store.SessionStore
	tradelocker   *tradelocker.Client
	refreshBuffer time.Duration
}

func NewAuth(cfg config.Config, sessions *store.SessionStore, tl *tradelocker.Client) *Auth {
	return &Auth{
		cfg:           cfg,
		sessions:      sessions,
		tradelocker:   tl,
		refreshBuffer: 60 * time.Second,
	}
}

func (a *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, ok := a.loadSession(w, r)
		if !ok {
			return
		}

		ctx := context.WithValue(r.Context(), SessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) Optional(r *http.Request) (store.Session, bool) {
	cookie, err := r.Cookie(a.cfg.SessionCookieName)
	if err != nil {
		return store.Session{}, false
	}

	id, err := uuid.Parse(cookie.Value)
	if err != nil {
		return store.Session{}, false
	}

	session, err := a.sessions.Get(r.Context(), id)
	if err != nil {
		return store.Session{}, false
	}

	session, err = a.ensureFreshTokens(r.Context(), session)
	if err != nil {
		return store.Session{}, false
	}

	return session, true
}

func (a *Auth) loadSession(w http.ResponseWriter, r *http.Request) (store.Session, bool) {
	cookie, err := r.Cookie(a.cfg.SessionCookieName)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "Not authenticated")
		return store.Session{}, false
	}

	id, err := uuid.Parse(cookie.Value)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "Invalid session")
		return store.Session{}, false
	}

	session, err := a.sessions.Get(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "Session expired")
		return store.Session{}, false
	}

	session, err = a.ensureFreshTokens(r.Context(), session)
	if err != nil {
		_ = a.sessions.Delete(r.Context(), session.ID)
		clearSessionCookie(w, a.cfg)
		httpx.Error(w, http.StatusUnauthorized, "Session expired")
		return store.Session{}, false
	}

	return session, true
}

func (a *Auth) ensureFreshTokens(ctx context.Context, session store.Session) (store.Session, error) {
	if time.Until(session.ExpiresAt) > a.refreshBuffer {
		return session, nil
	}

	sid := session.ID.String()
	tokens, err := a.tradelocker.Refresh(ctx, &sid, session.AccessToken, session.RefreshToken)
	if err != nil {
		return store.Session{}, err
	}

	session.AccessToken = tokens.AccessToken
	session.RefreshToken = tokens.RefreshToken
	session.ExpiresAt = tokens.ExpiresAt

	if err := a.sessions.Update(ctx, session); err != nil {
		return store.Session{}, err
	}

	return session, nil
}

func (a *Auth) SetSessionCookie(w http.ResponseWriter, session store.Session) {
	maxAge := int(time.Until(session.ExpiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     a.cfg.SessionCookieName,
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})
}

func clearSessionCookie(w http.ResponseWriter, cfg config.Config) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func SessionFromContext(ctx context.Context) (store.Session, bool) {
	session, ok := ctx.Value(SessionContextKey).(store.Session)
	return session, ok
}
