// Construido como parte da disciplina de Sistemas Distribuidos
// PUCRS - Escola Politecnica
// Professor: Fernando Dotti  (www.inf.pucrs.br/~fldotti)

/*
LANCAR N PROCESSOS EM SHELL's DIFERENTES, PARA CADA PROCESSO, O SEU PROPRIO ENDERECO EE O PRIMEIRO DA LISTA
go run chat.go 127.0.0.1:5001  127.0.0.1:6001    ...
go run chat.go 127.0.0.1:6001  127.0.0.1:5001    ...
go run chat.go ...  127.0.0.1:6001  127.0.0.1:5001
*/

package main

import (
	. "SD/PFD"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run pfd-test.go -1 127.0.0.1:5001 127.0.0.1:6001 127.0.0.1:7001")
		fmt.Println("go run pfd-test.go -1 127.0.0.1:6001 127.0.0.1:5001 127.0.0.1:7001")
		fmt.Println("go run pfd-test.go 5 127.0.0.1:7001 127.0.0.1:6001  127.0.0.1:5001")
		return
	}

	var registro []string
	failureAfter, _ := strconv.Atoi(os.Args[1])
	addresses := os.Args[2:]
	fmt.Println(addresses)

	pfd := &PerfectFailureDetector_Module{}
	var onCrashP = make(chan string)

	pfd.InitD(addresses, true, failureAfter, os.Args[2], onCrashP)

	time.Sleep(30 * time.Second)

	// avisa quando um processo P fizer Crash
	go func() {
		for {
			in := <-onCrashP
			registro = append(registro, in+":failed ")

			// imprime a mensagem recebida na tela
			fmt.Printf("          Fail from %v\n", in)
		}
	}()

	blq := make(chan int)
	<-blq
}
