module github.com/libp2p/go-libp2p-webrtc-direct/examples/libp2p-echo

go 1.16

replace github.com/libp2p/go-libp2p-webrtc-direct => ../../

require (
	github.com/ipfs/go-log v1.0.5
	github.com/libp2p/go-libp2p v0.15.1
	github.com/libp2p/go-libp2p-core v0.9.0
	github.com/libp2p/go-libp2p-mplex v0.4.1
	github.com/libp2p/go-libp2p-webrtc-direct v0.0.0-20201219114432-56b02029fbb8
	github.com/multiformats/go-multiaddr v0.4.0
	github.com/pion/webrtc/v3 v3.0.16
	github.com/whyrusleeping/go-logging v0.0.1
)
