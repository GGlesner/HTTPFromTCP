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
	defer fileHandle.Close()
	buffer := make([]byte, bufferSize)
	line := ""
	writeOuput := func(s string) {
		fmt.Printf("read: %s\n", s)
	}
	for {
		n, err := fileHandle.Read(buffer)
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
		for i := 1; i < m; i++ {
			writeOuput(parts[i-1])
		}
		line = parts[m-1]
	}
}
