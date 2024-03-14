package adaptors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lukeshay/g/auth"
)

type Session struct {
	SessionID    string    `json:"-" xml:"-" yaml:"-"`
	UserID       string    `json:"-" xml:"-" yaml:"-"`
	ExpiresAt    time.Time `json:"-" xml:"-" yaml:"-"`
	RefreshUntil time.Time `json:"-" xml:"-" yaml:"-"`
}

var _ auth.Session = &Session{}

func (s *Session) GetSessionID() string {
	return s.SessionID
}

func (s *Session) GetUserID() string {
	return s.UserID
}

func (s *Session) GetExpiresAt() time.Time {
	return s.ExpiresAt
}

func (s *Session) GetRefreshUntil() time.Time {
	return s.RefreshUntil
}

func (s *Session) SetSessionID(sessionID string) {
	s.SessionID = sessionID
}

func (s *Session) SetExpiresAt(expiresAt time.Time) {
	s.ExpiresAt = expiresAt
}

func (s *Session) Copy() auth.Session {
	return &Session{
		SessionID:    s.SessionID,
		UserID:       s.UserID,
		ExpiresAt:    s.ExpiresAt,
		RefreshUntil: s.RefreshUntil,
	}
}

type InMemoryAdapter struct {
	sessions sync.Map
}

func NewInMemoryAdapter() auth.SessionAdapter {
	return &InMemoryAdapter{
		sessions: sync.Map{},
	}
}

func (a *InMemoryAdapter) GetSession(ctx context.Context, sessionID string) (auth.Session, error) {
	value, found := a.sessions.Load(sessionID)
	if !found {
		return nil, fmt.Errorf("session not found")
	}

	return value.(*Session), nil
}

func (a *InMemoryAdapter) InsertSession(ctx context.Context, newSession auth.Session) error {
	session := newSession.(*Session)

	a.sessions.Store(session.SessionID, session)

	return nil
}

func (a *InMemoryAdapter) UpdateSession(ctx context.Context, newSession auth.Session) error {
	return a.InsertSession(ctx, newSession)
}

func (a *InMemoryAdapter) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	a.sessions.Range(func(key any, value any) bool {
		session := value.(*Session)
		if session.UserID == userID {
			a.sessions.Delete(key)
		}

		return true
	})

	return nil
}

func (a *InMemoryAdapter) DeleteSession(ctx context.Context, sessionID string) error {
	a.sessions.Delete(sessionID)

	return nil
}
