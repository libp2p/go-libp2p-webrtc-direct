package libp2pwebrtcdirect

import (
	"crypto/rand"
	"testing"

	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	mplex "github.com/libp2p/go-libp2p-mplex"
	peer "github.com/libp2p/go-libp2p-peer"
	secio "github.com/libp2p/go-libp2p-secio"

	logging "github.com/ipfs/go-log"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
)

func newUpgrader(t *testing.T) (*tptu.Upgrader, libp2pcrypto.PrivKey) {
	keyPriv, _, err := libp2pcrypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		t.Fatalf("error creating key: %v", err)
	}
	secTransp, err := secio.New(keyPriv)
	if err != nil {
		t.Fatalf("error creating secio transport: %v", err)
	}
	return &tptu.Upgrader{
		Muxer:  mplex.DefaultTransport,
		Secure: secTransp,
	}, keyPriv
}

func TestTransport(t *testing.T) {
	logging.SetLogLevel("*", "debug")
	aUpgrade, aKey := newUpgrader(t)
	bUpgrade, _ := newUpgrader(t)

	aId, err := peer.IDFromPublicKey(aKey.GetPublic())
	if err != nil {
		t.Fatalf("error getting id: %v", err)
	}

	ta := NewTransport(aUpgrade)
	tb := NewTransport(bUpgrade)

	addr := "/ip4/127.0.0.1/tcp/0/http/p2p-webrtc-direct"

	// TODO: Re-enable normal test suite when not hitting CI limits when using race detector
	// utils.SubtestTransport(t, ta, tb, addr, "peerA")
	SubtestTransport(t, ta, tb, addr, aId)
}
