package handler

//user_handler.go 는 사용자 프로필 조회 및 업데이트, 전체 사용자 목록 조회 등의 기능을 담당하는 핸들러입니다.

import (
	"API_Server/internal/model"
	"API_Server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GET /users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	profile, err := h.userSvc.GetProfile(c.Request.Context(), userID.(string))
	if err != nil || profile == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// PUT /users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	profile, err := h.userSvc.UpdateProfile(c.Request.Context(), userID.(string), req)
	if err != nil || profile == nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GET /users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	userID, _ := c.Get("userID")

	users, err := h.userSvc.GetAllUsers(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: "INTERNAL", Message: "서버 오류"})
		return
	}

	if users == nil {
		users = []model.UsersProfile{}
	}
	c.JSON(http.StatusOK, users)
}
