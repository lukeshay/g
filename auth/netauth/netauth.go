package netauth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lukeshay/g/auth"
)

type contextKey string

const SessionContextKey contextKey = "session"

type Validate func(context.Context, *http.Request, auth.Session) (context.Context, error)

type CookieOptions struct {
	Name   string
	Path   string
	Secure bool
}

type NetAuth struct {
	service       *auth.SessionService
	validate      Validate
	cookieOptions CookieOptions
}

type NewOptions struct {
	Adapter       auth.SessionAdapter
	Encrypter     auth.Encrypter
	Generator     auth.Generator
	CookieOptions CookieOptions
	Validate      Validate
}

func New(options NewOptions) *NetAuth {
	return &NetAuth{
		service: auth.NewSessionService(auth.NewSessionServiceOptions{
			Adapter:   options.Adapter,
			Encrypter: options.Encrypter,
			Generator: options.Generator,
		}),
		validate:      options.Validate,
		cookieOptions: options.CookieOptions,
	}
}

func (a *NetAuth) Service() *auth.SessionService {
	return a.service
}

func (a *NetAuth) CreateNewSession(ctx context.Context, w http.ResponseWriter, newSession auth.Session) (context.Context, auth.Session, error) {
	session, err := a.service.CreateSession(ctx, newSession)
	if err != nil {
		return ctx, nil, err
	}

	cookie, err := a.CreateCookie(session)
	if err != nil {
		return ctx, nil, err
	}

	http.SetCookie(w, cookie)

	ctx = context.WithValue(ctx, SessionContextKey, session)

	return ctx, session, nil
}

func (e *NetAuth) GetSession(ctx context.Context, r *http.Request) (context.Context, auth.Session, error) {
	var err error
	session, ok := ctx.Value(SessionContextKey).(auth.Session)
	if !ok {
		session, err = e.GetSessionFromCookies(ctx, r.Cookies())
		if err != nil {
			return ctx, nil, err
		}
	}

	ctx, err = e.validate(ctx, r, session)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(ctx, SessionContextKey, session)

	return ctx, session, nil
}

func (e *NetAuth) GetSessionAndRefresh(ctx context.Context, w http.ResponseWriter, r *http.Request, expiresAt time.Time) (context.Context, auth.Session, error) {
	ctx, session, err := e.GetSession(ctx, r)
	if err != nil {
		return ctx, nil, err
	}

	if expiresAt.After(session.GetRefreshUntil()) {
		expiresAt = session.GetRefreshUntil()
	}

	session.SetExpiresAt(expiresAt)

	err = e.service.UpdateSession(ctx, session)
	if err != nil {
		return ctx, nil, err
	}

	cookie, err := e.CreateCookie(session)
	if err != nil {
		return ctx, nil, err
	}

	http.SetCookie(w, cookie)

	ctx = context.WithValue(ctx, SessionContextKey, session)

	return ctx, session, nil
}

func (e *NetAuth) InvalidateSession(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx, session, err := e.GetSession(ctx, r)
	if err != nil {
		return ctx, nil
	}

	err = e.service.DeleteSession(ctx, session.GetSessionID())
	if err != nil {
		return ctx, err
	}

	http.SetCookie(w, e.EmptyCookie())

	ctx = context.WithValue(ctx, SessionContextKey, nil)

	return ctx, nil
}

func (a *NetAuth) GetSessionFromCookies(ctx context.Context, cookies []*http.Cookie) (auth.Session, error) {
	for _, cookie := range cookies {
		if cookie.Name == a.cookieOptions.Name {
			return a.GetSessionFromCookie(ctx, cookie)
		}
	}

	return nil, fmt.Errorf("session not found")
}

func (a *NetAuth) GetSessionFromCookie(ctx context.Context, cookie *http.Cookie) (auth.Session, error) {
	decryptedSessionID, err := a.service.DecryptSessionID(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("error decrypting session id: %s", err.Error())
	}

	return a.service.GetSession(ctx, decryptedSessionID)
}

func (a *NetAuth) CreateCookie(session auth.Session) (*http.Cookie, error) {
	encryptedSessionID, err := a.service.EncryptSessionID(session.GetSessionID())
	if err != nil {
		return a.createCookie(time.Now(), ""), fmt.Errorf("error encrypting session id: %s", err.Error())
	}

	return a.createCookie(session.GetExpiresAt(), encryptedSessionID), nil
}

func (a *NetAuth) EmptyCookie() *http.Cookie {
	return a.createCookie(time.Now(), "")
}

func (a *NetAuth) createCookie(expiresAt time.Time, value string) *http.Cookie {
	return &http.Cookie{
		Expires:  expiresAt,
		HttpOnly: true,
		Name:     a.cookieOptions.Name,
		Path:     a.cookieOptions.Path,
		Secure:   a.cookieOptions.Secure,
		Value:    value,
	}
}
