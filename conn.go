package libp2pwebrtcdirect

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"sync"
	"time"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	"github.com/pion/datachannel"
	"github.com/pion/webrtc/v2"
)

type connConfig struct {
	transport *Transport
	maAddr    ma.Multiaddr
	addr      net.Addr
	isServer  bool
	remoteID  peer.ID
}

func (c *connConfig) CopyWithNewRemoteID(remoteID peer.ID) *connConfig {
	return &connConfig{
		transport: c.transport,
		maAddr:    c.maAddr,
		addr:      c.addr,
		isServer:  c.isServer,
		remoteID:  remoteID,
	}
}

func newConnConfig(transport *Transport, maAddr ma.Multiaddr, isServer bool) (*connConfig, error) {
	httpMa := maAddr.Decapsulate(webrtcma)

	tcpMa := httpMa.Decapsulate(httpma)
	addr, err := manet.ToNetAddr(tcpMa)
	if err != nil {
		return nil, fmt.Errorf("failed to get net addr: %v", err)
	}

	return &connConfig{
		transport: transport,
		maAddr:    maAddr,
		addr:      addr,
		isServer:  isServer,
	}, nil
}

// Conn is a stream-multiplexing connection to a remote peer.
type Conn struct {
	config *connConfig

	peerConnection *webrtc.PeerConnection
	channel        datachannel.ReadWriteCloser

	lock   sync.RWMutex
	accept chan chan detachResult
	addr   net.Addr

	buf      []byte
	bufStart int
	bufEnd   int
}

func newConn(config *connConfig, pc *webrtc.PeerConnection, initChannel datachannel.ReadWriteCloser) *Conn {
	conn := &Conn{
		config:         config,
		peerConnection: pc,
		channel:        initChannel,
		accept:         make(chan chan detachResult),
		buf:            make([]byte, dcWrapperBufSize),
		addr:           config.addr,
	}

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		// We have to detach in OnDataChannel
		detachRes := detachChannel(dc)
		conn.accept <- detachRes
	})

	return conn
}

func dial(ctx context.Context, config *connConfig) (*Conn, error) {
	api := config.transport.api
	pc, err := api.NewPeerConnection(config.transport.webrtcOptions)
	if err != nil {
		return nil, err
	}

	dc, err := pc.CreateDataChannel("data", nil)
	if err != nil {
		return nil, err
	}

	detachRes := detachChannel(dc)

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		return nil, err
	}

	err = pc.SetLocalDescription(offer)
	if err != nil {
		return nil, err
	}

	offerEnc, err := encodeSignal(offer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "http://"+config.addr.String()+"/?signal="+offerEnc, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	answerEnc, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		return nil, err
	}

	answer, err := decodeSignal(string(answerEnc))
	if err != nil {
		return nil, err
	}

	if err := pc.SetRemoteDescription(answer); err != nil {
		return nil, err
	}

	select {
	case res := <-detachRes:
		if res.err != nil {
			return nil, res.err
		}
		return newConn(config, pc, res.dc), nil

	case <-ctx.Done():
		return newConn(config, pc, nil), ctx.Err()
	}
}

type detachResult struct {
	dc  datachannel.ReadWriteCloser
	err error
}

func detachChannel(dc *webrtc.DataChannel) chan detachResult {
	onOpenRes := make(chan detachResult)
	dc.OnOpen(func() {
		// Detach the data channel
		raw, err := dc.Detach()
		onOpenRes <- detachResult{raw, err}
	})

	return onOpenRes
}

// Close closes the stream muxer and the the underlying net.Conn.
func (c *Conn) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	var err error
	if c.peerConnection != nil {
		err = c.peerConnection.Close()
	}
	c.peerConnection = nil

	close(c.accept)

	newErr := c.channel.Close()
	if err == nil {
		err = newErr
	}
	return err
}

// IsClosed returns whether a connection is fully closed, so it can
// be garbage collected.
func (c *Conn) IsClosed() bool {
	c.lock.RLock()
	pc := c.peerConnection
	c.lock.RUnlock()
	return pc == nil
}

func (c *Conn) getPC() (*webrtc.PeerConnection, error) {
	c.lock.RLock()
	pc := c.peerConnection
	c.lock.RUnlock()

	if pc == nil {
		return nil, errors.New("Conn closed")
	}

	return pc, nil
}

// func (c *Conn) getMuxed() (smux.MuxedConn, error) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()

// 	if !c.isMuxed {
// 		return nil, nil
// 	}

// 	if c.muxedConn != nil {
// 		return c.muxedConn, nil
// 	}

