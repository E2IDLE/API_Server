package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string

	// TURN 서버 설정
	TurnHost   string
	TurnSecret string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/directp2p?sslmode=disable"),
		TurnHost:    getEnv("TURN_HOST", "turn:turn.directp2p.com:3478"),
		TurnSecret:  getEnv("TURN_SECRET", "supersecret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
