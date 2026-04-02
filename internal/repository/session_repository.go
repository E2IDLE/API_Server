package repository

import (
	"API_Server/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) Create(ctx context.Context, session *model.Session) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO sessions (session_id, invite_code, status, sender_id, receiver_id, sender_token, receiver_token, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		session.SessionID, session.InviteCode, session.Status,
		session.SenderID, session.ReceiverID,
		session.SenderToken, session.ReceiverToken,
		session.CreatedAt, session.UpdatedAt,
	)
	return err
}

func (r *SessionRepository) FindByID(ctx context.Context, sessionID string) (*model.Session, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT session_id, invite_code, status, sender_id, receiver_id, 
		        sender_token, receiver_token, created_at, updated_at
		 FROM sessions WHERE session_id = $1`, sessionID,
	)
	s := &model.Session{}
	err := row.Scan(&s.SessionID, &s.InviteCode, &s.Status,
		&s.SenderID, &s.ReceiverID,
		&s.SenderToken, &s.ReceiverToken,
		&s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return s, err
}

func (r *SessionRepository) UpdateReceiver(ctx context.Context, sessionID, receiverID, receiverToken, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE sessions SET receiver_id = $1, receiver_token = $2, status = $3, updated_at = NOW()
		 WHERE session_id = $4`,
		receiverID, receiverToken, status, sessionID,
	)
	return err
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, sessionID, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE sessions SET status = $1, updated_at = NOW() WHERE session_id = $2`,
		status, sessionID,
	)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE session_id = $1`, sessionID)
	return err
}

func (r *SessionRepository) FindByUserIDPaginated(ctx context.Context, userID string, page, pageSize int) ([]model.Session, int, error) {
	var totalCount int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM sessions WHERE sender_id = $1 OR receiver_id = $1`, userID,
	).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.pool.Query(ctx,
		`SELECT session_id, invite_code, status, sender_id, receiver_id,
		        sender_token, receiver_token, created_at, updated_at
		 FROM sessions WHERE sender_id = $1 OR receiver_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`, userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sessions []model.Session
	for rows.Next() {
		var s model.Session
		if err := rows.Scan(&s.SessionID, &s.InviteCode, &s.Status,
			&s.SenderID, &s.ReceiverID,
			&s.SenderToken, &s.ReceiverToken,
			&s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		sessions = append(sessions, s)
	}
	return sessions, totalCount, nil
}
