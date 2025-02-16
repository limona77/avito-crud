package api

import (
	"avito-crud/internal/model"
	authRepo "avito-crud/internal/repostiory/auth"
	"avito-crud/internal/service"
	shopService "avito-crud/internal/service/shop"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type info struct {
	infoService service.IInfoService
}

type InfoResponse struct {
	Info model.UserInfo `json:"info"`
}

const (
	TokenPrefix = "Bearer"
)

func (i *info) getInfo(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != TokenPrefix {
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
		if errors.Is(err, shopService.ErrInvalidToken) {
			ctx.JSON(401, gin.H{"error": shopService.ErrInvalidToken.Error()})
			return
		}
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, *userInfo)
}
