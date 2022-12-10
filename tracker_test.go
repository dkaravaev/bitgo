package main

import (
	"testing"
)

func TestGetTrackerInfo(t *testing.T) {
	bto, err := ParseTorrentFile(ReadFile(TestFileName, t))
	if err != nil {
		t.Errorf("%s", err)
	}

	tf, err := bto.toTorrentFile()
	if err != nil {
		t.Errorf("%s", err)
	}

	_, err = GetTrackerInfo(&tf)
	if err != nil {
		t.Errorf("%s", err)
	}
	// fmt.Println(ti.Peers)
}
