package store

import (
	"context"

	"github.com/google/uuid"
)

const (
	EventLoginAttempt       = "login_attempt"
	EventTradeLockerRequest = "tradelocker_request"
)

type APILogStore struct {
	pg *Postgres
}

func NewAPILogStore(pg *Postgres) *APILogStore {
	return &APILogStore{pg: pg}
}

func (l *APILogStore) Log(ctx context.Context, sessionID *uuid.UUID, eventType, method, path string, statusCode int, message string) {
	l.insert(ctx, sessionID, eventType, method, path, statusCode, message)
}

func (l *APILogStore) insert(ctx context.Context, sessionID *uuid.UUID, eventType, method, path string, statusCode int, message string) {
	_, _ = l.pg.pool.Exec(ctx, `
		INSERT INTO api_logs (session_id, event_type, method, path, status_code, message)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, sessionID, eventType, method, path, statusCode, message)
}
