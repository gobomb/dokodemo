package main

import (
	"doko/routers"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
)

var (
	serverWebAddress = "0.0.0.0:7777"
)

func startServerGin() {
	// 获取服务端 gin 实例
	sGin := gin.Default()

	// HTML 文件路由
	sGin.LoadHTMLGlob("./front/view/*")
	// 服务器路由
	routers.ServerRouters(sGin)
	// 启动 gin
	if err := sGin.Run(serverWebAddress); err != nil {
		log.Println(err)
	}
}

func main() {
	startServerGin()
}
