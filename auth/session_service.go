package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// SessionService is a struct that contains all of the necessary components to manage
// server side sessions. It is responsible for creating, retrieving, and
// invalidating sessions. It is also responsible for encrypting and decrypting
// session IDs. It is up to you to choose the implementation of the Adapter,
// Encrypter, and Generator interfaces. The following links provide
// implementations.
//
//   - [Adapters](./adapters)
//   - [Encrypter](./encrypters)
//   - [Generator](./generators)
type SessionService struct {
	adapter       Adapter
	encrypter     Encrypter
	generator     Generator
	cookieOptions CookieOptions
}

type CookieOptions struct {
	Name   string
	Path   string
	Secure bool
}

type NewSessionServiceOptions struct {
	Adapter       Adapter
	Encrypter     Encrypter
	Generator     Generator
	CookieOptions CookieOptions
}

// NewSessionService returns a new instance of Auth.
func NewSessionService(options NewSessionServiceOptions) *SessionService {
	return &SessionService{
		adapter:       options.Adapter,
		encrypter:     options.Encrypter,
		generator:     options.Generator,
		cookieOptions: options.CookieOptions,
	}
}

// 4. Return the session and cookie
func (a *SessionService) CreateSession(ctx context.Context, newSession Session) (Session, error) {
	sessionID, err := a.generator.Generate()
	if err != nil {
		return nil, fmt.Errorf("error generating session id: %s", err.Error())
	}

	insertedSession := newSession.Copy()
	insertedSession.SetSessionID(sessionID)

	err = a.adapter.InsertSession(ctx, insertedSession)
	if err != nil {
		return nil, fmt.Errorf("error inserting session: %s", err.Error())
	}

	return insertedSession, nil
}

func (a *SessionService) GetSession(ctx context.Context, sessionID string) (Session, error) {
	session, err := a.adapter.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %s", err.Error())
	}

	if session.GetExpiresAt().Before(time.Now()) {
		return nil, fmt.Errorf("session is expired: %s", session.GetExpiresAt().Format(time.RFC3339))
	}

	return session, nil
}

func (a *SessionService) GetSessionFromCookies(ctx context.Context, cookies []*http.Cookie) (Session, error) {
	for _, cookie := range cookies {
		if cookie.Name == a.cookieOptions.Name {
			return a.GetSessionFromCookie(ctx, cookie)
		}
	}

	return nil, fmt.Errorf("session not found")
}

func (a *SessionService) GetSessionFromCookie(ctx context.Context, cookie *http.Cookie) (Session, error) {
	decryptedSessionID, err := a.DecryptSessionID(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("error decrypting session id: %s", err.Error())
	}

	return a.GetSession(ctx, decryptedSessionID)
}

func (a *SessionService) UpdateSession(ctx context.Context, session Session) error {
	return a.adapter.UpdateSession(ctx, session)
}

func (a *SessionService) DeleteSession(ctx context.Context, sessionID string) error {
	return a.adapter.DeleteSession(ctx, sessionID)
}

func (a *SessionService) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	return a.adapter.DeleteSessionsByUserID(ctx, userID)
}

func (a *SessionService) EncryptSessionID(sessionID string) (string, error) {
	return a.encrypter.Encrypt(sessionID)
}

func (a *SessionService) DecryptSessionID(encryptedSessionID string) (string, error) {
	return a.encrypter.Decrypt(encryptedSessionID)
}

func (a *SessionService) CreateCookie(session Session) (*http.Cookie, error) {
	encryptedSessionID, err := a.EncryptSessionID(session.GetSessionID())
	if err != nil {
		return a.createCookie(time.Now(), ""), fmt.Errorf("error encrypting session id: %s", err.Error())
	}

	return a.createCookie(session.GetExpiresAt(), encryptedSessionID), nil
}

func (a *SessionService) EmptyCookie() *http.Cookie {
	return a.createCookie(time.Now(), "")
}

func (a *SessionService) createCookie(expiresAt time.Time, value string) *http.Cookie {
	return &http.Cookie{
		Expires:  expiresAt,
		HttpOnly: true,
		Name:     a.cookieOptions.Name,
		Path:     a.cookieOptions.Path,
		Secure:   a.cookieOptions.Secure,
		Value:    value,
	}
}
