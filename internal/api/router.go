package api

import (
	"avito-crud/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	app *gin.Engine,
	authService service.IAuthService,
	shopService service.IShopService,
	transferService service.ITransferService,
	infoService service.IInfoService,
) {
	api := app.Group("/api")
	newApiRoutes(
		api,
		authService,
		shopService,
		transferService,
		infoService,
	)
}

func newApiRoutes(r *gin.RouterGroup, authService service.IAuthService, shopService service.IShopService, transferService service.ITransferService, infoService service.IInfoService) {
	a := &auth{authService: authService}
	r.POST("/auth", a.auth)
	s := &shop{shopService: shopService}
	r.GET("/buy/:item", s.buyItem)
	t := &transaction{transferService: transferService}
	r.POST("/sendCoin", t.sendCoin)
	i := &info{infoService: infoService}
	r.GET("/info", i.getInfo)
}
