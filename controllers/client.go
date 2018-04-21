package controllers

import (
	"doko/client"
	"github.com/gin-gonic/gin"
	"doko/util"
)

var stopChan chan interface{}

func ClientServer(context *gin.Context) {

	stopChan = util.NewChan()
	go client.Main(stopChan)

	context.JSON(200, "start client success")
}

func StopClient(context *gin.Context) {
	stopChan<-1
	context.JSON(200, "start server success")
}
