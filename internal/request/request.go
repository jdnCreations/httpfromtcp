package request

import (
	"errors"
	"io"
	"strings"
)

const (
	stateInitalized = iota // 0
	stateDone // 1
)

type Request struct {
	RequestLine RequestLine
	state int // to track parsing state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateInitalized:
		requestLine, bytesConsumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		
		if bytesConsumed == 0 {
			return 0, nil
		}

		// successfully parsed request line
		r.RequestLine = requestLine
		r.state = stateDone
		return bytesConsumed, nil

	case stateDone:
		return 0, errors.New("error: trying to read data in a done state")

	default:
		return 0, errors.New("error: unknown state")
	}
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	dataToParse := string(data)

	crlfIndex := strings.Index(dataToParse, "\r\n")
	if crlfIndex == -1 {
		return RequestLine{}, 0, nil
	}

	// extract request line without \r\n
	requestLineStr := dataToParse[:crlfIndex]
	// calc bytes consumed incl \r\n
	bytesConsumed := crlfIndex + 2

	parts := strings.Split(requestLineStr, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("invalid request line format")
	}

	method := parts[0]
	if method != "GET" && method != "POST" {
		return RequestLine{}, 0, errors.New("invalid request method")
	}

	target := parts[1]

	httpVersion := parts[2]
	if !strings.HasPrefix(httpVersion, "HTTP/") {
		return RequestLine{}, 0, errors.New("invalid HTTP version")
	}

	versionParts := strings.Split(httpVersion, "/")
	if len(versionParts) != 2 {
		return RequestLine{}, 0, errors.New("invalid HTTP version format")
	}

	requestLine := RequestLine{
		Method: method,
		RequestTarget: target,
		HttpVersion: versionParts[1],
	}

	return requestLine, bytesConsumed, nil
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	request := Request{
		state: stateInitalized,
	}

	for request.state != stateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if readToIndex > 0 {
					parsed, parseErr := request.parse(buf[:readToIndex])
					if parseErr != nil {
						return nil, parseErr
					}

					if parsed > 0 {
						copy(buf, buf[parsed:readToIndex])
						readToIndex -= parsed
					}
				}
				request.state = stateDone
				break
			}
			return nil, err
		}

		readToIndex += n
		parsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if parsed > 0 {
			copy(buf, buf[parsed:readToIndex])
			readToIndex -= parsed
		}
		
	}

	return &request, nil
}