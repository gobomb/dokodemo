package cmd

import (
	"testing"
	"time"
	"github.com/qiniu/log"
	"doko/util"
	"doko/server"
	"doko/client"
	"sync"
)

// 测试 关闭服务器功能
func TestRunServer(t *testing.T) {
	s := sync.WaitGroup{}
	stopChan := util.NewChan()
	go func() {
		s.Add(1)
		time.Sleep(10 * time.Second)
		stopChan <- 1
		time.Sleep(2 * time.Second)
		go client.Main()
		log.Info("client ok!")
		s.Done()
	}()
	go client.Main()
	server.Main(stopChan)

	log.Warn("Server stop ok!")
	s.Wait()

}

func TestRerunServer(t *testing.T) {
	s := sync.WaitGroup{}
	stopChan := util.NewChan()
	go func() {
		s.Add(1)
		time.Sleep(10 * time.Second)
		stopChan <- 1
		time.Sleep(2 * time.Second)
		stopChan = util.NewChan()
		go server.Main(stopChan)
		log.Info("server rerun ok!")
		s.Done()
	}()
	go client.Main()
	server.Main(stopChan)

	log.Warn("Server stop ok!")
	s.Wait()

}
