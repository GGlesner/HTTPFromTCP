package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const serverAddr = "localhost:42069"

func main() {
	udpAdrr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP address %s: %v\n", serverAddr, err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpAdrr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing UDP address %s: %v\n", serverAddr, err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl-c to exit.\n", serverAddr)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Message sent: %s", message)
	}
}
