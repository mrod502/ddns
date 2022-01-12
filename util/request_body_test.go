package util

import (
	"testing"
	"time"
)

func TestEncode(t *testing.T) {

	priv, err := GenerateRSAKeyPair("foo")
	if err != nil {
		t.Fatal(err)
	}
	r := RequestBody{
		Timestamp: time.Now(),
		Domain:    "https://www.google.com/",
		APIKey:    []byte("foo"),
	}

	data, sig, err := r.Encode(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Decode(data, sig, priv)
	if err != nil {
		t.Fatal(err)
	}

}
