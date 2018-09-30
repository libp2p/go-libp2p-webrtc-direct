package libp2pwebrtcdirect

import (
	"fmt"
	"net"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

// laddr and raddr don't make much sense. Maybe they can be nil?
func wrapNetListener(listener net.Listener, laddr, raddr ma.Multiaddr) manet.Listener {
	return &maListener{
		Listener: listener,
		laddr:    laddr,
		raddr:    raddr,
	}
}

// maListener implements Listener
type maListener struct {
	net.Listener
	laddr ma.Multiaddr
	raddr ma.Multiaddr
}

// Accept waits for and returns the next connection to the listener.
// Returns a Multiaddr friendly Conn
func (l *maListener) Accept() (manet.Conn, error) {
	nconn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	fmt.Println("Accept: wrapNetConn")
	return wrapNetConn(nconn, l.laddr, l.raddr), nil
}

// Multiaddr returns the listener's (local) Multiaddr.
func (l *maListener) Multiaddr() ma.Multiaddr {
	return l.laddr
}

// laddr and raddr don't make much sense. Maybe they can be nil?
func wrapNetConn(conn net.Conn, laddr, raddr ma.Multiaddr) manet.Conn {
	endpts := maEndpoints{
		laddr: laddr,
		raddr: raddr,
	}

	return &struct {
		net.Conn
		maEndpoints
	}{conn, endpts}
}

type maEndpoints struct {
	laddr ma.Multiaddr
	raddr ma.Multiaddr
}

// LocalMultiaddr returns the local address associated with
// this connection
func (c *maEndpoints) LocalMultiaddr() ma.Multiaddr {
	return c.laddr
}

// RemoteMultiaddr returns the remote address associated with
// this connection
func (c *maEndpoints) RemoteMultiaddr() ma.Multiaddr {
	return c.raddr
}
