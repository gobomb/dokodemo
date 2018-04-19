package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"doko/controllers"
)

func ClientRouters(engine *gin.Engine) {
	routerGroup := engine.Group("/client")
	{
		// 静态资源文件路由
		//routerGroup.Static("/js", "./front/js")
		//routerGroup.Static("/css", "./front/css")
		//routerGroup.Handle(http.MethodGet, "index", controllers.GetIndex)
		routerGroup.Handle(http.MethodGet, "start", controllers.ClientServer)
		//routerGroup.Handle(http.MethodGet, "info", controllers.GetInfo)
	}
}
