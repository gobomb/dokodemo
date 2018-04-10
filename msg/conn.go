package msg

import (
	"doko/conn"
	"github.com/qiniu/log"
)

func ReadMsg(c conn.Conn) {
	var b []byte
	c.Read(b)
	log.Printf("[read from tcp]: %v", b)
}
