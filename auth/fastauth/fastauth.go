package fastauth

import (
	"time"

	"github.com/lukeshay/g/auth"
	"github.com/valyala/fasthttp"
)

type contextKey string

const SessionContextKey contextKey = "session"

type Validate func(*fasthttp.RequestCtx, auth.Session) error

type CookieOptions struct {
	Name   string
	Path   string
	Secure bool
}

type FastAuth struct {
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

func New(options NewOptions) *FastAuth {
	return &FastAuth{
		service: auth.NewSessionService(auth.NewSessionServiceOptions{
			Adapter:   options.Adapter,
			Encrypter: options.Encrypter,
			Generator: options.Generator,
		}),
		cookieOptions: options.CookieOptions,
		validate:      options.Validate,
	}
}

func (a *FastAuth) Service() *auth.SessionService {
	return a.service
}

func (a *FastAuth) CreateNewSession(ctx *fasthttp.RequestCtx, newSession auth.Session) (auth.Session, error) {
	session, err := a.service.CreateSession(ctx, newSession.Copy())
	if err != nil {
		return nil, err
	}

	cookie, err := a.CreateCookie(session)
	if err != nil {
		return nil, err
	}

	ctx.Response.Header.SetCookie(cookie)

	ctx.SetUserValue(SessionContextKey, session)

	return session, nil
}

func (a *FastAuth) GetSession(ctx *fasthttp.RequestCtx) (auth.Session, error) {
	var err error
	session, ok := ctx.UserValue(SessionContextKey).(auth.Session)
	if !ok {
		value := ctx.Request.Header.Cookie(a.cookieOptions.Name)

		sessionID, err := a.service.DecryptSessionID(string(value))
		if err != nil {
			return nil, err
		}

		session, err = a.service.GetSession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
	}

	err = a.validate(ctx, session)
	if err != nil {
		return nil, err
	}

	ctx.SetUserValue(SessionContextKey, session)

	return session, nil
}

func (e *FastAuth) GetSessionAndRefresh(ctx *fasthttp.RequestCtx, expiresAt time.Time) (auth.Session, error) {
	session, err := e.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	if expiresAt.After(session.GetRefreshUntil()) {
		expiresAt = session.GetRefreshUntil()
	}

	session.SetExpiresAt(expiresAt)

	err = e.service.UpdateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	cookie, err := e.CreateCookie(session)
	if err != nil {
		return nil, err
	}

	ctx.Response.Header.SetCookie(cookie)

	ctx.SetUserValue(SessionContextKey, session)

	return session, nil
}

func (e *FastAuth) InvalidateSession(ctx *fasthttp.RequestCtx) error {
	session, err := e.GetSession(ctx)
	if err != nil {
		return nil
	}

	err = e.service.DeleteSession(ctx, session.GetSessionID())
	if err != nil {
		return err
	}

	ctx.Response.Header.DelCookie(e.cookieOptions.Name)
	ctx.Response.Header.DelClientCookie(e.cookieOptions.Name)

	ctx.SetUserValue(SessionContextKey, nil)

	return nil
}

func (e *FastAuth) CreateCookie(session auth.Session) (*fasthttp.Cookie, error) {
	encryptedSessionID, err := e.service.EncryptSessionID(session.GetSessionID())
	if err != nil {
		return nil, err
	}

	cookie := &fasthttp.Cookie{}

	cookie.SetValue(encryptedSessionID)
	cookie.SetSecure(e.cookieOptions.Secure)
	cookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	cookie.SetExpire(session.GetExpiresAt())
	cookie.SetPath(e.cookieOptions.Path)
	cookie.SetKey(e.cookieOptions.Name)

	return cookie, nil
}
