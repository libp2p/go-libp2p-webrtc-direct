package libp2pwebrtcdirect

import (
	"github.com/pion/webrtc/v3"
	"testing"

	logging "github.com/ipfs/go-log"

	mplex "github.com/libp2p/go-libp2p-mplex"
	utils "github.com/libp2p/go-libp2p-testing/suites/transport"
	ma "github.com/multiformats/go-multiaddr"
)

func TestTransport(t *testing.T) {
	t.Skip("This test is failing, see https://github.com/libp2p/go-libp2p-webrtc-direct/issues/37")
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
