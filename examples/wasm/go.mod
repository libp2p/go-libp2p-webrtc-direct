module github.com/libp2p/go-libp2p-webrtc-direct/examples/wasm

go 1.16

replace github.com/libp2p/go-libp2p-webrtc-direct => ../../

require (
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/libp2p/go-libp2p-mplex v0.4.1
	github.com/libp2p/go-libp2p-webrtc-direct v0.0.0-00010101000000-000000000000
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/pion/webrtc/v3 v3.0.16
	github.com/shurcooL/go-goon v0.0.0-20210110234559-7585751d9a17 // indirect
)
