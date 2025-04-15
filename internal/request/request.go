package request

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/jdnCreations/httpfromtcp/internal/headers"
)

type requestState int

const (
	requestStateInitialized requestState = iota // 0
	requestStateDone // 1
	requestStateParsingHeaders
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers 
	state requestState // to track parsing state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, bytesConsumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		
		if bytesConsumed == 0 {
			return 0, nil
		}
	
		// successfully parsed request line
		r.RequestLine = requestLine
		r.state = requestStateParsingHeaders 
		return bytesConsumed, nil

	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateDone
		}
		return n, nil
	
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	
	default:
		return 0, errors.New("error: unknown state")
	}
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		totalBytesParsed += n

		if n == 0 {
			break
		}
	}

	return totalBytesParsed, nil
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
	request := &Request{
		state: requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for request.state != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.state, n)
				}
				break
			}
			return nil, err
		}
		readToIndex += n

		parsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[parsed:readToIndex])
		readToIndex -= parsed
		
	}
	return request, nil
}