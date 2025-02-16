package api

import (
	authRepo "avito-crud/internal/repostiory/auth"
	shopRepo "avito-crud/internal/repostiory/shop"
	"avito-crud/internal/service"
	shopService "avito-crud/internal/service/shop"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type shop struct {
	shopService service.IShopService
}

func (s *shop) buyItem(ctx *gin.Context) {
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

	item := ctx.Param("item")
	if item == "" {
		ctx.JSON(400, gin.H{"error": "item not provided"})
		return
	}

	err := s.shopService.BuyItem(ctx, token, item)
	if err != nil {
		if errors.Is(err, shopRepo.ErrInsufficientFunds) {
			ctx.JSON(400, gin.H{"error": shopRepo.ErrInsufficientFunds.Error()})
			return
		}
		if errors.Is(err, shopRepo.ErrMerchNotFound) {
			ctx.JSON(404, gin.H{"error": shopRepo.ErrMerchNotFound.Error()})
			return
		}
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
	ctx.JSON(200, gin.H{"message": "item bought"})
}
