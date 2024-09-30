package proxy

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"sync/atomic"
)

type Proxy struct {
	log           logrus.FieldLogger
	listen        net.Listener
	targetAddress string
	running       atomic.Bool

	connections   map[net.Addr]*Connection
	connectionsMu sync.RWMutex
}

func New(log logrus.FieldLogger, listenAddress, targetAddress string) (*Proxy, error) {
	listen, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		log:           log,
		connections:   map[net.Addr]*Connection{},
		listen:        listen,
		targetAddress: targetAddress,
	}, nil
}

func (p *Proxy) Run() error {
	p.running.Store(true)
	for p.running.Load() {
		conn, err := p.listen.Accept()
		if err != nil {
			return err
		}

		go func() {
			if err := p.handleNewConn(conn); err != nil {
				p.log.Errorf("an error occurred while handling connection: %v", err)
			}
		}()
	}
	return nil
}

func (p *Proxy) handleNewConn(conn net.Conn) error {
	defer conn.Close()

	target, err := net.Dial("tcp", p.targetAddress)
	if err != nil {
		return err
	}

	defer target.Close()

	p.log.Infof("new connection from %s to %s", conn.RemoteAddr(), target.RemoteAddr())

	connection := newConnection(p.log, conn, target)
	connection.running.Store(true)

	p.connectionsMu.RLock()
	if _, ok := p.connections[conn.RemoteAddr()]; ok {
		p.connectionsMu.RUnlock()
		return fmt.Errorf("connection from %s already exists", conn.RemoteAddr().String())
	}
	p.connectionsMu.RUnlock()

	p.connectionsMu.Lock()
	p.connections[conn.RemoteAddr()] = connection
	p.connectionsMu.Unlock()

	defer func() {
		p.log.Infof("connection from %s to %s closed", conn.RemoteAddr(), target.RemoteAddr())

		p.connectionsMu.Lock()
		delete(p.connections, conn.RemoteAddr())
		p.connectionsMu.Unlock()
	}()

	return startForwardPacket(connection)
}
