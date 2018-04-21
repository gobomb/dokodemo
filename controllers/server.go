package controllers

import (
	"doko/server"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
	"net/http"
	"doko/util"
)

var (
	StopChan   chan interface{}
	stopOkChan chan interface{}
	//startOkChan
)

func StartServer(context *gin.Context) {
	StopChan = util.NewChan()
	//startOkChan := util.NewChan()
	//stopOkChan = util.NewChan()
	go server.Main(StopChan)

	context.JSON(200, "start server success")
}

func StopServer(context *gin.Context) {
	StopChan <- 1
	//<-stopOkChan
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
