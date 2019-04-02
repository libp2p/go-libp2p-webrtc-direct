package libp2pwebrtcdirect

import (
	"math/rand"
	"testing"

	logging "github.com/ipfs/go-log"
	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	utils "github.com/libp2p/go-libp2p-transport/test"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pions/webrtc"
	mplex "github.com/whyrusleeping/go-smux-multiplex"
)

func newTestPeerID(t *testing.T) peer.ID {
	t.Helper()
	seededReader := rand.New(rand.NewSource(12345678))
	priv, _, err := crypto.GenerateEd25519Key(seededReader)
	if err != nil {
		t.Fatal(err)
	}
	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

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

	peerID := newTestPeerID(t)
	addr := "/ip4/127.0.0.1/tcp/0/http/p2p-webrtc/" + peer.IDB58Encode(peerID)
	utils.SubtestTransport(t, ta, tb, addr, peerID)
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
