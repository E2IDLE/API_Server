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
// API 응답용 구조체
type UsersProfile struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname"`
	ProfileImage *string   `json:"profileImage"`
	CreatedAt    time.Time `json:"createdAt"`
	AgentStatus  string    `json:"agentStatus"`
}

// 프로필 수정 요청
type UpdateProfileRequest struct {
	Nickname     *string `json:"nickname"`
	ProfileImage *string `json:"profileImage"`
}

// ══════════════════════════════════════
// Agents
// ══════════════════════════════════════

// 8. RegisterAgentRequest
type RegisterAgentRequest struct {
	DeviceName   string `json:"deviceName" binding:"required,max=100"`
	Platform     string `json:"platform" binding:"required,oneof=windows macos linux"`
	AgentVersion string `json:"agentVersion" binding:"required"`
	MultiAddress string `json:"multiaddress" binding:"required" `
}

// 9. Agent
type Agent struct {
	AgentID      string     `json:"agentId"`
	DeviceName   string     `json:"deviceName"`
	Platform     string     `json:"platform"`
	AgentVersion string     `json:"agentVersion"`
	RegisteredAt time.Time  `json:"registeredAt"`
	LastSeenAt   *time.Time `json:"lastSeenAt"`   // NULL 허용
	MultiAddress *string    `json:"multiaddress"` // NULL 허용
	Status       string     `json:"status"`
}

// ══════════════════════════════════════
// Sessions
// ══════════════════════════════════════

// 10. Session
type Session struct {
	SessionID     string    `json:"sessionId"`
	InviteCode    string    `json:"inviteCode"`
	Status        string    `json:"status"` // 상태
	SenderID      string    `json:"senderId"`
	ReceiverID    *string   `json:"receiverId"` // nullable
	SenderToken   string    `json:"-"`
	ReceiverToken string    `json:"-"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

//. CreateSessionRequest 연결 신청 요청
type CreateSessionRequest struct {
	UserID string `json:"userId"` // 연결 신청 대상 유저 ID (선택)
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
	ID           string    `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Nickname     string    `json:"nickname"`
	ProfileImage *string   `json:"profileImage"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

//친구 추가 및 기존 유저 관련
type Friend struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	FriendID  uint      `json:"friend_id" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at"`
}
