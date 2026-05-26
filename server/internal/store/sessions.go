package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const sessionKeyPrefix = "session:"

type Session struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Server       string    `json:"server"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type SessionStore struct {
	client *redis.Client
}

func NewSessionStore(client *redis.Client) *SessionStore {
	return &SessionStore{client: client}
}

func (s *SessionStore) Create(ctx context.Context, session Session) (Session, error) {
	session.ID = uuid.New()
	return session, s.save(ctx, session)
}

func (s *SessionStore) Get(ctx context.Context, id uuid.UUID) (Session, error) {
	data, err := s.client.Get(ctx, sessionKey(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return Session{}, ErrSessionNotFound
		}
		return Session{}, fmt.Errorf("get session: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return Session{}, fmt.Errorf("decode session: %w", err)
	}

	return session, nil
}

func (s *SessionStore) Update(ctx context.Context, session Session) error {
	return s.save(ctx, session)
}

func (s *SessionStore) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.client.Del(ctx, sessionKey(id)).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (s *SessionStore) save(ctx context.Context, session Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("encode session: %w", err)
	}

	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Hour
	}

	if err := s.client.Set(ctx, sessionKey(session.ID), data, ttl).Err(); err != nil {
		return fmt.Errorf("save session: %w", err)
	}

	return nil
}

func sessionKey(id uuid.UUID) string {
	return sessionKeyPrefix + id.String()
}
