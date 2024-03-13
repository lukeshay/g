package auth

import (
	"time"
)

// Session represents a user session. It contains the user ID, the session ID,
// the expiration time, and any other attributes that are associated with the
// session. The attributes can be used to store additional information about
// the session, such as the user's IP address, the user agent, and so on. The
// session ID is used to identify the session, and the user ID is used to link
// the session to the user. The expiration time is used to determine if the
// session is still valid. If the session is expired, it is considered invalid.
type Session struct {
	ID           string
	UserID       string
	ExpiresAt    time.Time
	RefreshUntil time.Time
	Attributes   map[string]interface{}
}

// SessionV2 is an interface that represents a user session. It contains the
// functions required for this package to interact with the session. The
// SessionV2 interface is used to allow for different implementations of the
// session to be used.
type SessionV2 interface {
	GetSessionID() string
	SetSessionID(string)
	GetUserID() string
	GetExpiresAt() time.Time
	SetExpiresAt(time.Time)
	GetRefreshUntil() time.Time
	Copy() SessionV2
}

// IsExpired returns true if the session is expired, and false otherwise.
func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

// IsRefreshable returns true if the session is refreshable, and false
// otherwise. A session is considered refreshable if it is not expired and its
// refresh until time is in the future.
func (s *Session) IsRefreshable() bool {
	return !s.IsExpired() && s.RefreshUntil.After(time.Now())
}

// IsValid returns true if the session is valid, and false otherwise. A session
// is considered valid if it is not expired.
func (s *Session) IsValid() bool {
	return !s.IsExpired()
}

type NewSession struct {
	UserID       string
	ExpiresAt    time.Time
	RefreshUntil time.Time
	Attributes   map[string]interface{}
}
