package client

import (
	"crypto/ecdsa"
	"fmt"
	"net/http"

	"github.com/mrod502/ddns/config"
	"github.com/mrod502/logger"
)

type Client struct {
	cfg    config.Config
	pubKey *ecdsa.PrivateKey
	log    logger.Client
}

func NewClient()

func (c *Client) Ping() error {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/ping", c.cfg.RemoteHost), nil)

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
