package libp2pwebrtcdirect

import (
	"context"
	"fmt"

	logging "github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-peer"
	tpt "github.com/libp2p/go-libp2p-transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/pions/dcnet"
	"github.com/pions/webrtc"
	mafmt "github.com/whyrusleeping/mafmt"
)

var log = logging.Logger("webrtcdirect-tpt")

var webrtcma, _ = ma.NewMultiaddr("/p2p-webrtc-direct")

// WebRTCDirectTransport is the TCP transport.
type WebRTCDirectTransport struct {
	webrtcOptions webrtc.RTCConfiguration
	// Connection upgrader for upgrading insecure stream connections to
	// secure multiplex connections.
	Upgrader *tptu.Upgrader
}

var _ tpt.Transport = &WebRTCDirectTransport{}

// NewWebRTCDirectTransport creates a WebRTC transport that signals over a direct HTTP connection.
func NewWebRTCDirectTransport(webrtcOptions webrtc.RTCConfiguration, upgrader *tptu.Upgrader) *WebRTCDirectTransport {
	return &WebRTCDirectTransport{
		webrtcOptions: webrtcOptions,
		Upgrader:      upgrader,
	}
}

// CanDial returns true if this transport believes it can dial the given
// multiaddr.
func (t *WebRTCDirectTransport) CanDial(addr ma.Multiaddr) bool {
	return mafmt.WebRTCDirect.Matches(addr)
}

// Dial dials the peer at the remote address.
func (t *WebRTCDirectTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (tpt.Conn, error) {
	if !t.CanDial(raddr) {
		return nil, fmt.Errorf("can't dial address %s", raddr)
	}
	httpMa := raddr.Decapsulate(webrtcma)
	_, httpAddr, err := manet.DialArgs(httpMa)
	if err != nil {
		return nil, fmt.Errorf("failed to get dial args: %v", err)
	}

	signaler := NewHTTPDirectSignaler(t.webrtcOptions, httpAddr)
	conn, err := dcnet.Dial(signaler)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	wrappedConn := wrapNetConn(conn, raddr, raddr)

	return t.Upgrader.UpgradeOutbound(ctx, t, wrappedConn, p)
}

// Listen listens on the given multiaddr.
func (t *WebRTCDirectTransport) Listen(laddr ma.Multiaddr) (tpt.Listener, error) {
	if !t.CanDial(laddr) {
		return nil, fmt.Errorf("can't listen on address %s", laddr)
	}
	httpMa := laddr.Decapsulate(webrtcma)
	_, httpAddr, err := manet.DialArgs(httpMa)
	if err != nil {
		return nil, fmt.Errorf("failed to get dial args: %v", err)
	}

	signaler := NewHTTPDirectSignaler(t.webrtcOptions, httpAddr)
	listener := dcnet.NewListener(signaler)

	wrappedListener := wrapNetListener(listener, laddr, laddr)

	return t.Upgrader.UpgradeListener(t, wrappedListener), nil
}

// Protocols returns the list of terminal protocols this transport can dial.
func (t *WebRTCDirectTransport) Protocols() []int {
	return []int{ma.P_P2P_WEBRTC_DIRECT}
}

// Proxy always returns false for the TCP transport.
func (t *WebRTCDirectTransport) Proxy() bool {
	return false
}

func (t *WebRTCDirectTransport) String() string {
	return "p2p-webrtc-direct"
}