// 	rawDC := c.initChannel
// 	if rawDC == nil {
// 		var err error
// 		rawDC, err = c.awaitAccept()
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	err := c.useMuxer(&dcWrapper{channel: rawDC, addr: c.config.addr, buf: make([]byte, dcWrapperBufSize)}, c.config.transport.muxer)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return c.muxedConn, nil
// }

// Note: caller should hold the conn lock.
// func (c *Conn) useMuxer(conn net.Conn, muxer smux.Multiplexer) error {
// 	muxed, err := muxer.NewConn(conn, c.config.isServer)
// 	if err != nil {
// 		return err
// 	}
// 	c.muxedConn = muxed

// 	return nil
// }

// func (c *Conn) checkInitChannel() datachannel.ReadWriteCloser {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	// Since a WebRTC offer can't be empty the offering side will have
// 	// an initial data channel opened. We return it here, the first time
// 	// OpenStream is called.
// 	if c.initChannel != nil {
// 		ch := c.initChannel
// 		c.initChannel = nil
// 		return ch
// 	}

// 	return nil
// }

func (c *Conn) awaitAccept() (datachannel.ReadWriteCloser, error) {
	detachRes, ok := <-c.accept
	if !ok {
		return nil, errors.New("Conn closed")
	}

	res := <-detachRes
	if res.err == nil {
		c.channel = res.dc
	}
	return res.dc, res.err
}

// LocalPeer returns our peer ID
func (c *Conn) LocalPeer() peer.ID {
	log.Debugf("local peer called, returning: %s", c.config.transport.localID)
	// TODO: Base on WebRTC security?
	return c.config.transport.localID
}

// LocalPrivateKey returns our private key
func (c *Conn) LocalPrivateKey() ic.PrivKey {
	// TODO: Base on WebRTC security?
	return nil

}

// RemotePeer returns the peer ID of the remote peer.
func (c *Conn) RemotePeer() peer.ID {
	log.Debugf("conn remote peer: %s", c.config.remoteID)
	// TODO: Base on WebRTC security?
	return c.config.remoteID
}

// RemotePublicKey returns the public key of the remote peer.
func (c *Conn) RemotePublicKey() ic.PubKey {
	// TODO: Base on WebRTC security?
	return nil
}

// LocalMultiaddr returns the local Multiaddr associated
// with this connection
func (c *Conn) LocalMultiaddr() ma.Multiaddr {
	return c.config.maAddr
}

// RemoteMultiaddr returns the remote Multiaddr associated
// with this connection
func (c *Conn) RemoteMultiaddr() ma.Multiaddr {
	log.Debugf("remote maddr: ", c.config.maAddr.String())
	return c.config.maAddr
}

// Transport returns the transport to which this connection belongs.
func (c *Conn) Transport() tpt.Transport {
	return c.config.transport
}

// Limit message size until we have a better
// packetizing strategy.
const dcWrapperBufSize = math.MaxUint16

// dcWrapper wraps datachannel.ReadWriteCloser to form a net.Conn
// type dcWrapper struct {
// 	channel datachannel.ReadWriteCloser
// 	addr    net.Addr

// 	buf      []byte
// 	bufStart int
// 	bufEnd   int
// }

func (c *Conn) Read(p []byte) (int, error) {
	if c.channel == nil {
		c.awaitAccept()
	}
	var err error

	if c.bufEnd == 0 {
		n := 0
		n, err = c.channel.Read(c.buf)
		c.bufEnd = n
	}

	n := 0
	if c.bufEnd-c.bufStart > 0 {
		n = copy(p, c.buf[c.bufStart:c.bufEnd])
		c.bufStart += n

		if c.bufStart >= c.bufEnd {
			c.bufStart = 0
			c.bufEnd = 0
		}
	}
	log.Debugf("read finished: %s", string(p))

	return n, err
}

func (c *Conn) Write(p []byte) (n int, err error) {
	if c.channel == nil {
		log.Debugf("awaiting accept")
		_, err := c.awaitAccept()
		if err != nil {
			return 0, fmt.Errorf("error awaiting accept: %v", err)
		}
	}
	if len(p) > dcWrapperBufSize {
		log.Debugf("write: %v", p[:dcWrapperBufSize])
		return c.channel.Write(p[:dcWrapperBufSize])
	}
	log.Debugf("write: (%v) %s", p, string(p))
	num, err := c.channel.Write(p)
	if err != nil {
		log.Errorf("error writing: %v", err)
	}
	return num, err
}

func (c *Conn) LocalAddr() net.Addr {
	return c.addr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}
