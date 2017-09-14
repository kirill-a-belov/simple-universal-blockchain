package blockchainer

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

const MSG_IN_BLOCK = 5

type BlockID struct {
	sync.RWMutex
	ID int
}

type Block struct {
	BlockId  int
	Time     time.Time
	Sign     []byte
	Messages []string
}

var blockID BlockID
var blockchain = make([]Block, 0)
var nodePrivatekey *ecdsa.PrivateKey

//This func collect messages and making blocks
func BlockchainWorker(in chan string) {
	//Get block id
	blockID.Lock()
	newBlockID := blockID.ID + 1
	blockID.ID = newBlockID
	blockID.Unlock()

	messages := []string{}
	fmt.Println("Started another blockchainer worker thread")
	//Main cycle
	for {
		msg := <-in
		if msg == "!Print_blockchain!\n" {
			printBlockchain()
		} else {
			messages = append(messages, msg)
		}

		//We have enough messages, lets pack it in block
		if len(messages) >= MSG_IN_BLOCK {
			newBlockTime := time.Now()

			newBlock := Block{
				BlockId:  newBlockID,
				Time:     newBlockTime,
				Messages: messages,
			}

			//Calculating hash and signing sign
			sign := []byte{}

			//1 - Make bytes from block ID
			binPos := make([]byte, 4)
			binary.LittleEndian.PutUint32(binPos, uint32(newBlockID))
			sign = append(sign, binPos...)

			//2 - Get previous block and hash it
			if len(blockchain)-1 > newBlockID {
				var binPrev bytes.Buffer
				binary.Write(&binPrev, binary.BigEndian, blockchain[newBlockID-1])
				prevBlockHash := sha256.Sum256([]byte(binPrev.Bytes()))
				sign = append(sign, prevBlockHash[:]...)
			}

			//3 - Hash creating time
			binTime, _ := newBlockTime.MarshalBinary()
			sign = append(sign, binTime...)

			//4 - Making Merkel tree and get it hash
			merkelHash := GetMerkelHash(messages)
			sign = append(sign, merkelHash[:]...)
			signHash := sha256.Sum256(sign)

			// 5 - Sign block sign by private key
			r, s, err := ecdsa.Sign(rand.Reader, nodePrivatekey, signHash[:])
			if err != nil {
				panic(err)
			}
			newBlock.Sign = r.Bytes()
			newBlock.Sign = append(newBlock.Sign, s.Bytes()...)

			//New block ready - save it into memory (file, db, cashe, etc.)
			blockchain = append(blockchain, newBlock)
			fmt.Println("*********************************")
			fmt.Println(" Creating and saved new block:")
			fmt.Println(" Block ID: ", newBlock.BlockId)
			fmt.Println(" Block creation time: ", newBlock.Time)
			fmt.Printf(" Block sign: %x\n", newBlock.Sign)
			fmt.Println(" Block messages:")
			for _, msg := range newBlock.Messages {
				fmt.Print("  ", msg)
			}
			fmt.Println("*********************************")

			//Clear storages
			blockID.Lock()
			newBlockID = blockID.ID + 1
			blockID.ID = newBlockID
			blockID.Unlock()

			messages = []string{}
		}
	}
}

func init() {
	blockID.ID = 0

	var err error
	nodePrivatekey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		if err != nil {
			fmt.Println("Error while generating nodes private key")
			panic(err)
		}
	}
}

func printBlockchain() {
	fmt.Println("Result blockchain is:")
	for _, block := range blockchain {
		fmt.Println("*********************************")
		fmt.Println(" Block ID: ", block.BlockId)
		fmt.Println(" Block creation time: ", block.Time)
		fmt.Printf(" Block sign: %x\n", block.Sign)
		fmt.Println(" Block messages:")
		for _, msg := range block.Messages {
			fmt.Print("  ", msg)
		}
	}
}
