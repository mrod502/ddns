package server

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrod502/ddns/config"
	"github.com/mrod502/ddns/util"
	"github.com/mrod502/logger"
)

type Server struct {
	cfg     config.Config
	privKey *rsa.PrivateKey
	*mux.Router
	log logger.Client
}

func NewServer(cfg config.Config) (s *Server, err error) {
	s = &Server{
		cfg:    cfg,
		Router: mux.NewRouter(),
	}
	s.HandleFunc("/ping", s.handlePing)
	return
}

func (s *Server) Start() error {
	privKey, err := util.LoadPrivKey(s.cfg.PrivateKeyPath)

	if err != nil {
		return err
	}
	s.privKey = privKey

	cli, err := logger.NewClient(s.cfg.ClientConfig)
	if err != nil {
		return err
	}
	s.log = cli
	err = s.log.Connect()
	if err != nil {
		return err
	}
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", s.cfg.Port), s.cfg.CertFilePath, s.cfg.KeyFilePath, s)
}

func New(cfg config.Config) *Server {
	s := &Server{
		cfg:    cfg,
		Router: mux.NewRouter(),
	}
	s.HandleFunc("/ping", s.handlePing)
	return s
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {

	body, err := util.DecodeFromHttpRequest(r, s.privKey)
	if err != nil {
		s.log.Write("process ping:", err.Error())
		s.log.Write("rand is:", r.Header.Get(util.HRand))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("%+v\n", *body)
	w.WriteHeader(http.StatusOK)
}
