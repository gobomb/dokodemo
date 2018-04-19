package controllers

import (
	"doko/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/qiniu/log"
)

func StartServer(context *gin.Context) {

	go server.Main()

	context.JSON(200, "start server success")
}

func GetInfo(context *gin.Context) {
	info:=server.GetInfo()
	log.Println(info)
	log.Info(info.CtlReg)

}

func GetIndex(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{})
}
