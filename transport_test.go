package libp2pwebrtcdirect

import (
	"reflect"
	"runtime"
	"testing"

	detectrace "github.com/jbenet/go-detect-race"
	peer "github.com/libp2p/go-libp2p-core/peer"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	ttransport "github.com/libp2p/go-libp2p-testing/suites/transport"
	ma "github.com/multiformats/go-multiaddr"
)

// The contents of this file are copied from libp2p/go-libp2p-core/transport/test
// in order to disable some tests while we investigate performance issues when
// running the tests on resource restricted environments like the Travis CI.

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func SubtestTransport(t *testing.T, ta, tb tpt.Transport, addr string, peerA peer.ID) {

	subtests := []func(t *testing.T, ta, tb tpt.Transport, maddr ma.Multiaddr, peerA peer.ID){}

	if detectrace.WithRace() {
		subtests = []func(t *testing.T, ta, tb tpt.Transport, maddr ma.Multiaddr, peerA peer.ID){
			ttransport.SubtestProtocols,
			ttransport.SubtestBasic,

			ttransport.SubtestCancel,
			ttransport.SubtestPingPong,

			// Stolen from the stream muxer test suite.
			ttransport.SubtestStress1Conn1Stream1Msg,
		}
	} else {
		subtests = []func(t *testing.T, ta, tb tpt.Transport, maddr ma.Multiaddr, peerA peer.ID){
			ttransport.SubtestStress1Conn1Stream100Msg,
			ttransport.SubtestStress1Conn100Stream100Msg,
			ttransport.SubtestStress50Conn10Stream50Msg,
			ttransport.SubtestStress1Conn1000Stream10Msg,
			ttransport.SubtestStress1Conn100Stream100Msg10MB,
			ttransport.SubtestStreamOpenStress,
			ttransport.SubtestStreamReset,
		}
	}

	maddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range subtests {
		t.Run(getFunctionName(f), func(t *testing.T) {
			f(t, ta, tb, maddr, peerA)
		})
	}
}
