package libp2pwebrtcdirect

import (
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
)

func init() {
	if err := ma.AddProtocol(Protocol); err != nil {
		panic(fmt.Errorf("error registering p2p-webrtc protocol: %s", err))
	}
}

const protoCode = 0x1992

var Protocol = ma.Protocol{
	Name:       "p2p-webrtc",
	Code:       protoCode,
	VCode:      ma.CodeToVarint(protoCode),
	Size:       -1,
	Transcoder: Transcoder,
}

var Transcoder ma.Transcoder = transcoder{}

// transcoder handles encoding/decoding peer ids as base58.
type transcoder struct {
	strtobyte func(string) ([]byte, error)
	bytetostr func([]byte) (string, error)
	validbyte func([]byte) error
}

func (t transcoder) StringToBytes(s string) ([]byte, error) {
	id, err := peer.IDB58Decode(s)
	if err != nil {
		return nil, err
	}
	return []byte(id), nil
}

func (t transcoder) BytesToString(b []byte) (string, error) {
	id, err := peer.IDFromBytes(b)
	if err != nil {
		return "", err
	}
	return peer.IDB58Encode(id), nil
}

func (t transcoder) ValidateBytes(b []byte) error {
	_, err := peer.IDFromBytes(b)
	return err
}
