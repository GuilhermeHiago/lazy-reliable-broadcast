package main

import (
	"fmt"
	"bufio"
	//"time"
	"os"
	"strings"
	PFD "SD/PFD"
	BEB "SD/BEB"
)

type LazyReliableBroadcast_Req_Message struct {
	Addresses []string
	Message   string
}

type LazyReliableBroadcast_Ind_Message struct {
	From    string
	Message string
}

type LazyReliableBroadcast_Module struct {
	Ind               chan LazyReliableBroadcast_Ind_Message
	Req               chan LazyReliableBroadcast_Req_Message
	bestEffortBroadcast BEB.BestEffortBroadcast_Module
	failureDetector   PFD.PerfectFailureDetector_Module
	correctNodes      []string
	address			  string
	dbg               bool
}

func (module *LazyReliableBroadcast_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . [ LRB msg : " + s + " ]")
	}
}

func (module *LazyReliableBroadcast_Module) Init(address string, peers []string) {
	module.InitD(address, peers, true)
}

func (module *LazyReliableBroadcast_Module) InitD(address string, peers []string, _dbg bool) {
	module.dbg = _dbg
	module.outDbg("Init LRB!")
	module.address = address
	
	module.bestEffortBroadcast = BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}
	module.bestEffortBroadcast.Init(address, 2) /* colocar o failafter passando como parâmetro no terminal */
	
	
	module.failureDetector = PFD.PerfectFailureDetector_Module{
	}
	module.failureDetector.InitD(peers, _dbg, -1, address) // REMEBER TO USE FAILURE AFTER
	
	module.correctNodes = make([]string, len(peers))
	copy(module.correctNodes, peers)
	
	module.outDbg("peers: " + strings.Join(peers, ", "));
	module.outDbg("addr: " + address);
	module.Start()
}



func (module *LazyReliableBroadcast_Module) Start() {
	go func() {
		for {
			select {
			case failedNode := <-module.failureDetector.Fail:
				module.outDbg("added to suspect list: " + failedNode)
				module.correctNodes = removeNodeFromList(module.correctNodes, failedNode)
			}
		}
	}()

	go func() {
		for {
			select {
			case y := <-module.Req:
				y.Message = module.address + ";" + y.Message
				module.outDbg(y.Message)
				module.Broadcast(y)

			case y := <-module.bestEffortBroadcast.Ind:
				originAddr := strings.Split(y.Message, ";")[0]
				found := false
				for _, node := range module.correctNodes {
					if node == originAddr {
						found = true
						break
					}
				}
				
				if !found {
					module.outDbg("origin failed: " + originAddr + " relaying message")
					reqMessage := LazyReliableBroadcast_Req_Message{
						Addresses: module.correctNodes,
						Message:   module.address + ";" + strings.Split(y.Message, ";")[1]}
					module.Broadcast(reqMessage)
				}
				module.Deliver(BEB2LRB(y))
			}
		}
	}()
}


func (module *LazyReliableBroadcast_Module) Broadcast(message LazyReliableBroadcast_Req_Message) {
	msg := LRB2BEB(message)
	msg.Addresses = module.correctNodes //*
	module.bestEffortBroadcast.Req <- msg
	//module.outDbg("Sent to " + addr)
	
}

func (module *LazyReliableBroadcast_Module) Deliver(message LazyReliableBroadcast_Ind_Message) {
	module.outDbg("message delivered: " + message.Message)
	module.Ind <- message
}

func LRB2BEB(message LazyReliableBroadcast_Req_Message) BEB.BestEffortBroadcast_Req_Message {
	return BEB.BestEffortBroadcast_Req_Message{
		Addresses: message.Addresses,
		Message: message.Message}
}
		
func BEB2LRB(message BEB.BestEffortBroadcast_Ind_Message) LazyReliableBroadcast_Ind_Message {
	return LazyReliableBroadcast_Ind_Message{
		From: message.From,
		Message: message.Message}
}

func removeNodeFromList(nodes []string, nodeToRemove string) []string {
	index := -1
	for i, node := range nodes {
		if node == nodeToRemove {
			index = i
			break
		}
	}
	if index != -1 {
		return append(nodes[:index], nodes[index+1:]...)
	}
	return nodes
}

/*
	Victor: colocar passar failureAfter do BEB como argumento no terminal
			colocar uma lista das mensagens recebidas para não duplicar no Lazy
			remover a alteração da origem na mensagem quando fazer relay no Lazy
			alterar para usar heartbeat request no Failure Detector
*/
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run Chat.go <your_address> <peer_address1> <peer_address2> ...")
		return
	}

	address := os.Args[1]
	peers := os.Args[2:]

	chatModule := LazyReliableBroadcast_Module{
		Ind: make(chan LazyReliableBroadcast_Ind_Message),
		Req: make(chan LazyReliableBroadcast_Req_Message),
	}
	chatModule.InitD(address, peers, true)
	
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

		reqMsg := LazyReliableBroadcast_Req_Message{
			Addresses: peers,
			Message:   text,
		}

		chatModule.Req <- reqMsg
	}
}