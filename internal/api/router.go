package api

import (
	"avito-crud/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(app *gin.Engine, authService service.IAuthService, shopService service.IShopService, transferService service.ITransferService) {
	api := app.Group("/api")
	newApiRoutes(api, authService, shopService, transferService)
}
func newApiRoutes(r *gin.RouterGroup, authService service.IAuthService, shopService service.IShopService, transferService service.ITransferService) {
	a := &auth{authService: authService}
	r.POST("/auth", a.auth)
	s := &shop{shopService: shopService}
	r.GET("/buy/:item", s.buyItem)
	t := &transaction{transferService: transferService}
	r.POST("/sendCoin", t.sendCoin)

}
