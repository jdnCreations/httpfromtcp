package response

import (
	"io"
	"log"
	"strconv"

	"github.com/jdnCreations/httpfromtcp/internal/headers"
)

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for header := range headers {
		value := headers.Get(header)
		log.Printf("Header: %s, Value: %s", header, value)
		_, err := w.Write([]byte(header + ": " + headers.Get(header) + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}



func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
	case StatusInternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
	default:
		_, err := w.Write([]byte("HTTP/1.1 " + strconv.Itoa(int(statusCode)) + " \r\n"))
		if err != nil {
			return err
		}
	}
	return nil
}