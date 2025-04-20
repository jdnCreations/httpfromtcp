package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jdnCreations/httpfromtcp/internal/headers"
	"github.com/jdnCreations/httpfromtcp/internal/request"
	"github.com/jdnCreations/httpfromtcp/internal/response"
	"github.com/jdnCreations/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w)	
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w)	
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, 
	"/httpbin") {
		handlerhttpbin(w, req)	
		return
	}
	handler200(w)
}

func handlerhttpbin(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusOK)

	reqUrl := fmt.Sprintf("https://httpbin.org/%s", strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin"))
	resp, err := http.Get(reqUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.Add("Transfer-Encoding", "chunked")
	h.Add("Trailer", "X-Content-Sha256, X-Content-Length")
	w.WriteHeaders(h)

	var responseBody []byte

	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading response body:", err)
			return
		}
		chunk := buffer[:n]
		responseBody = append(responseBody, chunk...)
		w.WriteChunkedBody(chunk)
	}
	_, err = w.WriteChunkedBodyDone()
	sha := sha256.Sum256(responseBody)
	t := headers.NewHeaders()

	t.Add("X-Content-Sha256", fmt.Sprintf("%x", sha))
	t.Add("X-Content-Length", fmt.Sprintf("%d", len(responseBody)))
	w.WriteTrailers(t)

	if err != nil {
		fmt.Printf("Error writing final chunk: %v\n", err)
	}
}

func handler400(w *response.Writer) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`
	<html>
	<head>
	<title>400 Bad Request</title>
	</head>
	<body>
	<h1>Bad Request</h1>
	<p>Your request honestly kinda sucked.</p>
	</body>
	</html>`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer) {
	w.WriteStatusLine(response.StatusInternalServerError)
	body := []byte(`
	<html>
	<head>
	<title>500 Internal Server Error</title>
	</head>
	<body>
	<h1>Internal Server Error</h1>
	<p>Okay, you know what? This one is on me.</p>
	</body>
	</html>`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(`
	<html>
	<head>
	<title>200 OK</title>
	</head>
	<body>
	<h1>Success!</h1>
	<p>Your request was an absolute banger.</p>
	</body>
	</html>`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}