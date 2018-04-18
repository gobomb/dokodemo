package conn

import (
	"net"
	"github.com/qiniu/log"
	"fmt"
	"sync"
	"io"
	"math/rand"
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
	log.Debugf("Renamed connection %s", oldId)
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
		log.Debugf("[net.Listen tcp failed:]%v", err)
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
				log.Debugf("Failed to accept new TCP connection : %v", err)
				continue
			}
			c := &loggedConn{
				tcp:  rawConn.(*net.TCPConn),
				Conn: rawConn,
			}

			log.Debugf("New connection from %v", rawConn.RemoteAddr())
			l.Conns <- c
		}
	}()
	return
}

func Dial(addr,typ string) (*loggedConn) {
	var rawConn net.Conn
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Debugf("[net.Dial error]: %v", err)
		return nil
	}

	conn := &loggedConn{
		tcp:  rawConn.(*net.TCPConn),
		Conn: rawConn,
		typ:  typ,
	}
	//.tcp=rawConn.(*net.TCPConn)

	log.Debugf("New connection to: %v", rawConn.RemoteAddr())

	return conn
}

func Join(c Conn, c2 Conn) (int64, int64) {
	var wait sync.WaitGroup
	pipe := func(to Conn, from Conn, bytesCopied *int64) {
		defer to.Close()
		defer from.Close()
		defer wait.Done()

		var err error
		*bytesCopied, err = io.Copy(to, from)
		if err != nil {
			log.Debugf("Copied %d bytes to %s before failing with error %v", *bytesCopied, to.Id(), err)
		} else {
			log.Debugf("Copied %d bytes to %s", *bytesCopied, to.Id())
		}
	}

	wait.Add(2)
	var fromBytes, toBytes int64
	go pipe(c, c2, &fromBytes)
	go pipe(c2, c, &toBytes)
	log.Debugf("oined with connection %s\n", c2.Id())
	wait.Wait()
	return fromBytes, toBytes
}

func Wrap(conn net.Conn, typ string) *loggedConn{
	return wrapConn(conn, typ)
}

func wrapConn(conn net.Conn,typ string)*loggedConn{
	switch c:= conn.(type){
	case *loggedConn:
		return c
	case *net.TCPConn:
		wrapped := &loggedConn{c,conn,rand.Int31(),typ}
		return wrapped
	}
	return nil
}
