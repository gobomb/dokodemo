package routers

import (
	"doko/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServerRouters(engine *gin.Engine) {
	routerGroup := engine.Group("/server")
	{
		routerGroup.Handle(http.MethodGet,"demo",controllers.Demo)
		routerGroup.Handle(http.MethodGet, "index", controllers.GetIndex)
		routerGroup.Handle(http.MethodGet, "status", controllers.StartServer)
		routerGroup.Handle(http.MethodGet, "info", controllers.GetInfo)
	}
}
