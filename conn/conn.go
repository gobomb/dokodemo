package conn

import (
	"net"
	"github.com/qiniu/log"
)

type Conn interface {
	net.Conn
	Id() string
	//SetType(string)
	//CloseRead() error
}

type loggedConn struct {
	tcp *net.TCPConn
	net.Conn
	//id  int32
	//typ string
}

type Listener struct {
	net.Addr
	Conns chan *loggedConn
}

func Listen(addr string) (l *Listener) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("[net.Listen tcp failed:]%v", err)
		return
	}
	l = &Listener{
		Addr:  listener.Addr(),
		Conns: make(chan *loggedConn),
	}

	go func() {
		for {
			rawConn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept new TCP connection : %v", err)
				continue
			}
			c := &loggedConn{
				tcp:  rawConn.(*net.TCPConn),
				Conn: rawConn,
			}

			log.Printf("New connection from %v", rawConn.RemoteAddr())
			l.Conns <- c
		}
	}()
	return
}

func Dial(addr string) (rawConn net.Conn) {
	//var rawConn net.Conn
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("[net.Dial error]: %v", err)
		return
	}

	log.Printf("New connection to: %v", rawConn.RemoteAddr())

	return
}
