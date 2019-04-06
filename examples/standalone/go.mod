module github.com/libp2p/go-libp2p-webrtc-direct/examples/standalone

go 1.12

require (
	github.com/libp2p/go-libp2p-transport v0.0.4
	github.com/libp2p/go-libp2p-webrtc-direct v0.0.0-00010101000000-000000000000
	github.com/libp2p/go-stream-muxer v0.0.1
	github.com/multiformats/go-multiaddr v0.0.2
	github.com/pion/webrtc/v2 v2.0.5
	github.com/whyrusleeping/go-smux-multiplex v3.0.16+incompatible
)

replace github.com/libp2p/go-libp2p-webrtc-direct => ../../
