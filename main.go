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

func main() {
	file, err := os.Open("messages.txt")
	check(err)	

	// set up 8 bytes so we can read this amount at a time
	bytesToRead := make([]byte, 8)
	currentLine := ""

	for { 
		bytesRead, err := file.Read(bytesToRead) 
		if err == io.EOF {
			fmt.Printf("read: %s\n", currentLine)
			break
		}

		splitData := strings.Split(string(bytesToRead[:bytesRead]), "\n")
		for i, part := range splitData {
			if i < len(splitData) - 1 {
				if i == 0 {
					fmt.Printf("read: %s\n", currentLine + part)
					currentLine = ""
				} else {
					fmt.Printf("read: %s\n", part)
				}
			} else {
				currentLine += part
			}
		}
	}
}
