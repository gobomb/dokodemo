package controllers

import (
	"github.com/gin-gonic/gin"
	"doko/server"
)

func StartServer(context *gin.Context) {

	go server.Main()

	context.JSON(200, "start server success")
}
