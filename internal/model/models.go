package model

import "time"

// ══════════════════════════════════════
// Auth
// ══════════════════════════════════════

// 1. RegisterRequest
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Nickname string `json:"nickname" binding:"required,min=2,max=30"`
}

// 2. RegisterResponse
type RegisterResponse struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

// 3. LoginRequest
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 4. LoginResponse
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// 5. ChangePasswordRequest
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8,max=72"`
}

// ══════════════════════════════════════
// Users
// ══════════════════════════════════════

// 6. UserProfile
type UserProfile struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname"`
	ProfileImage *string   `json:"profileImage"` // nullable
	CreatedAt    time.Time `json:"createdAt"`
	AgentStatus  string    `json:"agentStatus"` // online | offline
}

// 7. UpdateProfileRequest
type UpdateProfileRequest struct {
	Nickname     *string `json:"nickname" binding:"omitempty,min=2,max=30"`
	ProfileImage *string `json:"profileImage"` // nullable → null이면 삭제
}

// ══════════════════════════════════════
// Agents
// ══════════════════════════════════════

// 8. RegisterAgentRequest
type RegisterAgentRequest struct {
	DeviceName   string `json:"deviceName" binding:"required,max=100"`
	Platform     string `json:"platform" binding:"required,oneof=windows macos linux"`
	AgentVersion string `json:"agentVersion" binding:"required"`
}

// 9. Agent
type Agent struct {
	AgentID      string    `json:"agentId"`
	DeviceName   string    `json:"deviceName"`
	Platform     string    `json:"platform"`
	AgentVersion string    `json:"agentVersion"`
	RegisteredAt time.Time `json:"registeredAt"`
	LastSeenAt   time.Time `json:"lastSeenAt"`
	Status       string    `json:"status"` // online | offline
}

// ══════════════════════════════════════
// Sessions
// ══════════════════════════════════════

// 10. Session
type Session struct {
	SessionID  string    `json:"sessionId"`
	InviteCode string    `json:"inviteCode"`
	Status     string    `json:"status"` // waiting | connecting | connected | completed | error
	SenderID   string    `json:"senderId"`
	ReceiverID *string   `json:"receiverId"` // nullable
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// 11. JoinSessionRequest
type JoinSessionRequest struct {
	InviteCode string `json:"inviteCode" binding:"required"`
}

// 12. SessionHistoryResponse
type SessionHistoryResponse struct {
	Items      []Session  `json:"items"`
	Pagination Pagination `json:"pagination"`
}

// 13. Pagination
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

// ══════════════════════════════════════
// Candidates
// ══════════════════════════════════════

// 14. RegisterCandidateRequest
type RegisterCandidateRequest struct {
	Type     string `json:"type" binding:"required,oneof=host srflx relay"` // host | srflx | relay
	IP       string `json:"ip" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Protocol string `json:"protocol" binding:"required,oneof=udp tcp"` // udp | tcp
}

// 15. Candidate
type Candidate struct {
	CandidateID string    `json:"candidateId"`
	SessionID   string    `json:"sessionId"`
	UserID      string    `json:"userId"`
	Type        string    `json:"type"`
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	Protocol    string    `json:"protocol"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ══════════════════════════════════════
// TURN
// ══════════════════════════════════════

// 16. TurnCredentials
type TurnCredentials struct {
	URLs       []string  `json:"urls"`
	Username   string    `json:"username"`
	Credential string    `json:"credential"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

// ══════════════════════════════════════
// Error
// ══════════════════════════════════════

// 17. ErrorResponse
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ══════════════════════════════════════
// WebSocket 이벤트
// ══════════════════════════════════════

type WSMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// ── DB 내부 모델 (password hash 등) ──

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Nickname     string
	ProfileImage *string
	CreatedAt    time.Time
}

type AuthToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}
