package libp2pwebrtcdirect

import (
	"reflect"
	"runtime"
	"testing"

	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	utils "github.com/libp2p/go-libp2p-transport/test"
	ma "github.com/multiformats/go-multiaddr"
)

// The contents of this file are copied from libp2p/go-libp2p-transport/test
// in order to disable some tests while we investigate performance issues when
// running the tests on resource restricted environments like the Travis CI.

var Subtests = []func(t *testing.T, ta, tb tpt.Transport, maddr ma.Multiaddr, peerA peer.ID){
	utils.SubtestProtocols,
	utils.SubtestBasic,

	utils.SubtestCancel,
	// utils.SubtestPingPong,

	// Stolen from the stream muxer test suite.
	utils.SubtestStress1Conn1Stream1Msg,
	utils.SubtestStress1Conn1Stream100Msg,
	utils.SubtestStress1Conn100Stream100Msg,
	utils.SubtestStress50Conn10Stream50Msg,
	utils.SubtestStress1Conn1000Stream10Msg,
	utils.SubtestStress1Conn100Stream100Msg10MB,
	utils.SubtestStreamOpenStress,
	utils.SubtestStreamReset,
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
