package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jdnCreations/httpfromtcp/internal/request"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	check(err)
	defer listener.Close()
  for {
    conn, err := listener.Accept()
    check(err)
    
    req, err := request.RequestFromReader(conn)
    if err != nil {
			log.Fatal(err.Error())
		}	

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)

    conn.Close()
  }
}
