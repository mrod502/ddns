package client

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/mrod502/ddns/config"
	"github.com/mrod502/ddns/util"
	"github.com/mrod502/logger"
)

type Client struct {
	cfg    config.Config
	pubKey *rsa.PublicKey
	log    logger.Client
}

func (c *Client) Start() error {
	pub, err := util.LoadPubKey(c.cfg.PublicKeyPath)
	if err != nil {
		return err
	}
	c.pubKey = pub

	for {
		err = c.Ping()
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
	cli, err := logger.NewClient(cfg.ClientConfig)

	if err != nil {
		panic(err)
	}
	err = cli.Connect()
	if err != nil {
		panic(err)
	}
	return &Client{
		cfg: cfg,
		log: cli,
	}
}

func (c *Client) buildRequest() (*http.Request, error) {
	ts := time.Now()
	body := util.RequestBody{
		Timestamp: ts,
		Domain:    c.cfg.Domain,
	}

	encoded, rand, err := body.Encode(c.pubKey)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s:%d/ping", c.cfg.RemoteHost, c.cfg.Port), bytes.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	randStr := base64.URLEncoding.EncodeToString(rand)

	req.Header.Set(util.HRand, randStr)
	req.Header.Set(util.HTimestamp, fmt.Sprintf("%d", ts.UnixNano()))

	return req, nil
}
