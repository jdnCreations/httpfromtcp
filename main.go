package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineCh := make(chan string)
	go func() {
		defer f.Close()
		defer close(lineCh)
		bytesToRead := make([]byte, 8)
		currentLine := ""
		
		for { 
			bytesRead, err := f.Read(bytesToRead) 
			if err == io.EOF {
				lineCh <- currentLine
				break	
			}
		
			splitData := strings.Split(string(bytesToRead[:bytesRead]), "\n")
			for i, part := range splitData {
				if i < len(splitData) - 1 {
					if i == 0 {
						lineCh <- currentLine + part
						currentLine = ""
					} else {
						lineCh <- part
					}
				} else {
					currentLine += part
				}
			}
		}
	}()
	return lineCh
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	check(err)
	defer listener.Close()
  for {
    conn, err := listener.Accept()
    check(err)
    
    lines := getLinesChannel(conn)
    fmt.Println("connection accepted")
    for line := range lines {
      fmt.Printf("%s\n", line)
    }
    
    conn.Close()
    fmt.Println("Connection closed")
  }
}
