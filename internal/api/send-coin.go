package api

import (
	authRepo "avito-crud/internal/repostiory/auth"
	shopRepo "avito-crud/internal/repostiory/shop"
	"avito-crud/internal/service"
	shopService "avito-crud/internal/service/shop"
	"avito-crud/internal/service/transfer"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type transaction struct {
	transferService service.ITransferService
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

func (t *transaction) sendCoin(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		ctx.JSON(400, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	token := tokenParts[1]

	var req SendCoinRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := t.transferService.SendCoin(ctx, token, req.ToUser, req.Amount)
	if err != nil {
		if errors.Is(err, shopService.ErrInvalidToken) {
			ctx.JSON(401, gin.H{"error": shopService.ErrInvalidToken.Error()})
			return
		}
		if errors.Is(err, shopRepo.ErrInsufficientFunds) {
			ctx.JSON(400, gin.H{"error": shopRepo.ErrInsufficientFunds.Error()})
			return
		}
		if errors.Is(err, transfer.ErrSameUser) {
			ctx.JSON(400, gin.H{"error": transfer.ErrSameUser.Error()})
			return
		}
		if errors.Is(err, authRepo.ErrUserNotFound) {
			ctx.JSON(404, gin.H{"error": authRepo.ErrUserNotFound.Error()})
			return
		}
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{})
}
