package libp2pwebrtcdirect

import (
	"context"
	"fmt"

	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	ma "github.com/multiformats/go-multiaddr"
	mafmt "github.com/multiformats/go-multiaddr-fmt"
	"github.com/pion/webrtc/v2"
)

// Transport is the WebRTC transport.
type Transport struct {
	webrtcOptions webrtc.Configuration
	localID       peer.ID
	api           *webrtc.API
	Upgrader      *tptu.Upgrader
}

func NewTransport(upgrader *tptu.Upgrader) *Transport {
	s := webrtc.SettingEngine{}
	// Use Detach data channels mode
	s.DetachDataChannels()

	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	if upgrader == nil {
		return nil
	}

	return &Transport{
		webrtcOptions: webrtc.Configuration{},
		localID:       peer.ID(1),
		api:           api,
		Upgrader:      upgrader,
	}
}

// CanDial returns true if this transport believes it can dial the given
// multiaddr.
func (t *Transport) CanDial(addr ma.Multiaddr) bool {
	return mafmt.WebRTCDirect.Matches(addr)
}

// Dial dials the peer at the remote address.
func (t *Transport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	if !t.CanDial(raddr) {
		return nil, fmt.Errorf("can't dial address %s", raddr)
	}

	cfg, err := newConnConfig(t, raddr, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get dial args: %v", err)
	}

	cfg.remoteID = p

	conn, err := dial(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %v", err)
	}

	return t.Upgrader.UpgradeOutbound(ctx, t, conn, p)
}

// Listen listens on the given multiaddr.
func (t *Transport) Listen(laddr ma.Multiaddr) (transport.Listener, error) {
	if !t.CanDial(laddr) {
		return nil, fmt.Errorf("can't listen on address %s", laddr)
	}

	cfg, err := newConnConfig(t, laddr, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get dial args: %v", err)
	}

	l, err := newListener(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	return l, nil
}

// Protocols returns the list of terminal protocols this transport can dial.
func (t *Transport) Protocols() []int {
	return []int{ma.P_P2P_WEBRTC_DIRECT}
}

// Proxy always returns false for the TCP transport.
func (t *Transport) Proxy() bool {
	return false
}

func (t *Transport) String() string {
	return "p2p-webrtc-direct"
}
