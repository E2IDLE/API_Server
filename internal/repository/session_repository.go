package repository

import (
	"API_Server/internal/model"
	"context"
	"database/sql"
	"errors"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *model.Session) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO sessions (session_id, invite_code, status, sender_id, receiver_id, sender_token, receiver_token, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		session.SessionID, session.InviteCode, session.Status,
		session.SenderID, session.ReceiverID,
		session.SenderToken, session.ReceiverToken,
		session.CreatedAt, session.UpdatedAt,
	)
	return err
}

func (r *SessionRepository) FindByID(ctx context.Context, sessionID string) (*model.Session, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT session_id, invite_code, status, sender_id, receiver_id,
		        sender_token, receiver_token, created_at, updated_at
		 FROM sessions WHERE session_id = ?`, sessionID,
	)
	s := &model.Session{}
	err := row.Scan(&s.SessionID, &s.InviteCode, &s.Status,
		&s.SenderID, &s.ReceiverID,
		&s.SenderToken, &s.ReceiverToken,
		&s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return s, err
}

func (r *SessionRepository) UpdateReceiver(ctx context.Context, sessionID, receiverID, receiverToken, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET receiver_id = ?, receiver_token = ?, status = ?, updated_at = datetime('now')
		 WHERE session_id = ?`,
		receiverID, receiverToken, status, sessionID,
	)
	return err
}

func (r *SessionRepository) UpdateStatus(ctx context.Context, sessionID, status string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET status = ?, updated_at = datetime('now') WHERE session_id = ?`,
		status, sessionID,
	)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE session_id = ?`, sessionID)
	return err
}

func (r *SessionRepository) FindByUserIDPaginated(ctx context.Context, userID string, page, pageSize int) ([]model.Session, int, error) {
	var totalCount int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM sessions WHERE sender_id = ? OR receiver_id = ?`, userID, userID,
	).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx,
		`SELECT session_id, invite_code, status, sender_id, receiver_id,
		        sender_token, receiver_token, created_at, updated_at
		 FROM sessions WHERE sender_id = ? OR receiver_id = ?
		 ORDER BY created_at DESC
		 LIMIT ? OFFSET ?`, userID, userID, pageSize, offset,
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
	return sessions, totalCount, rows.Err()
}
