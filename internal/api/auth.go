package api

import (
	repo "avito-crud/internal/repostiory/auth"
	"avito-crud/internal/service"
	service2 "avito-crud/internal/service/auth"
	"errors"
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
		c.JSON(400, gin.H{"error": errors.New("incorrect format").Error()})
		return
	}

	token, err := a.authService.Login(c, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			c.JSON(401, gin.H{"error": repo.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, service2.ErrInvalidCredentials) {
			c.JSON(401, gin.H{"error": service2.ErrInvalidCredentials.Error()})
			return
		}
		c.JSON(500, gin.H{"error": errors.New("internal server error").Error()})
		return
	}
	resp := AuthResponse{Token: token}
	c.JSON(200, resp)
}
