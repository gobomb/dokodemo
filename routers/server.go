package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"doko/controllers"
)

func ServerRouters(engine *gin.Engine){
	routerGroup := engine.Group("/server")
	{
		//start:=routerGroup.Group("/start")
		routerGroup.Handle(http.MethodPost,"start",controllers.StartServer )
	}
}
