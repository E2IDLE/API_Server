package repository

import (
	"API_Server/internal/model"
	"context"
	"database/sql"
)

type AgentRepository struct {
	db *sql.DB
}

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) Create(ctx context.Context, userID string, agent *model.Agent) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO agents (agent_id, user_id, device_name, platform, agent_version, registered_at, last_seen_at, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		agent.AgentID, userID, agent.DeviceName, agent.Platform,
		agent.AgentVersion, agent.RegisteredAt, agent.LastSeenAt, agent.Status,
	)
	return err
}

func (r *AgentRepository) FindByUserID(ctx context.Context, userID string) ([]model.Agent, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT agent_id, device_name, platform, agent_version, 
     COALESCE(multiaddress, ''), registered_at, 
     COALESCE(last_seen_at, registered_at), status
     FROM agents WHERE user_id = ? ORDER BY registered_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []model.Agent
	for rows.Next() {
		var a model.Agent
		if err := rows.Scan(&a.AgentID, &a.DeviceName, &a.Platform,
			&a.AgentVersion, &a.MultiAddress, &a.RegisteredAt, &a.LastSeenAt, &a.Status); err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}
	return agents, rows.Err()
}

func (r *AgentRepository) UpdateLastSeen(ctx context.Context, agentID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE agents SET last_seen_at = datetime('now'), status = 'online' WHERE agent_id = ?`,
		agentID,
	)
	return err
}
