package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"
)

var testFileName string = "data/debian-edu-11.5.0-amd64-netinst.iso.torrent"

func readFile(fileName string, t *testing.T) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		t.Errorf("File %s not found!\n", fileName)
	}

	return file
}

func TestParseTorrentFile(t *testing.T) {
	bto, err := ParseTorrentFile(readFile(testFileName, t))
	if err != nil {
		t.Errorf("%s", err)
	}

	announce := "http://bttracker.debian.org:6969/announce"
	if announce != bto.Announce {
		t.Errorf("Wrong announce! Expected: %s, got: %s\n", announce, bto.Announce)
	}

	bs := []byte(bto.Info.Pieces)
	piecesNum := (len(bs) / 20)
	totalSize := piecesNum * bto.Info.PieceLength

	totalSizeMB := int(float64(totalSize) * 1E-6)

	if totalSizeMB != 464 {
		t.Errorf("Wrong file size! Expected: 464 (MB), got: %d\n", totalSizeMB)
	}
}

func TestToTorrentFile(t *testing.T) {
	bto, err := ParseTorrentFile(readFile(testFileName, t))
	if err != nil {
		t.Errorf("%s", err)
	}

	_, err = bto.toTorrentFile()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestBuildTrackerURL(t *testing.T) {
	bto, err := ParseTorrentFile(readFile(testFileName, t))
	if err != nil {
		t.Errorf("%s", err)
	}

	tf, err := bto.toTorrentFile()
	if err != nil {
		t.Errorf("%s", err)
	}

	random := make([]byte, 20)
	rand.Read(random)

	var peerID [20]byte
	copy(peerID[:], random)

	url, _ := tf.buildTrackerURL(peerID, 10000)

	fmt.Println(url)
}
