package main

import (
	"bufio"
	"bytes"
	"encoding/binary"

	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

func SendPacket(address string, packet interface{}) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		log.Errorf("dialing: %s", err)
	}

	buffer := new(bytes.Buffer)
	err = binary.Write(buffer, binary.LittleEndian, packet)
	if err != nil {
		log.Errorf("writing packet to binary buffer: %s", err)
		return nil, err
	}

	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		log.Errorf("sending packet: %s", err)
		return nil, err
	}

	b, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Errorf("reading bytes from response: %s", err)
		return nil, err
	}

	return b, nil
}
