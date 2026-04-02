package repository

import (
	"API_Server/internal/model"
	"context"
	"database/sql"
)

type CandidateRepository struct {
	db *sql.DB
}

func NewCandidateRepository(db *sql.DB) *CandidateRepository {
	return &CandidateRepository{db: db}
}

func (r *CandidateRepository) Create(ctx context.Context, c *model.Candidate) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO candidates (candidate_id, session_id, user_id, type, ip, port, protocol, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		c.CandidateID, c.SessionID, c.UserID, c.Type, c.IP, c.Port, c.Protocol, c.CreatedAt,
	)
	return err
}

func (r *CandidateRepository) FindBySessionID(ctx context.Context, sessionID string) ([]model.Candidate, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT candidate_id, session_id, user_id, type, ip, port, protocol, created_at
		 FROM candidates WHERE session_id = ? ORDER BY created_at ASC`, sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candidates []model.Candidate
	for rows.Next() {
		var c model.Candidate
		if err := rows.Scan(&c.CandidateID, &c.SessionID, &c.UserID,
			&c.Type, &c.IP, &c.Port, &c.Protocol, &c.CreatedAt); err != nil {
			return nil, err
		}
		candidates = append(candidates, c)
	}
	return candidates, rows.Err()
}
