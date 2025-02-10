package api

import (
	"avito-crud/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type auth struct {
	authService service.IAuthService
	log         *slog.Logger
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type AuthResponse struct {
	Token string `json:"token"`
}

func newAuthRoutes(r *gin.RouterGroup, authService service.IAuthService) {
	a := &auth{authService: authService}
	r.POST("/auth", a.auth)
}

func (a *auth) auth(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := a.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp := AuthResponse{Token: token}
	c.JSON(200, resp)
}
