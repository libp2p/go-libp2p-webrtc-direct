package libp2pwebrtcdirect

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/libp2p/go-conn-security/insecure"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	utils "github.com/libp2p/go-libp2p-transport/test"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pions/webrtc"
	mplex "github.com/whyrusleeping/go-smux-multiplex"
)

var Subtests = []func(t *testing.T, ta, tb tpt.Transport, maddr ma.Multiaddr, peerA peer.ID){
	// utils.SubtestProtocols,
	// utils.SubtestBasic,
	// utils.SubtestCancel,
	// utils.SubtestPingPong,

	// Stolen from the stream muxer test suite.
	utils.SubtestStress1Conn1Stream1Msg,
	// utils.SubtestStress1Conn1Stream100Msg,
	// utils.SubtestStress1Conn100Stream100Msg,
	// utils.SubtestStress50Conn10Stream50Msg,
	// utils.SubtestStress1Conn1000Stream10Msg,
	// utils.SubtestStress1Conn100Stream100Msg10MB,
	// utils.SubtestStreamOpenStress,
	// utils.SubtestStreamReset,
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

func TestWebRTCDirectTransport(t *testing.T) {
	ta := NewWebRTCDirectTransport(
		webrtc.RTCConfiguration{},
		&tptu.Upgrader{
			Secure: insecure.New("peerA"),
			Muxer:  new(mplex.Transport),
		},
	)
	tb := NewWebRTCDirectTransport(
		webrtc.RTCConfiguration{},
		&tptu.Upgrader{
			Secure: insecure.New("peerB"),
			Muxer:  new(mplex.Transport),
		},
	)

	addr := "/ip4/127.0.0.1/tcp/50000/http/p2p-webrtc-direct"
	SubtestTransport(t, ta, tb, addr, "peerA")
}

func TestWebRTCDirectTransportCantListenUtp(t *testing.T) {
	utpa, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5001")
	if err != nil {
		t.Fatal(err)
	}

	tpt := NewWebRTCDirectTransport(
		webrtc.RTCConfiguration{},
		&tptu.Upgrader{
			Secure: insecure.New("peerB"),
			Muxer:  new(mplex.Transport),
		},
	)

	_, err = tpt.Listen(utpa)
	if err == nil {
		t.Fatal("shouldnt be able to listen on utp addr with tcp transport")
	}

}
