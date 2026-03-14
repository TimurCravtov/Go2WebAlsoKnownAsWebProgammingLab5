package connect

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"time"
)

func RequestTCPRaw(host string, port int, withTls bool, content []byte, readStrategy func(net.Conn) ([]byte, error)) ([]byte, error) {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	
	var conn net.Conn
	var err error

	if withTls {
		conf := &tls.Config{InsecureSkipVerify: false, ServerName: host}
		conn, err = tls.Dial("tcp", address, conf)
	} else {
		conn, err = net.DialTimeout("tcp", address, 5*time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if _, err = conn.Write(content); err != nil {
		return nil, fmt.Errorf("write error: %w", err)
	}

	return readStrategy(conn)
}
