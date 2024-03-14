package main

import (
	"time"

	"github.com/lukeshay/g/auth"
	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	ID           string    `bun:",pk"`
	UserID       string    `bun:",notnull"`
	ExpiresAt    time.Time `bun:",notnull"`
	RefreshUntil time.Time `bun:",notnull"`
}

var _ auth.Session = (*Session)(nil)

func (s *Session) GetSessionID() string {
	return s.ID
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

func (s *Session) SetSessionID(id string) {
	s.ID = id
}

func (s *Session) Copy() auth.Session {
	return &Session{
		ID:           s.ID,
		UserID:       s.UserID,
		ExpiresAt:    s.ExpiresAt,
		RefreshUntil: s.RefreshUntil,
	}
}

func (s *Session) SetExpiresAt(expiresAt time.Time) {
	s.ExpiresAt = expiresAt
}
