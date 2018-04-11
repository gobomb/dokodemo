package conn

import (
	"net"
	"github.com/qiniu/log"
	"fmt"
)

type Conn interface {
	net.Conn
	Id() string
	SetType(string)
	CloseRead() error
}

type loggedConn struct {
	tcp *net.TCPConn
	net.Conn
	id  int32
	typ string
}

type Listener struct {
	net.Addr
	Conns chan *loggedConn
}

func (c *loggedConn) Id() string {
	return fmt.Sprintf("%s:%x", c.typ, c.id)
}

func (c *loggedConn) SetType(typ string) {
	oldId := c.Id()
	c.typ = typ
	log.Printf("Renamed connection %s", oldId)
}

func (c *loggedConn) CloseRead() error {
	// XXX: use CloseRead() in Conn.Join() and in Control.shutdown() for cleaner
	// connection termination. Unfortunately, when I've tried that, I've observed
	// failures where the connection was closed *before* flushing its write buffer,
	// set with SetLinger() set properly (which it is by default).
	return c.tcp.CloseRead()
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

func Dial(addr string) ( *loggedConn) {
	var rawConn net.Conn
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("[net.Dial error]: %v", err)
		return nil
	}

	conn:=&loggedConn{
		tcp:rawConn.(*net.TCPConn),
		Conn:rawConn,
		typ:"auth",
	}
	//.tcp=rawConn.(*net.TCPConn)

	log.Printf("New connection to: %v", rawConn.RemoteAddr())

	return conn
}
