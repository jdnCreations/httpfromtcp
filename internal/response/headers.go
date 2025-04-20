package response

import (
	"strconv"

	"github.com/jdnCreations/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Parse([]byte("Content-Length: " + strconv.Itoa(contentLen) + "\r\n"))
	h.Parse([]byte("Connection: close\r\n"))
	h.Parse([]byte("Content-Type: text/plain\r\n"))
	return h
}