package repository

import (
	"API_Server/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	pool *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{pool: pool}
}

func (r *TokenRepository) Create(ctx context.Context, token *model.AuthToken) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO auth_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		token.Token, token.UserID, token.ExpiresAt,
	)
	return err
}

func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*model.AuthToken, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT token, user_id, expires_at FROM auth_tokens
		 WHERE token = $1 AND expires_at > NOW()`, token,
	)
	t := &model.AuthToken{}
	err := row.Scan(&t.Token, &t.UserID, &t.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (r *TokenRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM auth_tokens WHERE token = $1`, token)
	return err
}

func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM auth_tokens WHERE user_id = $1`, userID)
	return err
}
