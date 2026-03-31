package repository

import (
	"API_Server/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CandidateRepository struct {
	pool *pgxpool.Pool
}

func NewCandidateRepository(pool *pgxpool.Pool) *CandidateRepository {
	return &CandidateRepository{pool: pool}
}

func (r *CandidateRepository) Create(ctx context.Context, c *model.Candidate) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO candidates (candidate_id, session_id, user_id, type, ip, port, protocol, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		c.CandidateID, c.SessionID, c.UserID, c.Type, c.IP, c.Port, c.Protocol, c.CreatedAt,
	)
	return err
}

func (r *CandidateRepository) FindBySessionID(ctx context.Context, sessionID string) ([]model.Candidate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT candidate_id, session_id, user_id, type, ip, port, protocol, created_at
		 FROM candidates WHERE session_id = $1 ORDER BY created_at ASC`, sessionID,
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
	return candidates, nil
}
