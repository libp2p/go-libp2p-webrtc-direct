package libp2pwebrtcdirect

import (
	"encoding/json"
	"fmt"

	multibase "github.com/multiformats/go-multibase"
	"github.com/pion/webrtc/v3"
)

func encodeSignal(desc webrtc.SessionDescription) (string, error) {
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

func decodeSignal(descEnc string) (webrtc.SessionDescription, error) {
	var desc webrtc.SessionDescription

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
