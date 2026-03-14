package connect

import (
	"fmt"
	"strconv"
	"time"
	"net"
	"bufio"
	"bytes"
	"io"
	"strings"
)

type HttpResponse struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

// core method of this module
func Request(method string, url string, body []byte, headers map[string]string) (*HttpResponse, error) {
    host, port, secured := parseURL(url)

    httpRequest := []byte(fmt.Sprintf("%s / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n", method, host))
    for key, value := range headers {
        httpRequest = append(httpRequest, []byte(fmt.Sprintf("%s: %s\r\n", key, value))...)
    }
    httpRequest = append(httpRequest, []byte("\r\n")...)
    
    if len(body) > 0 {
        httpRequest = append(httpRequest, body...)
    }

    rawBytes, err := RequestTCPRaw(host, port, secured, httpRequest, HttpTCPReadStrategy)
    if err != nil {
        return nil, err
    }

    return ParseRawToResponse(rawBytes)
}

func Get(url string, body []byte, headers map[string]string) (*HttpResponse, error) {
	return Request("GET", url, body, headers)
}

func Post(url string, body []byte, headers map[string]string) (*HttpResponse, error) {
	return Request("POST", url, body, headers)
}

func Delete(url string, body []byte, headers map[string]string) (*HttpResponse, error) {
	return Request("DELETE", url, body, headers)
}

func Put(url string, body []byte, headers map[string]string) (*HttpResponse, error) {
	return Request("PUT", url, body, headers)
}

func parseURL(url string) (string, int, bool) {

	var host string
	var port int
	var secured bool

	if len(url) > 7 && url[:7] == "http://" {
		host = url[7:]
		secured = false
		port = 80
	} else if len(url) > 8 && url[:8] == "https://" {
		host = url[8:]
		port = 443
		secured = true
	} else {
		host = url
		port = 80
		secured = false;
	}
	return host, port, secured
}

func HttpTCPReadStrategy(conn net.Conn) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	
	reader := bufio.NewReader(conn)
	var fullResponse bytes.Buffer
	var isChunked bool
	contentLength := -1

	// read headers
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading headers: %w", err)
		}
		fullResponse.Write(line)

		if bytes.Equal(line, []byte("\r\n")) {
			break
		}

		// content-length or chunked?
		lowerLine := strings.ToLower(string(line))
		if strings.HasPrefix(lowerLine, "content-length:") {
			parts := strings.Split(lowerLine, ":")
			if len(parts) > 1 {
				fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &contentLength)
			}
		} else if strings.Contains(lowerLine, "transfer-encoding: chunked") {
			isChunked = true
		}
	}

	// read based on strategy
	if isChunked {
		body, err := readChunkedBody(reader)
		if err != nil {
			return nil, err
		}
		fullResponse.Write(body)
	} else if contentLength > 0 {
		body := make([]byte, contentLength)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			return nil, fmt.Errorf("error reading fixed body: %w", err)
		}
		fullResponse.Write(body)
	} else {
		// neither: read all
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		body, _ := io.ReadAll(reader) 
		fullResponse.Write(body)
	}

	return fullResponse.Bytes(), nil
}

func readChunkedBody(reader *bufio.Reader) ([]byte, error) {
	var body bytes.Buffer
	for {

		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		parts := strings.Split(strings.TrimSpace(line), ";")
		size, err := strconv.ParseInt(parts[0], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid chunk size: %w", err)
		}

		if size == 0 {
			reader.Discard(2) 
			break
		}

		chunk := make([]byte, size)
		_, err = io.ReadFull(reader, chunk)
		if err != nil {
			return nil, err
		}
		body.Write(chunk)

		reader.Discard(2)
	}
	return body.Bytes(), nil
}


func ParseRawToResponse(raw []byte) (*HttpResponse, error) {
    reader := bufio.NewReader(bytes.NewReader(raw))
    resp := &HttpResponse{
        Headers: make(map[string]string),
    }

    statusLine, err := reader.ReadString('\n')
    if err != nil {
        return nil, fmt.Errorf("failed to read status line: %w", err)
    }
    statusParts := strings.SplitN(strings.TrimSpace(statusLine), " ", 3)
    if len(statusParts) >= 2 {
        resp.StatusCode, _ = strconv.Atoi(statusParts[1])
    }
    if len(statusParts) == 3 {
        resp.StatusText = statusParts[2]
    }

    for {
        line, err := reader.ReadString('\n')
        if err != nil || line == "\r\n" || line == "\n" {
            break
        }
        parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
        if len(parts) == 2 {
            resp.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
        }
    }

    body, err := io.ReadAll(reader)
    if err != nil {
        return nil, fmt.Errorf("failed to read body: %w", err)
    }
    resp.Body = body

    return resp, nil
}