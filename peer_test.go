package main

import (
	"net"
	"testing"
)

func TestPeerString(t *testing.T) {
	peer := PeerInfo{net.IPv4(0x0, 0x1, 0x2, 0x3), 256}

	if peer.String() != "0.1.2.3:256" {
		t.Errorf("Wrong parsing")
	}
}
