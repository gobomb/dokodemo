package main

import(
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
)

func main(){
	//common.StartUp()

	engine:=gin.Default()
	//engine.Use()

	if err := engine.Run("0.0.0.0:7777"); err != nil {
		log.Println(err)
	}
	}