package encrypt

import (
	"crypto/hmac"
	"crypto/sha256"
)

func GetSHA256Signature(msg, key []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(msg); err != nil {
		return nil, err
	}
	data := mac.Sum(nil)
	return data, nil
}
