package main

import (
	"context"
	"fmt"

	mplex "github.com/libp2p/go-libp2p-mplex"
	direct "github.com/libp2p/go-libp2p-webrtc-direct"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/webrtc/v3"
)

func main() {
	maddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/9090/http/p2p-webrtc-direct")
	check(err)

	transport := direct.NewTransport(
		webrtc.Configuration{},
		new(mplex.Transport),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := transport.Dial(ctx, maddr, "peerA")
	check(err)
	defer c.Close()
	fmt.Println("[dialer] Opened connection")

	s, err := c.OpenStream(context.Background())
	check(err)
	fmt.Println("[dialer] Opened stream")

	_, err = s.Write([]byte("hey, how is it going. I am the dialer"))
	check(err)

	err = s.Close()
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
