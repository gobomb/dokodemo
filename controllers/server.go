package controllers

import (
	"doko/server"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
	"net/http"
	"doko/util"
)



func StartServer(context *gin.Context) {

	go server.Main(util.StopChan)

	context.JSON(200, "start server success")
}

func StopServer(context *gin.Context) {
	util.StopChan<-1
	context.JSON(200, "stop server success")
}

func Demo(context *gin.Context) {
	context.HTML(http.StatusOK, "demo.html", gin.H{})
}

func GetInfo(context *gin.Context) {
	info := server.GetInfo()
	log.Println(info)
	//log.Info(info.CtlReg)
	// 变量不可导出
	context.JSON(200, info)
}

func GetIndex(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{})
}
