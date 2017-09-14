package main

import (
	"fmt"
	"net"
)

func main() {
	//Open connection to a server
	connAddr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345}
	conn, err := net.DialTCP("tcp4", nil, connAddr)
	if err != nil {
		fmt.Println("Error while connecting to: ", connAddr.String())
		panic(err)
	}

	//Sending cycle
	msg := ""
	for {
		fmt.Println("Enter your message to save in blockchain:")
		fmt.Scanf("%v", &msg)

		conn.Write([]byte(msg + "\n"))
		msg = ""
	}
}
