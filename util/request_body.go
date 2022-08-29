package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type RequestBody struct {
	Timestamp time.Time
	Domain    string
	APIKey    []byte
	Data      []byte
}

func (r RequestBody) Encode(pub *rsa.PublicKey) (data, nonce []byte, err error) {
	msgBytes, err := msgpack.Marshal(r)
	if err != nil {
		return
	}

	if err != nil {
		return
	}
	nonce = GenRand()

	data, err = rsa.EncryptOAEP(
		sha256.New(),
		bytes.NewReader(nonce),
		pub,
		msgBytes,
		[]byte("msg"),
	)

	if err != nil {
		return
	}

	return
}

func Decode(data, nonce []byte, priv *rsa.PrivateKey) (r *RequestBody, err error) {

	decoded, err := rsa.DecryptOAEP(sha256.New(), bytes.NewReader(nonce), priv, data, []byte("msg"))

	if err != nil {

		return nil, err
	}

	r = new(RequestBody)

	err = msgpack.Unmarshal(decoded, r)

	return
}

func GenRand() []byte {
	var b []byte = make([]byte, 0, 0x400)

	for i := 0; i < 0x400; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(0xff))
		b = append(b, byte(num.Int64()))
	}
	return b
}

func DecodeFromHttpRequest(req *http.Request, priv *rsa.PrivateKey) (*RequestBody, error) {
	data, err := io.ReadAll(req.Body)
	if err != nil {

		return nil, err
	}
	decodedRand, err := base64.URLEncoding.DecodeString(req.Header.Get(HRand))
	if err != nil {

		return nil, err
	}

	return Decode(data, decodedRand, priv)
}
