package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server!")

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	fmt.Println("Message from server:", message)
}
