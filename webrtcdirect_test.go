package libp2pwebrtcdirect

import (
	"reflect"
	"runtime"
	"testing"

	logging "github.com/ipfs/go-log"

	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	utils "github.com/libp2p/go-libp2p-transport/test"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pions/webrtc"
	mplex "github.com/whyrusleeping/go-smux-multiplex"
)

var Subtests = []func(t *testing.T, ta, tb tpt.Transport, maddr ma.Multiaddr, peerA peer.ID){
	utils.SubtestProtocols,
	utils.SubtestBasic,

	utils.SubtestCancel,
	utils.SubtestPingPong,

	// Stolen from the stream muxer test suite.
	utils.SubtestStress1Conn1Stream1Msg,
	// utils.SubtestStress1Conn1Stream100Msg, // Flaky (WIP on SCTP issues)
	// utils.SubtestStress1Conn100Stream100Msg, // Flaky (WIP on SCTP issues)
	// utils.SubtestStress50Conn10Stream50Msg, // TODO
	// utils.SubtestStress1Conn1000Stream10Msg, // TODO
	// utils.SubtestStress1Conn100Stream100Msg10MB, // TODO
	// utils.SubtestStreamOpenStress, // Passes with higher timeout
	utils.SubtestStreamReset,
}

func TestTransport(t *testing.T) {
	logging.SetLogLevel("*", "warning")

	ta := NewTransport(
		webrtc.RTCConfiguration{},
		new(mplex.Transport),
	)
	tb := NewTransport(
		webrtc.RTCConfiguration{},
		new(mplex.Transport),
	)

	addr := "/ip4/127.0.0.1/tcp/0/http/p2p-webrtc-direct"
	SubtestTransport(t, ta, tb, addr, "peerA")
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func SubtestTransport(t *testing.T, ta, tb tpt.Transport, addr string, peerA peer.ID) {
	maddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range Subtests {
		t.Run(getFunctionName(f), func(t *testing.T) {
			f(t, ta, tb, maddr, peerA)
		})
	}
}
func TestTransportCantListenUtp(t *testing.T) {
	utpa, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/50000")
	if err != nil {
		t.Fatal(err)
	}

	tpt := NewTransport(
		webrtc.RTCConfiguration{},
		new(mplex.Transport),
	)

	_, err = tpt.Listen(utpa)
	if err == nil {
		t.Fatal("shouldnt be able to listen on utp addr with tcp transport")
	}

}
