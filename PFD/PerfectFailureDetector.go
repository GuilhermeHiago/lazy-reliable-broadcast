/*
  Construido como parte da disciplina: Sistemas Distribuidos - PUCRS - Escola Politecnica
  Autores: Guilherme H C Santos & Victor Putrich
  Modulo representando Perfect Failure Detector tal como definido em:
    Introduction to Reliable and Secure Distributed Programming
    Christian Cachin, Rachid Gerraoui, Luis Rodrigues
  * Semestre 2023/1 - Primeira versao.
*/
package PerfectFailureDetector

import (
	"fmt"
	"os"
	"time"
	"strings"
	PP2PLink "SD/PP2PLink"
)

type PerfectFailureDetector_Module struct {
	alive         map[string]bool
	detected      map[string]bool
	timeoutCounter map[string]int
	peers         []string
	dbg           bool
	timeout       chan int
	failureAfter  int
	address		  string
	p2pLink       PP2PLink.PP2PLink
	onCrashP	  chan string
}

func (module *PerfectFailureDetector_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . [ PFD msg : " + s + " ]")
		os.Stdout.Sync()
	}
}

func (module *PerfectFailureDetector_Module) Init(peers []string, failureAfter int, address string, onCrashP chan string) {
	module.InitD(peers, true, failureAfter, address, onCrashP)
}

func (module *PerfectFailureDetector_Module) InitD(peers []string, _dbg bool, failureAfter int, address string, onCrashP chan string) {
	module.dbg = _dbg
	module.peers = peers
	module.alive = make(map[string]bool)
	module.detected = make(map[string]bool)
	module.timeoutCounter = make(map[string]int) // Add timeoutCounter map
	module.timeout = make(chan int)
	module.failureAfter = failureAfter
	module.address = address
	module.onCrashP = onCrashP

	for _, peer := range peers {
		module.alive[peer] = true
		module.detected[peer] = false
		module.timeoutCounter[peer] = 0 // Initialize timeoutCounter for each peer
	}

	module.outDbg("Init PFD!")

	module.p2pLink = PP2PLink.PP2PLink{
		Req: make(chan PP2PLink.PP2PLink_Req_Message),
		Ind: make(chan PP2PLink.PP2PLink_Ind_Message),
	}

	module.p2pLink.InitD(address, _dbg)
	module.Start()
}

func (module *PerfectFailureDetector_Module) StartTimer() {
	time.Sleep(1 * time.Second)
	module.timeout <- 1
}

func (module *PerfectFailureDetector_Module) Start() {
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

				//module.outDbg("trigger heard")
				for _, peer := range module.peers {
					if !module.alive[peer] && !module.detected[peer] {
						module.timeoutCounter[peer]++ // Increment timeoutCounter if the process is not alive and not detected
						if module.timeoutCounter[peer] >= 5 { // Check if the counter reached the threshold (5)
							module.detected[peer] = true
							module.outDbg(fmt.Sprintf("Process %s failed", peer))
							module.onCrashP <- peer
						}
					}
				}

				module.sendHeartbeatMessages()

				//module.printPeers()
				module.setAllNodesToFalse()
				go module.StartTimer()

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
    	module.outDbg(from)
	}
	//module.outDbg(from)
	module.alive[from] = true
}
	
// func main() {
// 	if len(os.Args) < 3 {
// 		fmt.Println("Usage:   go run PerfectFailureDetector.go  failureAfter thisProcessIpAddress:port otherProcessIpAddress:port")
// 		fmt.Println("go run PerfectFailureDetector.go -1 127.0.0.1:5001 127.0.0.1:6001 127.0.0.1:7001")
// 		fmt.Println("go run PerfectFailureDetector.go -1 127.0.0.1:6001 127.0.0.1:5001 127.0.0.1:7001")
// 		fmt.Println("go run PerfectFailureDetector.go 5 127.0.0.1:7001 127.0.0.1:6001  127.0.0.1:5001")
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

// 	time.Sleep(30 * time.Second)
// }