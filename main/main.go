package main

import (
	//"github.com/gin-gonic/gin"
	//"github.com/qiniu/log"
	//"doko/routers"
	"doko/server"
)

func main() {

	//engine := gin.Default()
	//
	//routers.ServerRouters(engine)
	//if err := engine.Run("0.0.0.0:7777"); err != nil {
	//	log.Println(err)
	//}
	server.Main()
}
