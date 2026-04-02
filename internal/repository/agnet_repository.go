package repository

import (
	"API_Server/internal/model"
	"context"
	"database/sql"
	"time"
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
		 multiaddress, registered_at, last_seen_at, status
		 FROM agents WHERE user_id = ? ORDER BY registered_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []model.Agent
	for rows.Next() {
		var a model.Agent
		var registeredAt string
		var lastSeenAt sql.NullString
		var multiAddr sql.NullString

		if err := rows.Scan(&a.AgentID, &a.DeviceName, &a.Platform,
			&a.AgentVersion, &multiAddr, &registeredAt, &lastSeenAt, &a.Status); err != nil {
			return nil, err
		}

		// 문자열 → time.Time 변환
		t, err := time.Parse("2006-01-02 15:04:05", registeredAt)
		if err != nil {
			return nil, err
		}
		a.RegisteredAt = t

		if lastSeenAt.Valid {
			t2, err := time.Parse("2006-01-02 15:04:05", lastSeenAt.String)
			if err != nil {
				return nil, err
			}
			a.LastSeenAt = &t2
		}

		if multiAddr.Valid {
			a.MultiAddress = &multiAddr.String
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
