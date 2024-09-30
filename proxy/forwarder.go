package proxy

import (
	"fmt"
)

const (
	bufferSize = 65535
)

func startForwardPacket(c *Connection) error {
	go func() {
		for c.running.Load() {
			data := make([]byte, bufferSize)
			n, err := c.conn.Read(data)
			if err != nil {
				c.log.Errorf("an error occurred while reading data from client: %v", err)
				c.running.Store(false)
				return
			}

			_, err = c.serverConn.Write(data[:n])
			if err != nil {
				c.log.Errorf("an error occurred while writing data (client -> server): %v", err)
				c.running.Store(false)
				return
			}
		}
	}()

	for c.running.Load() {
		data := make([]byte, bufferSize)
		n, err := c.serverConn.Read(data)
		if err != nil {
			return fmt.Errorf("an error occurred while reading data from server: %v", err)
		}

		_, err = c.conn.Write(data[:n])
		if err != nil {
			return fmt.Errorf("an error occurred while writing data (server -> client): %v", err)
		}
	}

	c.running.Store(false)
	return nil
}
