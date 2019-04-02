package libp2pwebrtcdirect

import (
	"fmt"
	"net"

	peer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

func newMultiaddrFromNetAddr(netAddr net.Addr, peerID peer.ID) (ma.Multiaddr, error) {
	tcpMa, err := manet.FromNetAddr(netAddr)
	if err != nil {
		return nil, fmt.Errorf("failed create ma: %v", err)
	}
	httpMa := tcpMa.Encapsulate(httpma)
	webrtcMa, err := ma.NewMultiaddr("/p2p-webrtc/" + peer.IDB58Encode(peerID))
	if err != nil {
		return nil, err
	}
	maAddr := httpMa.Encapsulate(webrtcMa)
	return maAddr, nil
}

func getPeerIDFromMultiAddr(addr ma.Multiaddr) (peer.ID, error) {
	idString, err := addr.ValueForProtocol(protoCode)
	if err != nil {
		return "", err
	}
	return peer.IDB58Decode(idString)
}
