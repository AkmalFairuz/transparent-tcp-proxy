package proxy

import (
	"github.com/sirupsen/logrus"
	"net"
	"sync/atomic"
)

type Connection struct {
	log        logrus.FieldLogger
	running    atomic.Bool
	conn       net.Conn
	serverConn net.Conn
}

func newConnection(log logrus.FieldLogger, conn net.Conn, serverConn net.Conn) *Connection {
	return &Connection{
		log:        log,
		conn:       conn,
		serverConn: serverConn,
	}
}
