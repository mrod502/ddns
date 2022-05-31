package auth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

type TimestampedPayload interface {
	Timestamp() time.Time
	Payload() []byte
}

func TimestampPayload(obj TimestampedPayload) []byte {
	res := sha256.Sum256([]byte(fmt.Sprintf("%s\n%s", obj.Timestamp(), string(obj.Payload()))))
	return res[:]
}

func Authorize(obj TimestampedPayload, key *ecdsa.PrivateKey) (string, error) {
	key.Sign(rand.Reader, TimestampPayload(obj), nil)
	b, err := ecdsa.SignASN1(rand.Reader, key, TimestampPayload(obj))

	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
