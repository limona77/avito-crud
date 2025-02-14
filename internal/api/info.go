package api

import (
	"avito-crud/internal/model"
	authRepo "avito-crud/internal/repostiory/auth"
	"avito-crud/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

type info struct {
	infoService service.IInfoService
	log         *slog.Logger
}

type InfoResponse struct {
	Info model.UserInfo `json:"info"`
}

func (i *info) getInfo(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Проверяем, что токен начинается с "Bearer "
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		ctx.JSON(400, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	token := tokenParts[1]

	userInfo, err := i.infoService.GetInfo(ctx, token)
	if err != nil {
		if errors.Is(err, authRepo.ErrUserNotFound) {
			ctx.JSON(404, gin.H{"error": authRepo.ErrUserNotFound.Error()})
			return
		}
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, *userInfo)
}
