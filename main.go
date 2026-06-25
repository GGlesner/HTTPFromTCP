package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	filePath   = "./messages.txt"
	bufferSize = 8
)

func main() {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error openning file: %s. error: %s\n", filePath, err.Error())
	}
	writeOuput := func(s string) {
		fmt.Printf("read: %s\n", s)
	}

	for line := range parsLines(fileHandle) {
		writeOuput(line)
	}
}

func parsLines(file io.ReadCloser) <-chan string {
	lines := make(chan string)
	buffer := make([]byte, bufferSize)
	line := ""
	sendLines := func() {
		defer file.Close()
		defer close(lines)
		for {
			n, err := file.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return
				}
				fmt.Printf("error reading file: %s. error: %s\n", filePath, err.Error())
			}
			parts := strings.Split(string(buffer[:n]), "\n")
			m := len(parts)
			line += parts[0]
			parts[0] = line
			for i := range m - 1 {
				lines <- parts[i]
			}
			line = parts[m-1]
		}
	}
	go sendLines()
	return lines
}
