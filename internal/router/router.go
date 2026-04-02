package router

import (
	"API_Server/internal/handler"
	"API_Server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	authH *handler.AuthHandler,
	userH *handler.UserHandler,
	agentH *handler.AgentHandler,
	sessionH *handler.SessionHandler,
	candidateH *handler.CandidateHandler,
	turnH *handler.TurnHandler,
	wsH *handler.WSHandler,
	authMw *middleware.AuthMiddleware,
) *gin.Engine {

	r := gin.Default()

	// ── Auth (인증 불필요) ──
	auth := r.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	// ── Auth (인증 필요) ──
	authProtected := r.Group("/auth", authMw.Required())
	{
		authProtected.POST("/logout", authH.Logout)
		authProtected.DELETE("/me", authH.DeleteAccount)
		authProtected.PUT("/password", authH.ChangePassword)
	}

	// ── Users ──
	users := r.Group("/users", authMw.Required())
	{
		users.GET("/me", userH.GetProfile)
		users.PUT("/me", userH.UpdateProfile)
		users.GET("", userH.GetAllUsers)

		// ── Agents (nested) ──
		users.POST("/me/agents", agentH.RegisterAgent)
		users.GET("/me/agents", agentH.ListAgents)
	}

	// ── Sessions ──
	sessions := r.Group("/sessions", authMw.Required())
	{
		sessions.POST("", sessionH.CreateSession)
		sessions.GET("/history", sessionH.GetHistory)
		sessions.POST("/:sessionId/join", sessionH.JoinSession)
		sessions.GET("/:sessionId", sessionH.GetSession)
		sessions.DELETE("/:sessionId", sessionH.DeleteSession)

		// ── Candidates (nested) ──
		sessions.POST("/:sessionId/candidates", candidateH.RegisterCandidate)
		sessions.GET("/:sessionId/candidates", candidateH.ListCandidates)

		// ── TURN (nested) ──
		sessions.POST("/:sessionId/turn-credentials", turnH.IssueTurnCredentials)
	}

	// ── WebSocket ──
	r.GET("/ws", wsH.HandleWebSocket)

	return r
}
