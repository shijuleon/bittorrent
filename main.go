package main

import (
	log "github.com/sirupsen/logrus"
)

const (
	peerAddress = "localhost:30613"
)

func main() {
	log.SetLevel(log.DebugLevel)
	if ok, err := Handshake(peerAddress, "5a8ce26e8a19a877d8ccc927fcc18e34e1f5ff67", "4a5ce26f8a13a877d8ccc987fcc18e24e1f5ff37"); !ok {
		log.Fatalf("handshake failed: %s", err)
	}
}
