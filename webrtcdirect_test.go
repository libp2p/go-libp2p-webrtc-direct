package libp2pwebrtcdirect

import (
	"testing"

	logging "github.com/ipfs/go-log"

	utils "github.com/libp2p/go-libp2p-transport/test"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pions/webrtc"
	mplex "github.com/whyrusleeping/go-smux-multiplex"
)

func TestTransport(t *testing.T) {
	logging.SetLogLevel("*", "warning")

	ta := NewTransport(
		webrtc.Configuration{},
		new(mplex.Transport),
	)
	tb := NewTransport(
		webrtc.Configuration{},
		new(mplex.Transport),
	)

	addr := "/ip4/127.0.0.1/tcp/0/http/p2p-webrtc-direct"
	utils.SubtestTransport(t, ta, tb, addr, "peerA")
}

func TestTransportCantListenUtp(t *testing.T) {
	utpa, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/50000")
	if err != nil {
		t.Fatal(err)
	}

	tpt := NewTransport(
		webrtc.Configuration{},
		new(mplex.Transport),
	)

	_, err = tpt.Listen(utpa)
	if err == nil {
		t.Fatal("shouldnt be able to listen on utp addr with tcp transport")
	}

}
