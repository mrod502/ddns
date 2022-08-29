package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/mrod502/ddns/client"
	"github.com/mrod502/ddns/config"
	"github.com/mrod502/ddns/server"
)

func TestMain(t *testing.T) {

	s := server.New(config.Config{
		CertFilePath:   "cert.pem",
		KeyFilePath:    "key.pem",
		Port:           3391,
		PrivateKeyPath: "util/ddns_key",
		PublicKeyPath:  "util/ddns_key.pub",
		APIKey:         "LemmeIn",
	})
	go func() {
		err := s.Start()
		if err != nil {
			fmt.Println("FAILED TO SERVE", err.Error())
			t.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	cli := client.New(config.Config{
		PingInterval:  time.Second,
		Port:          3391,
		RemoteHost:    "localhost",
		PublicKeyPath: "util/ddns_key.pub",
	})
	go func() {
		if err := cli.Start(); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(time.Minute * 100)
	err := cli.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
