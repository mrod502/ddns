package server

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mrod502/ddns/config"
	"github.com/mrod502/ddns/interfaces"
	"github.com/mrod502/ddns/logger"
	"github.com/mrod502/ddns/util"
)

type Server struct {
	cfg     config.Config
	privKey *rsa.PrivateKey
	*mux.Router
	log    interfaces.Logger
	lastIp string
}

func NewServer(cfg config.Config) (s *Server, err error) {
	s = &Server{
		cfg:    cfg,
		log:    logger.New(),
		Router: mux.NewRouter(),
	}
	s.HandleFunc("/ping", s.handlePing)
	s.HandleFunc("/current", s.handleCurrent)
	return
}

func (s *Server) Start() error {
	privKey, err := util.LoadPrivKey(s.cfg.PrivateKeyPath)
	fmt.Println("loaded privKey")

	if err != nil {
		return err
	}
	s.privKey = privKey
	//return http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), s)
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", s.cfg.Port), s.cfg.CertFilePath, s.cfg.KeyFilePath, s)
}

func New(cfg config.Config) *Server {
	s := &Server{
		cfg:    cfg,
		Router: mux.NewRouter(),
		log:    logger.New(),
	}

	s.HandleFunc("/ping", s.handlePing)
	s.HandleFunc("/current", s.handleCurrent)
	return s
}

func (s *Server) verifySignature(w http.ResponseWriter, r *http.Request) (*util.RequestBody, error) {
	body, err := util.DecodeFromHttpRequest(r, s.privKey)
	if err != nil {
		s.log.Write("process ping:", err.Error())
		s.log.Write("rand is:", r.Header.Get(util.HRand))
		return body, err
	}
	return body, nil
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	_, err := s.verifySignature(w, r)
	if err != nil {
		s.log.Write("ERROR:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.lastIp = parseIp(r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleCurrent(w http.ResponseWriter, r *http.Request) {
	//_, err := s.verifySignature(w, r)
	//if err != nil {
	//		s.log.Write("error")
	//		http.Error(w, err.Error(), http.StatusBadRequest)
	//		return
	//	}
	vars := r.URL.Query()
	apiKey := vars.Get("apiKey")
	if apiKey != s.cfg.APIKey {
		s.log.Write("expected", s.cfg.APIKey, "got", apiKey)
		http.Error(w, "fuck off\n", http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("%s\n", s.lastIp)))
}

func parseIp(inp string) string {
	out := strings.Split(inp, ":")
	if len(out) == 0 {
		return ""
	}
	return out[0]
}
