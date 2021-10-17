module github.com/libp2p/go-libp2p-webrtc-direct/examples/libp2p-echo

go 1.16

replace github.com/libp2p/go-libp2p-webrtc-direct => ../../

require (
	github.com/ipfs/go-log v1.0.5
	github.com/ipfs/go-log/v2 v2.1.3
	github.com/koron/go-ssdp v0.0.2 // indirect
	github.com/libp2p/go-conn-security-multistream v0.2.1
	github.com/libp2p/go-libp2p v0.13.0
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/libp2p/go-libp2p-mplex v0.4.1
	github.com/libp2p/go-libp2p-transport-upgrader v0.4.2
	github.com/libp2p/go-libp2p-webrtc-direct v0.0.0-20201219114432-56b02029fbb8
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/pion/webrtc/v3 v3.0.16
	github.com/whyrusleeping/go-logging v0.0.1
	golang.org/x/text v0.3.5 // indirect
)
