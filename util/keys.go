package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func ParseRsaPublicKeyFromPem(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("key type is not RSA")
}

func ExportRsaPublicKey(pubkey *rsa.PublicKey) (*pem.Block, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return nil, err
	}

	return &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubkey_bytes,
	}, nil
}

func ParseRsaPrivateKeyFromPem(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func GenerateRSAKeyPair(outFile string) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil || len(outFile) == 0 {
		return key, err
	}

	privFile, err := os.Create(outFile)
	if err != nil {
		return key, err
	}
	pubFile, err := os.Create(outFile + ".pub")
	if err != nil {
		return key, err
	}

	pubBlock, err := ExportRsaPublicKey(&key.PublicKey)
	if err != nil {
		return key, err
	}

	pem.Encode(pubFile, pubBlock)

	var pemkey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key)}

	err = pem.Encode(privFile, pemkey)

	return key, err

}

func LoadPubKey(path string) (*rsa.PublicKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseRsaPublicKeyFromPem(b)
}

func LoadPrivKey(path string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseRsaPrivateKeyFromPem(b)
}
