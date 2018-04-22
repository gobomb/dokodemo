package controllers

import (
	"doko/server"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
	"net/http"
	"doko/util"
)

var (
	StopChan chan interface{}
	onceChan chan interface{}
	//startOkChan
)

func init() {
	onceChan = util.NewChan(1)
}

func StartServer(context *gin.Context) {
	select {
	case onceChan <- 1:
		StopChan = util.NewChan(0)
		go server.Main(StopChan)
		context.JSON(200, "start server success")
	default:
		context.JSON(200, "the server has started")
	}
}

func StopServer(context *gin.Context) {
	select {
	case <-onceChan:
		StopChan <- 1
		context.JSON(200, "stop server success")
	default:
		context.JSON(200, "the server has stopped")
	}

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
