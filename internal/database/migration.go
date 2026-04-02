package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) {
	queries := []string{
		// users
		`CREATE TABLE IF NOT EXISTS users (
			id            TEXT PRIMARY KEY,
			email         TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			nickname      TEXT NOT NULL,
			profile_image TEXT,
			created_at    DATETIME NOT NULL DEFAULT (datetime('now'))
		)`,

		// auth_tokens
		`CREATE TABLE IF NOT EXISTS auth_tokens (
			token      TEXT PRIMARY KEY,
			user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			expires_at DATETIME NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_auth_tokens_user ON auth_tokens(user_id)`,

		// agents
		`CREATE TABLE IF NOT EXISTS agents (
			agent_id      TEXT PRIMARY KEY,
			user_id       TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			device_name   TEXT NOT NULL,
			platform      TEXT NOT NULL CHECK (platform IN ('windows','macos','linux')),
			agent_version TEXT NOT NULL,
			multiaddress  TEXT,
			registered_at DATETIME NOT NULL DEFAULT (datetime('now')),
			last_seen_at  DATETIME NOT NULL DEFAULT (datetime('now')),
			status        TEXT NOT NULL DEFAULT 'offline' CHECK (status IN ('online','offline'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_agents_user ON agents(user_id)`,

		// sessions
		`CREATE TABLE IF NOT EXISTS sessions (
			session_id     TEXT PRIMARY KEY,
			invite_code    TEXT UNIQUE NOT NULL,
			status         TEXT NOT NULL DEFAULT 'waiting'
			               CHECK (status IN ('waiting','connecting','connected','completed','error')),
			sender_id      TEXT NOT NULL REFERENCES users(id),
			receiver_id    TEXT REFERENCES users(id),
			sender_token   TEXT,
			receiver_token TEXT,
			created_at     DATETIME NOT NULL DEFAULT (datetime('now')),
			updated_at     DATETIME NOT NULL DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_sender   ON sessions(sender_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_receiver ON sessions(receiver_id)`,

		// candidates
		`CREATE TABLE IF NOT EXISTS candidates (
			candidate_id TEXT PRIMARY KEY,
			session_id   TEXT    NOT NULL REFERENCES sessions(session_id) ON DELETE CASCADE,
			user_id      TEXT    NOT NULL REFERENCES users(id),
			type         TEXT    NOT NULL CHECK (type IN ('host','srflx','relay')),
			ip           TEXT    NOT NULL,
			port         INTEGER NOT NULL,
			protocol     TEXT    NOT NULL CHECK (protocol IN ('udp','tcp')),
			created_at   DATETIME NOT NULL DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_candidates_session ON candidates(session_id)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("migration 실패: %v\nquery: %s", err, q)
		}
	}
	log.Println("✅ migrations 완료")
}
