package controllers

import (
	"doko/client"
	"github.com/gin-gonic/gin"
)

func ClientServer(context *gin.Context) {
	go client.Main()

	context.JSON(200, "start client success")
}
