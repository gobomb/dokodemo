package routers

import (
	"doko/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServerRouters(engine *gin.Engine) {
    // 页面路由
	engine.Handle(http.MethodGet, "/", controllers.GetIndex)
	engine.Handle(http.MethodGet,"ports",controllers.GetPortPage)
	engine.Handle(http.MethodGet, "index", controllers.GetIndex)
	engine.Handle(http.MethodGet, "clients", controllers.GetClientPage)

	// 接口路由
	engine.Handle(http.MethodGet, "status-on", controllers.StartServer)
	engine.Handle(http.MethodGet, "status-off", controllers.StopServer)
	engine.Handle(http.MethodGet, "info", controllers.GetInfo)
	engine.Handle(http.MethodGet,"gotty",controllers.Gotty)

}
