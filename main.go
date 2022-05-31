package main

import (
	"flag"

	"github.com/mrod502/ddns/client"
	"github.com/mrod502/ddns/config"
	"github.com/mrod502/ddns/server"
	"github.com/mrod502/ddns/util"
)

type runner interface {
	Start() error
}

func main() {

	gk := flag.Bool("genkey", false, "generate key pair")
	fname := flag.String("o", "ddns_key", "name for key file")

	srv := flag.Bool("server", false, "run as server")
	fpath := flag.String("cfg", "", "config file path")
	flag.Parse()

	if *gk {
		util.GenerateRSAKeyPair(*fname)
		return
	}
	var r runner

	cfg, err := config.Parse(*fpath)
	if err != nil {
		panic(err)
	}

	if *srv {
		r = server.New(cfg)
	} else {
		r = client.New(cfg)
	}

	err = r.Start()

	if err != nil {
		panic(err)
	}

}
