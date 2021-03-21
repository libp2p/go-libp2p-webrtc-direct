package libp2pwebrtcdirect

import (
	"io"
	"time"

	"github.com/pion/datachannel"
)

// Stream is a bidirectional io pipe within a connection.
type Stream struct {
	channel datachannel.ReadWriteCloser
}

func newStream(channel datachannel.ReadWriteCloser) *Stream {
	return &Stream{channel: channel}
}

// Read implements the io.Reader.
func (s *Stream) Read(p []byte) (int, error) {
	i, err := s.channel.Read(p)
	if err != nil {
		// pions/datachannel retuns an error when the underlying transport
		// is closed. Here we turn this into EOF.
		return i, io.EOF
	}
	return i, nil
}

// Write implements the io.Writer.
func (s *Stream) Write(p []byte) (int, error) {
	return s.channel.Write(p)
}

// CloseRead closes the stream for writing. Reading will still work (that
// is, the remote side can still write).
func (s *Stream) CloseRead() error {
	// TODO: figure out close vs reset
	return nil
}

// CloseWrite closes the stream for writing. Reading will still work (that
// is, the remote side can still write).
func (s *Stream) CloseWrite() error {
	// TODO: figure out close vs reset
	return nil
}

// Close closes the stream for writing. Reading will still work (that
// is, the remote side can still write).
func (s *Stream) Close() error {
	// TODO: figure out close vs reset
	return nil
}

// Reset closes both ends of the stream. Use this to tell the remote
// side to hang up and go away.
func (s *Stream) Reset() error {
	// TODO: figure out close vs reset
	return s.channel.Close()
}

// SetDeadline is a stub
func (s *Stream) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline is a stub
func (s *Stream) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline is a stub
func (s *Stream) SetWriteDeadline(t time.Time) error {
	return nil
}
