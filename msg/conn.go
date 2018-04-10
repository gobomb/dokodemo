package msg

import (
	"doko/conn"
	"github.com/qiniu/log"
)

func ReadMsg(c conn.Conn) {
	var b =make([]byte,6)
	n,err:=c.Read(b)
	log.Println(n)
	log.Println(err)
	log.Printf("[read from tcp]: %v", b)
}
