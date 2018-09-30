package libp2pwebrtcdirect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	multibase "github.com/multiformats/go-multibase"
	"github.com/pions/dcnet"
	"github.com/pions/webrtc"
)

func NewHTTPDirectSignaler(config webrtc.RTCConfiguration, address string) *HTTPDirectSignaler {
	ctx, cancel := context.WithCancel(context.Background())
	return &HTTPDirectSignaler{
		config:  webrtc.RTCConfiguration{},
		address: address,
		ctx:     ctx,
		cancel:  cancel,
	}
}

type HTTPDirectSignaler struct {
	config  webrtc.RTCConfiguration
	address string
	ctx     context.Context
	cancel  func()
}

func (r *HTTPDirectSignaler) Dial() (*webrtc.RTCDataChannel, net.Addr, error) {
	c, err := webrtc.New(r.config)
	if err != nil {
		return nil, nil, err
	}

	var dc *webrtc.RTCDataChannel
	dc, err = c.CreateDataChannel("data", nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: migrate to OnNegotiationNeeded when available
	offer, err := c.CreateOffer(nil)
	if err != nil {
		return nil, nil, err
	}

	offerEnc, err := Encode(offer)
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.Get("http://" + r.address + "/?signal=" + offerEnc)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	answerEnc, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	answer, err := Decode(string(answerEnc))
	if err != nil {
		return nil, nil, err
	}

	if err := c.SetRemoteDescription(answer); err != nil {
		return nil, nil, err
	}
	return dc, &dcnet.NilAddr{}, nil
}

func (r *HTTPDirectSignaler) Accept() (*webrtc.RTCDataChannel, net.Addr, error) {
	c, err := webrtc.New(r.config)
	if err != nil {
		return nil, nil, err
	}
	//c.OnICEConnectionStateChange = func(connectionState ice.ConnectionState) {
	//	fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	//}

	var dc *webrtc.RTCDataChannel
	res := make(chan *webrtc.RTCDataChannel)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		signals, ok := r.Form["signal"]
		if !ok || len(signals) != 1 {
			fmt.Println("Failed get offer")
			return
		}

		offer, err := Decode(signals[0])
		if err != nil {
			fmt.Println("Failed to decode offer:", err)
			return
		}

		if err := c.SetRemoteDescription(offer); err != nil {
			fmt.Println("Failed to set remote description:", err)
			return
		}

		answer, err := c.CreateAnswer(nil)
		if err != nil {
			fmt.Println("Failed to create answer:", err)
			return
		}

		answerEnc, err := Encode(answer)
		if err != nil {
			fmt.Println("Failed to encode answer:", err)
			return
		}

		_, err = fmt.Fprint(w, answerEnc)
		if err != nil {
			fmt.Println("Failed to send answer:", err)
			return
		}

		c.OnDataChannel = func(d *webrtc.RTCDataChannel) {
			res <- d
		}

	})

	srv := &http.Server{
		Addr:    r.address,
		Handler: mux,
	}

	go srv.ListenAndServe()

	select {
	case dc = <-res:
	case <-r.ctx.Done():
		return nil, nil, errors.New("signaler closed")
	}
	return dc, &dcnet.NilAddr{}, nil
}

func Encode(desc webrtc.RTCSessionDescription) (string, error) {
	descData, err := json.Marshal(desc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal description: %v", err)
	}

	descEnc, err := multibase.Encode(multibase.Base58BTC, descData)
	if err != nil {
		return "", fmt.Errorf("failed to encode description: %v", err)
	}
	return descEnc, nil
}

func Decode(descEnc string) (webrtc.RTCSessionDescription, error) {
	var desc webrtc.RTCSessionDescription

	_, descData, err := multibase.Decode(descEnc)
	if err != nil {
		return desc, fmt.Errorf("failed to decode description: %v", err)
	}

	err = json.Unmarshal(descData, &desc)
	if err != nil {
		return desc, fmt.Errorf("failed to unmarshal description: %v", err)
	}

	return desc, nil
}

func (r *HTTPDirectSignaler) Close() error {
	r.cancel()
	return nil
}

func (r *HTTPDirectSignaler) Addr() net.Addr {
	return &dcnet.NilAddr{}
}
