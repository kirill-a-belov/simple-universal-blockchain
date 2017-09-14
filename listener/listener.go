package listener

import (
	"bufio"
	"fmt"
	"net"
)

//This func open socket and listen port for messages
func Listen(bcChan chan string) {

	listenAddr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345}
	//Open listen port
	listner, err := net.ListenTCP("tcp4", listenAddr)
	if err != nil {
		fmt.Println("Error while openig listen port: ", listenAddr.String())
		panic(err)
	}
	//Accepting remote connection
	conn, err := listner.AcceptTCP()
	if err != nil {
		fmt.Println("Error while openig connection")
		panic(err)
	} else {
		fmt.Println("New connection from: ", conn.RemoteAddr().String())
	}

	//Make connection data buffer
	connBuf := bufio.NewReader(conn)

	//Listener cycle
	for {
		msg, err := connBuf.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading message from connection: ", err)

		} else {
			fmt.Print("Recived message: ", msg)
			//Sending new message to blockchain worker
			bcChan <- msg
		}
	}
}
