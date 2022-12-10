package main

import (
	"bytes"
	"crypto/sha1"
	"io"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type bencodeInfo struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func (bto bencodeTorrent) computeInfoHash() [20]byte {
	var buffer bytes.Buffer
	bencode.Marshal(&buffer, bto.Info)

	return sha1.Sum(buffer.Bytes())
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	var ret TorrentFile

	ret.Announce = bto.Announce
	ret.PieceLength = bto.Info.PieceLength
	ret.Length = bto.Info.Length
	ret.Name = bto.Info.Name

	bs := []byte(bto.Info.Pieces)
	SHALength := 20

	ret.PieceHashes = make([][20]byte, len(bs)/SHALength)
	for i := 0; i < len(ret.PieceHashes); i++ {
		copy(ret.PieceHashes[i][:], bs[i*SHALength:(i+1)*SHALength])
	}

	ret.InfoHash = bto.computeInfoHash()

	return ret, nil
}

func ParseTorrentFile(r io.Reader) (*bencodeTorrent, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, nil
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, Port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
