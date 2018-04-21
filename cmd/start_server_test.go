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
		stopChanC :=make(chan interface{})
		go client.Main(stopChanC)
		log.Info("client ok!")
		s.Done()
	}()
	stopChanC :=make(chan interface{})
	go client.Main(stopChanC)
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
		time.Sleep(2*time.Second)
		stopChanC :=util.NewChan()
		go client.Main(stopChanC)
		time.Sleep(10*time.Second)
		s.Done()
	}()
	stopChanC :=util.NewChan()
	go client.Main(stopChanC)
	server.Main(stopChan)

	log.Warn("Server stop ok!")
	s.Wait()

}
