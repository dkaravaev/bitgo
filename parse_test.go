package main

import (
	"os"
	"testing"
)

func TestParseTorrentFile(t *testing.T) {
	fileName := "data/debian-edu-11.5.0-amd64-netinst.iso.torrent"

	file, err := os.Open(fileName)
	if err != nil {
		t.Errorf("File %s not found!\n", fileName)
	}
	defer file.Close()

	bto, err := ParseTorrentFile(file)
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
