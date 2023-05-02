#!/bin/bash

gnome-terminal -- bash -c "go run LazyReliableBroadcast.go -1 127.0.0.1 127.0.0.2 127.0.0.3 127.0.0.4" &
gnome-terminal -- bash -c "go run LazyReliableBroadcast.go -1 127.0.0.2 127.0.0.1 127.0.0.3 127.0.0.4" &
gnome-terminal -- bash -c "go run LazyReliableBroadcast.go -1 127.0.0.3 127.0.0.1 127.0.0.2 127.0.0.4" &
gnome-terminal -- bash -c "go run LazyReliableBroadcast.go 5 127.0.0.4 127.0.0.1 127.0.0.2 127.0.0.3; exec bash"

