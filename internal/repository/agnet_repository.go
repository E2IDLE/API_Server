package repository

import (
	"API_Server/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentRepository struct {
	pool *pgxpool.Pool
}

func NewAgentRepository(pool *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{pool: pool}
}

func (r *AgentRepository) Create(ctx context.Context, userID string, agent *model.Agent) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO agents (agent_id, user_id, device_name, platform, agent_version, registered_at, last_seen_at, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		agent.AgentID, userID, agent.DeviceName, agent.Platform,
		agent.AgentVersion, agent.RegisteredAt, agent.LastSeenAt, agent.Status,
	)
	return err
}

func (r *AgentRepository) FindByUserID(ctx context.Context, userID string) ([]model.Agent, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT agent_id, device_name, platform, agent_version, registered_at, last_seen_at, status
		 FROM agents WHERE user_id = $1 ORDER BY registered_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []model.Agent
	for rows.Next() {
		var a model.Agent
		if err := rows.Scan(&a.AgentID, &a.DeviceName, &a.Platform,
			&a.AgentVersion, &a.RegisteredAt, &a.LastSeenAt, &a.Status); err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}
	return agents, nil
}

func (r *AgentRepository) UpdateLastSeen(ctx context.Context, agentID string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE agents SET last_seen_at = NOW(), status = 'online' WHERE agent_id = $1`,
		agentID,
	)
	return err
}
