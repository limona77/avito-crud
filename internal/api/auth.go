package api

import (
	repo "avito-crud/internal/repostiory/auth"
	"avito-crud/internal/service"
	service2 "avito-crud/internal/service/auth"
	"errors"

	"github.com/gin-gonic/gin"
)

type auth struct {
	authService service.IAuthService
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type AuthResponse struct {
	Token string `json:"token"`
}

func (a *auth) auth(ctx *gin.Context) {
	var req AuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": errors.New("incorrect format").Error()})
		return
	}

	token, err := a.authService.Auth(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			ctx.JSON(404, gin.H{"error": repo.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, service2.ErrInvalidCredentials) {
			ctx.JSON(401, gin.H{"error": service2.ErrInvalidCredentials.Error()})
			return
		}
		ctx.JSON(500, gin.H{"error": errors.New("internal server error").Error()})
		return
	}
	resp := AuthResponse{Token: token}
	ctx.JSON(200, resp)
}
