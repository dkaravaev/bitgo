package main

import (
	"crypto/sha1"
	"io"

	"github.com/jackpal/bencode-go"
)

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

func ParseTorrentFile(r io.Reader) (*bencodeTorrent, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, nil
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

	// todo: Add bencoded info block SHA-1 calculcation
	hasher := sha1.New()
	hasher.Write(bs)

	return ret, nil
}
