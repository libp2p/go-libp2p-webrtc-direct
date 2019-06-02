module github.com/libp2p/go-libp2p-webrtc-direct/examples/libp2p-echo

go 1.12

require (
	github.com/ipfs/go-log v0.0.1
	github.com/libp2p/go-libp2p v0.0.3
	github.com/libp2p/go-libp2p-core/crypto v0.0.1
	github.com/libp2p/go-libp2p-host v0.0.1
	github.com/libp2p/go-libp2p-net v0.0.1
	github.com/libp2p/go-libp2p-core/peer v0.1.0
	github.com/libp2p/go-libp2p-core/peerstore v0.0.1
	github.com/libp2p/go-libp2p-webrtc-direct v0.0.0-00010101000000-000000000000
	github.com/multiformats/go-multiaddr v0.0.2
	github.com/pion/webrtc/v2 v2.0.5
	github.com/whyrusleeping/go-logging v0.0.0-20170515211332-0457bb6b88fc
	github.com/whyrusleeping/go-smux-multiplex v3.0.16+incompatible
)

replace github.com/libp2p/go-libp2p-webrtc-direct => ../../
