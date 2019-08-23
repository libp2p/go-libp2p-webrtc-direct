package libp2pwebrtcdirect

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

// Listener is an interface closely resembling the net.Listener interface.
type Listener struct {
	config *connConfig
	accept chan *Conn

	srv *http.Server
}

func newListener(config *connConfig) (*Listener, error) {
	lnet, lnaddr, err := manet.DialArgs(config.maAddr)
	if err != nil {
		return nil, err
	}

	ln, err := net.Listen(lnet, lnaddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	// Update the addr after listening
	tcpMa, err := manet.FromNetAddr(ln.Addr())
	if err != nil {
		return nil, fmt.Errorf("failed create ma: %v", err)
	}

	httpMa := tcpMa.Encapsulate(httpma)
	maAddr := httpMa.Encapsulate(webrtcma)

	config.addr = ln.Addr()
	config.maAddr = maAddr

	l := &Listener{
		config: config,
		accept: make(chan *Conn, 1),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", l.handler)

	srv := &http.Server{
		Handler: mux,
	}

	go func() {
		log.Debug("srv serving")
		srvErr := srv.Serve(ln)
		if srvErr != nil {
			log.Warningf("failed to start server: %v", srvErr)
		}
	}()

	l.srv = srv
	return l, nil
}

func (l *Listener) handler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handler")
	r.ParseForm()
	signals, ok := r.Form["signal"]
	if !ok || len(signals) != 1 {
		log.Warningf("failed to handle request: failed to parse signal")
		return
	}

	answer, err := l.handleSignal(signals[0])
	if err != nil {
		log.Warningf("failed to handle request: failed to setup connection: %v", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, err = fmt.Fprint(w, answer)
	if err != nil {
		log.Warningf("failed to handle request: failed to send answer: %v", err)
		return
	}
}

func (l *Listener) handleSignal(offerStr string) (string, error) {
	log.Debugf("handle signal")
	offer, err := decodeSignal(offerStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode offer: %v", err)
	}

	api := l.config.transport.api
	pc, err := api.NewPeerConnection(l.config.transport.webrtcOptions)
	if err != nil {
		return "", err
	}

	if err := pc.SetRemoteDescription(offer); err != nil {
		return "", fmt.Errorf("failed to set remote description: %v", err)
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create answer: %v", err)
	}

	err = pc.SetLocalDescription(answer)
	if err != nil {
		return "", fmt.Errorf("failed to set local description: %v", err)
	}

	answerEnc, err := encodeSignal(answer)
	if err != nil {
		return "", fmt.Errorf("failed to encode answer: %v", err)
	}

	c := newConn(l.config.CopyWithNewRemoteID(peer.ID(offer.SDP)), pc, nil)
	go func() {
		err := c.awaitAccept(30 * time.Second)
		if err != nil {
			log.Warningf("error awaiting channel accept: %v", err)
			return
		}
		l.accept <- c
	}()
	log.Debug("signal handled")
	return answerEnc, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (transport.CapableConn, error) {
	conn, ok := <-l.accept
	if !ok {
		return nil, errors.New("Listener closed")
	}

	return l.config.transport.Upgrader.UpgradeInbound(context.Background(), l.config.transport, conn)
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() error {
	if l.srv == nil {
		return nil // don't close twice
	}
	err := l.srv.Shutdown(context.Background())
	if err != nil {
		return err
	}
	l.srv = nil

	close(l.accept)

	return nil
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.config.addr
}

// Multiaddr returns the listener's network Multi address.
func (l *Listener) Multiaddr() ma.Multiaddr {
	return l.config.maAddr
}
