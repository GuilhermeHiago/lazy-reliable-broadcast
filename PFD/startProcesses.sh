#! /bin/bash

go run PerfectFailureDetector.go 4 127.0.0.1:8050 127.0.0.1:8051 127.0.0.1:8052 >> log1.txt & 
go run PerfectFailureDetector.go -1 127.0.0.1:8051 127.0.0.1:8050 127.0.0.1:8052 >> log2.txt &
go run PerfectFailureDetector.go -1 127.0.0.1:8052 127.0.0.1:8051 127.0.0.1:8050 >> log3.txt &
