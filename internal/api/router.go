package api

import (
	"avito-crud/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(app *gin.Engine, authService service.IAuthService) {
	auth := app.Group("/api")
	newAuthRoutes(auth, authService)
}
