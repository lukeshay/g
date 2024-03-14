package main

import (
	"context"

	"github.com/lukeshay/g/auth"
	"github.com/uptrace/bun"
)

type Adapter struct {
	db *bun.DB
}

var _ auth.Adapter = (*Adapter)(nil)

func (a *Adapter) GetSession(ctx context.Context, sessionID string) (auth.Session, error) {
	session := new(Session)
	err := a.db.NewSelect().Model(session).Where("id = ?", sessionID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (a *Adapter) InsertSession(ctx context.Context, newSession auth.Session) error {
	_, err := a.db.NewInsert().Model(newSession).Exec(ctx)

	return err
}

func (a *Adapter) UpdateSession(ctx context.Context, newSession auth.Session) error {
	_, err := a.db.NewUpdate().Model(newSession).Column("expires_at").WherePK().Exec(ctx)

	return err
}

func (a *Adapter) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	_, err := a.db.NewDelete().Model((*Session)(nil)).Where("user_id = ?", userID).Exec(ctx)

	return err
}

func (a *Adapter) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := a.db.NewDelete().Model(&Session{ID: sessionID}).WherePK().Exec(ctx)

	return err
}
