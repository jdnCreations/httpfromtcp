package main

import (
	"fmt"
	"io"
	"log"
	"os"
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
	file, err := os.Open("messages.txt")
	check(err)

	channel := getLinesChannel(file)
	for v := range channel {
		fmt.Printf("read: %s\n", v)
	}	
}
