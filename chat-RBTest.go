package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	LRB "SD/RB"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run Chat.go <your_address> <peer_address1> <peer_address2> ...")
		return
	}

	address := os.Args[1]
	peers := os.Args[2:]

	chatModule := LRB.LazyReliableBroadcast_Module{
		Ind: make(chan LRB.LazyReliableBroadcast_Ind_Message),
		Req: make(chan LRB.LazyReliableBroadcast_Req_Message),
	}
	chatModule.InitD(address, true)

	go func() {
		for {
			indMsg := <-chatModule.Ind
			fmt.Printf("%s: %s\n", indMsg.From, indMsg.Message)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")

		reqMsg := LRB.LazyReliableBroadcast_Req_Message{
			Addresses: peers,
			Message:   text,
		}

		chatModule.Req <- reqMsg
	}
}
