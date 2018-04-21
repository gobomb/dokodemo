package controllers

import (
	"testing"
	"doko/util"
	"doko/server"
	"github.com/qiniu/log"
	"time"
)

// 测试 状态返回
func TestStartServer(t *testing.T) {
	StopChan = util.NewChan()
	go server.Main(StopChan)
	time.Sleep(4 * time.Second)
	info := server.GetInfo()
	log.Info(info)
	time.Sleep(10 * time.Second)
	StopChan <- 1
	time.Sleep(6 * time.Second)
	info = server.GetInfo()
	log.Info(info)
}
func TestStopServer(t *testing.T) {
	StopChan <- 1
}
func TestGetInfo(t *testing.T) {
	info := server.GetInfo()
	log.Info(info)
}
