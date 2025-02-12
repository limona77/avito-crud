package api

import (
	authRepo "avito-crud/internal/repostiory/auth"
	shopRepo "avito-crud/internal/repostiory/shop"
	"avito-crud/internal/service"
	shopService "avito-crud/internal/service/shop"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

type shop struct {
	shopService service.IShopService
	log         *slog.Logger
}

func (s *shop) buyItem(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Проверяем, что токен начинается с "Bearer "
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(400, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	token := tokenParts[1]

	item := c.Param("item")
	if item == "" {
		c.JSON(400, gin.H{"error": "item not provided"})
		return
	}

	err := s.shopService.BuyItem(c, token, item)
	if err != nil {
		if errors.Is(err, shopRepo.ErrInsufficientFunds) {
			c.JSON(400, gin.H{"error": shopRepo.ErrInsufficientFunds.Error()})
			return
		}
		if errors.Is(err, shopRepo.ErrMerchNotFound) {
			c.JSON(404, gin.H{"error": shopRepo.ErrMerchNotFound.Error()})
			return
		}
		if errors.Is(err, authRepo.ErrUserNotFound) {
			c.JSON(404, gin.H{"error": authRepo.ErrUserNotFound.Error()})
			return
		}
		if errors.Is(err, shopService.ErrInvalidToken) {
			c.JSON(400, gin.H{"error": shopService.ErrInvalidToken.Error()})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "item bought"})
}
