package main

import(
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
	"doko/routers"
)

func main(){
	//common.StartUp()

	engine:=gin.Default()
	//engine.Use()

		routers.ServerRouters(engine)
	if err := engine.Run("0.0.0.0:7777"); err != nil {
		log.Println(err)
	}
	}