package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	const filePath = "messages.txt"
	fileHandle, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error openning file: %s. error: %s\n", filePath, err.Error())
	}
	defer fileHandle.Close()
	buffer := make([]byte, 8)
	for {
		n, err := fileHandle.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("error reading file: %s. error: %s\n", filePath, err.Error())
		}
		fmt.Printf("read: %s\n", string(buffer[:n]))
	}
}
