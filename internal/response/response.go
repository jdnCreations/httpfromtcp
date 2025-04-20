package response

import (
	"io"
	"strconv"

	"github.com/jdnCreations/httpfromtcp/internal/headers"
)

func (w *Writer) WriteTrailers(h headers.Headers) error {
	_, err := w.writer.Write([]byte("X-Content-Sha256" + ": " + h.Get("x-content-sha256") + "\r\n"))
	if err != nil {
		return err
	}

	_, err = w.writer.Write([]byte("X-Content-Length" + ": " + h.Get("x-content-length") + "\r\n"))
	if err != nil {
		return err
	}
	_, err = w.writer.Write([]byte("\r\n"))
	return err 
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