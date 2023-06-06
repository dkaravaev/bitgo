package main

import (
	"crypto/rand"
	"fmt"
	"io"
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

func PeerHandshakeSize() int {
	return 1 + len("BitTorrent protocol") + 8 + 20 + 20
}

// Serialize serializes the handshake to a buffer
func (ph *PeerHandshake) Serialize() []byte {
	ph.Pstr = "BitTorrent protocol"
	buf := make([]byte, PeerHandshakeSize())

	offset := 0
	buf[offset] = byte(len(ph.Pstr))
	offset += 1
	offset += copy(buf[offset:], ph.Pstr)
	offset += copy(buf[offset:], make([]byte, 8))
	offset += copy(buf[offset:], ph.InfoHash[:])
	offset += copy(buf[offset:], ph.PID[:])

	return buf
}

func ReadHandshake(r io.Reader) (*PeerHandshake, error) {
	var ph PeerHandshake
	buf := make([]byte, PeerHandshakeSize())

	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	if n != PeerHandshakeSize() {
		return nil, fmt.Errorf("invalid peer handshake size %d", n)
	}

	offset := 0

	pLen := int(buf[offset])
	offset += 1

	ph.Pstr = string(buf[offset : offset+pLen])
	offset += pLen

	offset += copy(ph.InfoHash[:], buf[offset:offset+len(ph.InfoHash)])
	offset += copy(ph.PID[:], buf[offset:offset+len(ph.PID)])

	return &ph, nil
}

func (peer PeerInfo) String() string {
	return peer.IP.String() + ":" + strconv.Itoa(int(peer.Port))
}

func (pconn *PeerConn) ConnectToPeer(peer PeerInfo) (err error) {
	pconn.conn, err = net.DialTimeout("tcp", peer.String(), 3*time.Second)
	return err
}
