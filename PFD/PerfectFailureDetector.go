/*
  Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
  Autores: Guilherme H C Santos & Victor Putrich
  Modulo representando Perfect Failure Detector tal como definido em:
    Introduction to Reliable and Secure Distributed Programming
    Christian Cachin, Rachid Gerraoui, Luis Rodrigues
  * Semestre 2023/1 - Primeira versao.
*/

package PerfectFailureDetector

/*
	Victor: alterar perfect failure para usar heartbeat request 
*/

import (
	"fmt"
	"os"
	"time"
	//"strconv"
	"strings"
	PP2PLink "SD/PP2PLink"
)

type PerfectFailureDetector_Module struct {
	//Modules
	p2pLink       PP2PLink.PP2PLink
	
	// args
	alive         map[string]bool
	detected      map[string]bool
	peers         []string
	timeout       chan int
	
	// implementation args
	timeoutCounter map[string]int
	failureAfter  int
	heartbeat 	  chan int
	address		  string
	Fail          chan string
	dbg           bool
}
func (module *PerfectFailureDetector_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ PFD msg : " + s + " ]")
		os.Stdout.Sync()
	}
}

func (module *PerfectFailureDetector_Module) InitD(peers []string, _dbg bool, failureAfter int, address string) {
	
	
	// Append port to the address
	module.address = address + ":8000"
	
	module.peers = make([]string, len(peers)) // Fix this line to create a slice of strings
	// Append port to each address in peers
	for i, peer := range peers {
		module.peers[i] = peer + ":8000" 
	}
	
	module.dbg = _dbg
	//module.peers = peers
	module.alive = make(map[string]bool)
	module.detected = make(map[string]bool)
	module.timeoutCounter = make(map[string]int) // Add timeoutCounter map
	module.timeout = make(chan int)
	module.heartbeat = make(chan int)
	module.failureAfter = failureAfter
	module.Fail = make(chan string, len(peers))
	for _, peer := range peers {
		module.alive[peer] = true
		module.detected[peer] = false
		module.timeoutCounter[peer] = 0 // Initialize timeoutCounter for each peer
	}

	module.outDbg("Init PFD!")

	module.p2pLink = PP2PLink.PP2PLink{
		Req: make(chan PP2PLink.PP2PLink_Req_Message),
		Ind: make(chan PP2PLink.PP2PLink_Ind_Message)}
	module.p2pLink.Init(module.address)

	module.outDbg("peers: " + strings.Join(peers, ", "));
	module.outDbg("addr: " + address);
	module.Start()
}

func (module *PerfectFailureDetector_Module) StartTimer() {
	time.Sleep(300 * time.Millisecond)
	module.timeout <- 1
}

func (module *PerfectFailureDetector_Module) heartbeatTimer() {
	time.Sleep(50 * time.Millisecond)
	module.heartbeat <- 1
}

func (module *PerfectFailureDetector_Module) Start() {
	go module.heartbeatTimer()
	go module.StartTimer()

	go func() {
		for {
			select {
			case <-module.timeout:
				if module.failureAfter != -1 {
					module.failureAfter -= 1
					if module.failureAfter == 0 {
						module.outDbg("Process failed")
						return
					}
				}

				for _, peer := range module.peers {
					if !module.alive[peer] && !module.detected[peer] {
						module.timeoutCounter[peer]++
						if module.timeoutCounter[peer] >= 2 {
							peerAddress := strings.Split(peer, ":")[0]
							module.detected[peer] = true
							module.outDbg(fmt.Sprintf("Process %s failed", peerAddress))
							module.Fail <- peerAddress // Notify failure
						}
					}
				}
				//module.printPeers()
				module.setAllNodesToFalse()
				go module.StartTimer()
			
			case <-module.heartbeat:
				module.sendHeartbeatMessages()
				go module.heartbeatTimer()
			
			case indMsg := <-module.p2pLink.Ind:
				module.receiveHeartbeatMessage(indMsg)
			}
		}
	}()
}

func (module *PerfectFailureDetector_Module) setAllNodesToFalse() {
	for peer := range module.alive {
		module.alive[peer] = false
	}
}

func (module *PerfectFailureDetector_Module) printPeers() {
	fmt.Println("List of peers:")
	for _, peer := range module.peers {
		fmt.Println(peer, module.alive[peer], module.detected[peer])
	}

	os.Stdout.Sync()
}

func (module *PerfectFailureDetector_Module) sendHeartbeatMessages() {
    for _, peer := range module.peers {
        if !module.detected[peer] {
            msg := PP2PLink.PP2PLink_Req_Message{
                To: peer,
                Message: "heartbeat;" + module.address + ";" + peer,
            }
            module.p2pLink.Req <- msg
        }
    }
}
	
func (module *PerfectFailureDetector_Module) receiveHeartbeatMessage(msg PP2PLink.PP2PLink_Ind_Message) {
    msgParts := strings.Split(msg.Message, ";")
    if len(msgParts) < 3 || msgParts[0] != "heartbeat" {
        fmt.Println("Error: Received message is not in the expected format")
        return
    }
	
    from := msgParts[1]
	if from != module.address{
    	//module.outDbg(from)
	}
	//module.outDbg(from)
	module.alive[from] = true
}
	
// func main() {
// 	if len(os.Args) < 3 {
// 		fmt.Println("Usage:   go run PerfectFailureDetector.go  failureAfter thisProcessIpAddress:port otherProcessIpAddress:port")
// 		fmt.Println("Example: go run PerfectFailureDetector.go -1 127.0.0.1:8050 127.0.0.1:8051")
// 		fmt.Println("Example: go run PerfectFailureDetector.go 2 127.0.0.1:8051 127.0.0.1:8050")
// 		return
// 	}

// 	failureAfter, err := strconv.Atoi(os.Args[1])
// 	if err != nil {
//     	fmt.Println("Invalid failureAfter argument")
//     	return
// 	}

// 	peers := os.Args[2:] // Replace these with actual process IP:Port addresses
// 	pfd := &PerfectFailureDetector_Module{}



// 	pfd.InitD(peers, true, failureAfter, os.Args[2])

// 	time.Sleep(60 * time.Second)
// }