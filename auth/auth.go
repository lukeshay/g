package auth

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const SessionContextKey contextKey = "session"

type Validate func(context.Context, *http.Request, Session) (context.Context, error)

type Auth struct {
	service  *SessionService
	validate Validate
}

type NewOptions struct {
	Adapter       Adapter
	Encrypter     Encrypter
	Generator     Generator
	CookieOptions CookieOptions
	Validate      Validate
}

func New(options NewOptions) *Auth {
	return &Auth{
		service: NewSessionService(NewSessionServiceOptions{
			Adapter:       options.Adapter,
			Encrypter:     options.Encrypter,
			Generator:     options.Generator,
			CookieOptions: options.CookieOptions,
		}),
		validate: options.Validate,
	}
}

func (a *Auth) Service() *SessionService {
	return a.service
}

func (a *Auth) CreateNewSession(ctx context.Context, w http.ResponseWriter, newSession Session) (context.Context, Session, error) {
	session, err := a.service.CreateSession(ctx, newSession)
	if err != nil {
		return ctx, nil, err
	}

	cookie, err := a.service.CreateCookie(session)
	if err != nil {
		return ctx, nil, err
	}

	http.SetCookie(w, cookie)

	ctx = context.WithValue(ctx, SessionContextKey, session)

	return ctx, session, nil
}

func (e *Auth) GetSession(ctx context.Context, r *http.Request) (context.Context, Session, error) {
	var err error
	session, ok := ctx.Value(SessionContextKey).(Session)
	if !ok {
		session, err = e.service.GetSessionFromCookies(ctx, r.Cookies())
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

func (e *Auth) GetSessionAndRefresh(ctx context.Context, w http.ResponseWriter, r *http.Request, expiresAt time.Time) (context.Context, Session, error) {
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

	cookie, err := e.service.CreateCookie(session)
	if err != nil {
		return ctx, nil, err
	}

	http.SetCookie(w, cookie)

	ctx = context.WithValue(ctx, SessionContextKey, session)

	return ctx, session, nil
}

func (e *Auth) InvalidateSession(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ctx, session, err := e.GetSession(ctx, r)
	if err != nil {
		return ctx, nil
	}

	err = e.service.DeleteSession(ctx, session.GetSessionID())
	if err != nil {
		return ctx, err
	}

	http.SetCookie(w, e.service.EmptyCookie())

	ctx = context.WithValue(ctx, SessionContextKey, nil)

	return ctx, nil
}
