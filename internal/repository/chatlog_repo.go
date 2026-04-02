package repository

import "database/sql"

type ChatLogRepository struct {
	db *sql.DB
}

func NewChatLogRepository(db *sql.DB) *ChatLogRepository {
	return &ChatLogRepository{db: db}
}

func (r *ChatLogRepository) Save(sessionID, role, content string) error {
	_, err := r.db.Exec(
		"INSERT INTO chat_logs (session_id, role, content) VALUES (?, ?, ?)",
		sessionID, role, content,
	)
	return err
}
