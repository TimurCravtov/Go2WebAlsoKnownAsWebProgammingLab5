package connect

import (
	"net"
	"time"
	"fmt"
	"crypto/tls"
	"strconv"
	"io"
)

func RequestTCPRaw(host string, port int, withTls bool, content []byte) ([]byte, error) {

	address := net.JoinHostPort(host, strconv.Itoa(port))
	
	var conn net.Conn
	var err error

	// configure tls if needed and establish the connection
	if withTls {
		conf := &tls.Config{
			ServerName: host,
		}
		conn, err = tls.Dial("tcp", address, conf)
	} else {
		conn, err = net.DialTimeout("tcp", address, 5*time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	// send the payload
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Write(content)
	if err != nil {
		return nil, fmt.Errorf("write error: %w", err)
	}

	// read the payload
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	response, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	return response, nil
}