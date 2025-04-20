package response

import (
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	length := len(p)

	hexHeader := fmt.Sprintf("%02x\r\n", length)
	chunk := append([]byte(hexHeader), p...)
	chunk = append(chunk, []byte("\r\n")...)

	return w.writer.Write(chunk)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	i, err := w.writer.Write([]byte("0\r\n"))
	return i, err
}