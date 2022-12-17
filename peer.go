package main

import (
	"crypto/rand"
	"net"
	"strconv"
	"time"
)

type PeerID = [20]byte

type PeerInfo struct {
	IP   net.IP
	Port uint16
}

type PeerHandshake struct {
	Pstr     string
	InfoHash [20]byte
	PID      PeerID
}

type PeerConn struct {
	conn net.Conn
}

func MakeRandomPeerID() PeerID {
	var ret [20]byte

	random := make([]byte, 20)
	rand.Read(random)

	copy(ret[:], random)

	return ret
}

// Serialize serializes the handshake to a buffer
func (ph *PeerHandshake) Serialize() []byte {
	size := 1 + len(ph.Pstr) + 8 + len(ph.InfoHash) + len(ph.PID)
	buf := make([]byte, size)

	offset := 0
	buf[offset] = byte(len(ph.Pstr))
	offset += 1
	offset += copy(buf[offset:], ph.Pstr)
	offset += copy(buf[offset:], make([]byte, 8))
	offset += copy(buf[offset:], ph.InfoHash[:])
	offset += copy(buf[offset:], ph.PID[:])

	return buf
}

func (peer PeerInfo) String() string {
	return peer.IP.String() + ":" + strconv.Itoa(int(peer.Port))
}

func (pconn *PeerConn) ConnectToPeer(peer PeerInfo) (err error) {
	pconn.conn, err = net.DialTimeout("tcp", peer.String(), 3*time.Second)
	return err
}
