package main

import (
	"crypto/rand"
	"net"
)

type PeerID = [20]byte

func MakeRandomPeerID() PeerID {
	var ret [20]byte

	random := make([]byte, 20)
	rand.Read(random)

	copy(ret[:], random)

	return ret
}

type Peer struct {
	IP   net.IP
	Port uint16
}
