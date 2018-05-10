package controllers

import (
	"doko/server"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/log"
	"net/http"
	"doko/util"
	"os"
	"fmt"
	"time"
)

var (
	StopChan chan interface{}
	onceChan chan interface{}
	//startOkChan
)

func init() {
	procMap = make(map[int]*os.Process)
	onceChan = util.NewChan(1)
}

func StartServer(context *gin.Context) {
	select {
	case onceChan <- 1:
		StopChan = util.NewChan(0)
		go server.Main(StopChan)
		context.JSON(200, true)
	default:
		context.JSON(200,false)
	}
}

func StopServer(context *gin.Context) {
	select {
	case <-onceChan:
		StopChan <- 1
		context.JSON(200, true)
	default:
		context.JSON(200, false)
	}

}

func GetPortPage(context *gin.Context) {
	context.HTML(http.StatusOK, "port.html", gin.H{})
}


func GetClientPage(context *gin.Context) {
	context.HTML(http.StatusOK, "client.html", gin.H{})
}


func GetInfo(context *gin.Context) {
	info := server.GetInfo()
	log.Println(info)
	//log.Info(info.CtlReg)
	// 变量不可导出
	context.JSON(200, info)
}

func GetIndex(context *gin.Context) {
	context.HTML(http.StatusOK, "client.html", gin.H{})
}

func Gotty(context *gin.Context) {
	pid := execShell()
	context.JSON(200, pid)
}

func StopGotty(context *gin.Context) {
	var pid int
	context.Bind(pid)
	stopShell(pid)
	context.JSON(200, pid)
}

var procMap map[int]*os.Process

func stopShell(pid int) {
	process := procMap[pid]
	time.Sleep(10 * time.Second)
	if err := process.Kill(); err != nil {
		log.Error(err)
	}
	log.Print("kill success")
}

func execShell() int {
	// 1) os.StartProcess //
	/*********************/
	/* Linux: */
	env := os.Environ()
	procAttr := &os.ProcAttr{
		Env: env,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}
	// 1st example: list files
	process, err := os.StartProcess("/Users/cym/go/bin/gotty", []string{"", "-c", "a:b", "zsh"}, procAttr)
	if err != nil {
		fmt.Printf("Error %v starting process!", err) //
		os.Exit(1)
	}
	log.Printf("The process id is %v", process)
	procMap[process.Pid] = process
	return process.Pid

}
