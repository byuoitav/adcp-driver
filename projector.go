package adcp

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/pooled"
)

// Projector is the base level object for a projector controlled over ADCP
type Projector struct {
	poolInit sync.Once
	pool     *pooled.Pool
	Address  string
}

const (
	// CR is a carriage return
	CR = '\r'
	// LF is a line feed
	LF = '\n'
)

func getConnection(key interface{}) (pooled.Conn, error) {
	address, ok := key.(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}

	conn, err := net.DialTimeout("tcp", address+":53595", 10*time.Second)
	if err != nil {
		return nil, err
	}

	// read the NOKEY line
	pconn := pooled.Wrap(conn)
	b, err := pconn.ReadUntil(LF, 5*time.Second)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if strings.TrimSpace(string(b)) != "NOKEY" {
		conn.Close()
		return nil, fmt.Errorf("unexpected message when opening connection: %s", b)
	}

	return pconn, nil
}

// SendCommand sends the byte array to the desired address of the projector
func (p *Projector) SendCommand(ctx context.Context, addr string, cmd []byte) (string, error) {
	p.poolInit.Do(func() {
		// create the pool
		p.pool = pooled.NewPool(45*time.Second, 400*time.Millisecond, getConnection)
	})

	var resp []byte
	err := p.pool.Do(addr, func(conn pooled.Conn) error {
		conn.SetWriteDeadline(time.Now().Add(3 * time.Second))

		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return err
		case n != len(cmd):
			return fmt.Errorf("wrote %v/%v bytes of command 0x%x", n, len(cmd), cmd)
		}

		resp = make([]byte, 5)

		resp, err = conn.ReadUntil(LF, 3*time.Second)
		if err != nil {
			return err
		}

		conn.Log().Debugf("Response from command: 0x%x", resp)

		return nil
	})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(resp)), nil
}
