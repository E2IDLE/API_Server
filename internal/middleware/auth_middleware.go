package middleware

// AuthMiddleware 는 인증이 필요한 라우트에 적용되는 미들웨어입니다.

import (
	"API_Server/internal/model"
	"API_Server/internal/repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenRepo *repository.TokenRepository
}

func NewAuthMiddleware(tokenRepo *repository.TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{tokenRepo: tokenRepo}
}

// Required 는 인증이 필수인 라우트에 사용
func (m *AuthMiddleware) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "인증 토큰이 필요합니다.",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		authToken, err := m.tokenRepo.FindByToken(c.Request.Context(), token)
		if err != nil || authToken == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "유효하지 않은 토큰입니다.",
			})
			return
		}

		// 컨텍스트에 사용자 ID 저장
		c.Set("userID", authToken.UserID)
		c.Set("token", token)
		c.Next()
	}
}
