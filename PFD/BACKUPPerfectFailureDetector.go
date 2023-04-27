package main

import (
	"fmt"
	"time"
)

const (
	INITIAL_DELTA = 1000 * time.Millisecond
	MAX_DELTA     = 32000 * time.Millisecond
)

type PerfectFailureDetector_Ind_Message struct {
	From    string
	Message string
}

type PerfectFailureDetector_Module struct {
	aliveNodes       map[string]bool
	heartbeatTimeout map[string]time.Time
	deltaTimeout     map[string]time.Duration
	ind              chan PerfectFailureDetector_Ind_Message
	heartbeatChan chan PerfectFailureDetector_Heartbeat_Message
	dbg              bool
}

type PerfectFailureDetector_Heartbeat_Message struct {
	From string
}

func (module *PerfectFailureDetector_Module) Init(aliveNodes []string) {
	module.InitD(aliveNodes, true)
}

func (module *PerfectFailureDetector_Module) InitD(aliveNodes []string, _dbg bool) {
	module.dbg = _dbg
	module.outDbg("Init PFD!")
	module.aliveNodes = make(map[string]bool)
	module.heartbeatTimeout = make(map[string]time.Time)
	module.deltaTimeout = make(map[string]time.Duration)
	module.heartbeatChan = make(chan PerfectFailureDetector_Heartbeat_Message)
	
	for _, node := range aliveNodes {
		module.aliveNodes[node] = true
		module.heartbeatTimeout[node] = time.Now().Add(INITIAL_DELTA)
		module.deltaTimeout[node] = INITIAL_DELTA
	}
	module.ind = make(chan PerfectFailureDetector_Ind_Message)
	module.Start()
}

func (module *PerfectFailureDetector_Module) Start() {
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			now := time.Now()

			select {
			case heartbeat := <-module.heartbeatChan:
				module.aliveNodes[heartbeat.From] = true
			default:
				for node, timeout := range module.heartbeatTimeout {
					if now.After(timeout) {
						if module.aliveNodes[node] && !module.detectedNodes[node] {
							module.detectedNodes[node] = true
							module.ind <- PerfectFailureDetector_Ind_Message{From: node, Message: "CRASHED"}
						}
						module.sendHeartbeatRequest(node)
						newDelta := module.deltaTimeout[node]
						if newDelta > MAX_DELTA {
							newDelta = MAX_DELTA
						}
						module.deltaTimeout[node] = newDelta
						module.heartbeatTimeout[node] = time.Now().Add(newDelta)
					}
				}
				module.aliveNodes = make(map[string]bool)
				module.restartTimerWithDelta()
			}
		}
	}()
}

func (module *PerfectFailureDetector_Module) sendHeartbeatRequest(node string) {
	// Send a heartbeat request to the given node
}

func (module *PerfectFailureDetector_Module) restartTimerWithDelta() {
	// Restart the timer with delta
}

func (module *PerfectFailureDetector_Module) Heartbeat(from string) {
	if _, ok := module.heartbeatTimeout[from]; ok {
		module.heartbeatTimeout[from] = time.Now().Add(module.deltaTimeout[from])
		module.deltaTimeout[from] = INITIAL_DELTA
		if !module.aliveNodes[from] {
			module.aliveNodes[from] = true
			module.ind <- PerfectFailureDetector_Ind_Message{From: from, Message: "RECOVERED"}
		}
	}
}

func (module *PerfectFailureDetector_Module) outDbg(format string, a ...interface{}) {
	if module.dbg {
		fmt.Printf(format+"\n", a...)
	}
}

func main() {
	addresses := []string{"node1", "node2", "node3"}
	pfd := &PerfectFailureDetector_Module{}
	pfd.InitD(addresses, true)

	go func() {
		for indMsg := range pfd.ind {
			fmt.Printf("Node %s: %s\n", indMsg.From, indMsg.Message)
		}
	}()

	time.Sleep(2 * time.Second)
	pfd.Heartbeat("node1")
	time.Sleep(2 * time.Second)
	pfd.Heartbeat("node2")
	time.Sleep(2 * time.Second)
	pfd.Heartbeat("node3")

	time.Sleep(5 * time.Second)
}
