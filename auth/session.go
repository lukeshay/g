package auth

import (
	"time"
)

// Session is an interface that represents a user session. It contains the
// functions required for this package to interact with the session. The
// Session interface is used to allow for different implementations of the
// session to be used.
//
// We recommend storing the user agent and IP address so you can verify that
// the session is being used by the same device and from the same location.
type Session interface {
	GetSessionID() string
	SetSessionID(string)
	GetUserID() string
	GetExpiresAt() time.Time
	SetExpiresAt(time.Time)
	GetRefreshUntil() time.Time
	Copy() Session
}
