package controllers

import (
	"doko/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartServer(context *gin.Context) {

	go server.Main()

	context.JSON(200, "start server success")
}

func GetInfo(context *gin.Context) {
	server.GetInfo()

}

func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
