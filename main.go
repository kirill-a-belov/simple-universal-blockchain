package main

import (
	"./blockchainer"
	"./listener"
	"fmt"
)

func main() {
	fmt.Println("Starting Simple universal blockchain server!")
	fmt.Println()

	//Open channel for workers
	bcChan := make(chan string)
	go blockchainer.BlockchainWorker(bcChan)

	listener.Listen(bcChan)
}
