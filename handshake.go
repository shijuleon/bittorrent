package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type HandshakeRequest struct {
	ProtocolLen    uint8
	ProtocolString [19]byte
	Reserved       [8]byte
	InfoHash       [20]byte
	PeerID         [20]byte
}

type HandshakeResponse struct {
	Protocol [20]byte
	Reserved [8]byte
	Infohash [20]byte
	PeerID   [20]byte
}

const (
	ProtocolString = "BitTorrent protocol"
)

func Handshake(address string, infoHash string, peerID string) (ok bool, err error) {
	infoHashBytes, err := hex.DecodeString(infoHash)
	if err != nil {
		log.Fatalf("Error: Decoding hex from InfoHash: %s", err)
	}

	peerIDBytes, err := hex.DecodeString(peerID)
	if err != nil {
		log.Fatalf("Error: Decoding hex from PeerID: %s", err)
	}

	packet := &HandshakeRequest{
		ProtocolLen: uint8(19),
		Reserved:    [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}

	copy(packet.ProtocolString[:], ProtocolString)
	copy(packet.InfoHash[:], infoHashBytes)
	copy(packet.PeerID[:], peerIDBytes)

	response, err := SendPacket(address, packet)
	if err != nil {
		log.Errorf("sending packet: %s", err)
	}

	handshakeResponse, err := serializeHandshakeResponse(response)
	if err != nil {
		log.Errorf("serializing handshake response: %s", err)
	}
	if !bytes.Equal(handshakeResponse.Protocol[:], append([]byte{0x13}, packet.ProtocolString[:]...)) {
		return false, fmt.Errorf("unexpected handshake response: %s", handshakeResponse.Protocol)
	}

	return true, nil
}

func serializeHandshakeResponse(response []byte) (*HandshakeResponse, error) {
	r := bytes.NewReader(response)
	handshakeResponse := &HandshakeResponse{}

	if err := binary.Read(r, binary.LittleEndian, handshakeResponse); err != nil {
		log.Errorf("binary read failed: %s", err)
		return nil, err
	}

	log.Debugf("pstr: %s\n", string(handshakeResponse.Protocol[:]))
	log.Debugf("reserved: %x\n", handshakeResponse.Reserved[:])
	log.Debugf("info_hash: %x\n", handshakeResponse.Infohash[:])
	log.Debugf("peer_id: %x\n", handshakeResponse.PeerID[:])

	return handshakeResponse, nil
}
