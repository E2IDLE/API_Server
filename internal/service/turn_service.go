package service

import (
	"API_Server/internal/config"
	"API_Server/internal/model"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
	//"honnef.co/go/tools/config"
)

type TurnService struct {
	cfg *config.Config
}

func NewTurnService(cfg *config.Config) *TurnService {
	return &TurnService{cfg: cfg}
}

// IssueTurnCredentials 는 time-limited TURN credentials(RFC 보안 방식)을 발급
func (s *TurnService) IssueTurnCredentials(sessionID string) (*model.TurnCredentials, error) {
	ttl := 24 * time.Hour
	expiresAt := time.Now().UTC().Add(ttl)
	timestamp := expiresAt.Unix()

	// username = "timestamp:sessionId"
	username := fmt.Sprintf("%d:%s", timestamp, sessionID)

	// HMAC-SHA1 credential
	mac := hmac.New(sha1.New, []byte(s.cfg.TurnSecret))
	mac.Write([]byte(username))
	credential := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return &model.TurnCredentials{
		URLs:       []string{s.cfg.TurnHost},
		Username:   username,
		Credential: credential,
		ExpiresAt:  expiresAt,
	}, nil
}
