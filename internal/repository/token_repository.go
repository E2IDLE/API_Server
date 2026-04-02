package repository

import (
	"API_Server/internal/model"
	"context"
	"database/sql"
	"errors"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, token *model.AuthToken) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO auth_tokens (token, user_id, expires_at) VALUES (?, ?, ?)`,
		token.Token, token.UserID, token.ExpiresAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*model.AuthToken, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT token, user_id, expires_at FROM auth_tokens
		 WHERE token = ? AND expires_at > datetime('now')`, token,
	)
	t := &model.AuthToken{}
	err := row.Scan(&t.Token, &t.UserID, &t.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (r *TokenRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth_tokens WHERE token = ?`, token)
	return err
}

func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth_tokens WHERE user_id = ?`, userID)
	return err
}
