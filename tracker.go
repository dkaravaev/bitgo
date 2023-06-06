package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/jackpal/bencode-go"
)

const Port uint16 = 20001

type TrackerInfo struct {
	Interval int
	Peers    []PeerInfo
}

type bencodeTrackerInfo struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func ParseTrackerResponse(resp *http.Response) (*bencodeTrackerInfo, error) {
	var buffer bytes.Buffer

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// bencodedBody := string(bytesBody)
	// fmt.Println(bencodedBody)

	buffer.Write(bytesBody)

	bti := bencodeTrackerInfo{}
	err = bencode.Unmarshal(&buffer, &bti)

	return &bti, err
}

func ParsePeers(peersStr string) ([]PeerInfo, error) {
	peersBin := []byte(peersStr)

	const peerSize = 6
	numPeers := len(peersBin) / peerSize
	if len(peersBin)%peerSize != 0 {
		return nil, fmt.Errorf("received malformed peers")
	}

	peers := make([]PeerInfo, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])
	}

	return peers, nil
}

func GetTrackerInfo(t *TorrentFile) (TrackerInfo, error) {
	var ti TrackerInfo

	urlString, _ := t.buildTrackerURL(MakeRandomPeerID(), Port)

	client := &http.Client{
		Timeout: 3*time.Second,
	}

	resp, err := client.Get(urlString)
	if err != nil {
		return ti, err
	}
	defer resp.Body.Close()

	bti, err := ParseTrackerResponse(resp)
	if err != nil {
		return ti, err
	}

	ti.Interval = bti.Interval
	ti.Peers, err = ParsePeers(bti.Peers)

	return ti, err
}
