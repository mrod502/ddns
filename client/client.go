package client

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/mrod502/ddns/config"
	"github.com/mrod502/ddns/interfaces"
	"github.com/mrod502/ddns/logger"
	"github.com/mrod502/ddns/util"
	"github.com/mrod502/encmsg/encoder"
	"github.com/vmihailenco/msgpack/v5"
)

type Client struct {
	cfg    config.Config
	pubKey *rsa.PublicKey
	log    interfaces.Logger
	*encoder.Encoder
}

func (c *Client) Start() error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	for {
		err := c.Ping()
		if err != nil {
			c.log.Write(err.Error())
		}
		time.Sleep(c.cfg.PingInterval)
	}

}

func (c *Client) Ping() error {
	req, err := c.buildRequest()
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}

func New(cfg config.Config) *Client {
	pub, err := util.LoadPubKey(cfg.PublicKeyPath)
	if err != nil {
		panic(err)
	}

	return &Client{
		cfg:     cfg,
		pubKey:  pub,
		log:     logger.New(),
		Encoder: encoder.New(encoder.NewRsaEncrypter(pub), msgpack.Marshal),
	}
}

func (c *Client) buildRequest() (*http.Request, error) {
	ts := time.Now()
	body := util.RequestBody{
		Timestamp: ts,
		Domain:    c.cfg.Domain,
	}
	b, err := c.Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s:%d/ping", c.cfg.RemoteHost, c.cfg.Port), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return req, nil
}
