package repository

import (
	"API_Server/internal/model"
	"context"
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password_hash, nickname, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.PasswordHash, user.Nickname, user.CreatedAt,
	)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, nickname, profile_image, created_at
		 FROM users WHERE email = ?`, email,
	)
	u := &model.User{}
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Nickname, &u.ProfileImage, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, nickname, profile_image, created_at
		 FROM users WHERE id = ?`, id,
	)
	u := &model.User{}
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Nickname, &u.ProfileImage, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET nickname = ?, profile_image = ? WHERE id = ?`,
		user.Nickname, user.ProfileImage, user.ID,
	)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = ? WHERE id = ?`,
		passwordHash, userID,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}

func (r *UserRepository) FindAll(ctx context.Context, excludeUserID string) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, email, nickname, profile_image, created_at
		 FROM users WHERE id != ?
		 ORDER BY created_at DESC`, excludeUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Nickname, &u.ProfileImage, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
