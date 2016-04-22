package golf

import (
	"encoding/hex"
	"testing"
)

func TestRandomBytes(t *testing.T) {
	randomBytesLengths := []int{0, 32, 64, 1024}
	for length := range randomBytesLengths {
		randomBytes := randomBytes(length)
		if len(randomBytes) != length {
			t.Errorf("Could not create random byte array with length %v", length)
		}
	}
}

func TestDecodeXSRF(t *testing.T) {
	maskedToken := newXSRFToken()
	maskedTokenBytes, _ := hex.DecodeString(maskedToken)
	mask, token, _ := decodeXSRFToken(maskedToken)
	result := append(mask, websocketMask(mask, token)...)
	if !compareToken(maskedTokenBytes, result) {
		t.Errorf("Could not genearte correct XSRF token. %v != %v", maskedTokenBytes, result)
	}
}
