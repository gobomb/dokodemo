package routers

import (
	"doko/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServerRouters(engine *gin.Engine) {
	routerGroup := engine.Group("/server")
	{
		// 静态资源文件路由
		routerGroup.Static("/js", "./front/js")
		routerGroup.Static("/css", "./front/css")
		routerGroup.Handle(http.MethodGet, "index", controllers.GetIndex)
		routerGroup.Handle(http.MethodGet, "start", controllers.StartServer)
		routerGroup.Handle(http.MethodGet, "info", controllers.GetInfo)
	}
}
