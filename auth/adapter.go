package auth

import (
	"context"
)

// Adapter is responsible for creating, updating, and deleting sessions from a
// datastore of your choice. This can be a SQL database, a NoSQL database, an
// in-memory cache, etc. It is up to you to implement this interface for your
// specific use case. The adapter should just be simple operations. The logic
// for wether for calling these operations are handled elsewhere.
type Adapter interface {
	// GetSession retrieves the session with the given sessionID from the
	// database and returns it.
	GetSession(ctx context.Context, sessionID string) (Session, error)
	// InsertSession inserts a new session into the database with the given
	// values.
	InsertSession(ctx context.Context, newSession Session) error
	// UpdateSession updates the session with the given sessionID to have a new
	// expiration time.
	UpdateSession(ctx context.Context, newSession Session) error
	// DeleteSessionsByUserID invalidates all sessions for the given user.
	DeleteSessionsByUserID(ctx context.Context, userID string) error
	// DeleteSession invalidates the session with the given sessionID.
	DeleteSession(ctx context.Context, sessionID string) error
}
