package golf

import (
	"crypto/rand"
	"encoding/hex"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func randomBytes(strlen int) []byte {
	b := make([]byte, strlen)
	_, err := rand.Read(b)
	if err != nil {
		// panic on failure since this indicates a failure of the system's
		// CSPRNG
		panic(err)
	}
	return b
}

func decodeXSRFToken(maskedToken string) ([]byte, []byte, error) {
	maskedTokenBytes, err := hex.DecodeString(maskedToken)
	if err != nil {
		return nil, nil, err
	}
	mask := maskedTokenBytes[0:4]
	token := websocketMask(mask, maskedTokenBytes[4:])
	return mask, token, nil
}

func websocketMask(mask, data []byte) []byte {
	for i, v := range data {
		data[i] = v ^ mask[i%4]
	}
	return data
}

func compareToken(tokenA, tokenB []byte) bool {
	if tokenA == nil && tokenB == nil {
		return true
	}
	if tokenA == nil || tokenB == nil {
		return false
	}
	if len(tokenA) != len(tokenB) {
		return false
	}
	for i := range tokenA {
		if tokenA[i] != tokenB[i] {
			return false
		}
	}
	return true
}

func newXSRFToken() string {
	tokenBytes := randomBytes(32)
	maskBytes := randomBytes(4)
	maskedTokenBytes := append(maskBytes, websocketMask(maskBytes, tokenBytes)...)
	maskedToken := hex.EncodeToString(maskedTokenBytes)
	return maskedToken
}
