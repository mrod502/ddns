package server

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mrod502/ddns/config"
	"github.com/mrod502/ddns/interfaces"
	"github.com/mrod502/ddns/logger"
	"github.com/mrod502/ddns/util"
	"github.com/mrod502/encmsg/decoder"
	"github.com/vmihailenco/msgpack/v5"
)

type Server struct {
	cfg     config.Config
	privKey *rsa.PrivateKey
	*mux.Router
	interfaces.Logger
	lastIp string
	*decoder.Decoder
}

func NewServer(cfg config.Config) (s *Server, err error) {
	s = &Server{
		cfg:    cfg,
		Logger: logger.New(),
		Router: mux.NewRouter(),
	}
	s.HandleFunc("/ping", s.handlePing)
	s.HandleFunc("/current", s.handleCurrent)
	return
}

func (s *Server) Start() error {
	privKey, err := util.LoadPrivKey(s.cfg.PrivateKeyPath)

	if err != nil {
		return err
	}
	fmt.Println("loaded privKey")

	s.privKey = privKey
	return http.ListenAndServeTLS(fmt.Sprintf(":%d", s.cfg.Port), s.cfg.CertFilePath, s.cfg.KeyFilePath, s)
}

func New(cfg config.Config) *Server {
	priv, _ := util.LoadPrivKey(cfg.PrivateKeyPath)
	s := &Server{
		cfg:     cfg,
		Router:  mux.NewRouter(),
		Logger:  logger.New(),
		Decoder: decoder.New(decoder.NewRsaDecrypter(priv), msgpack.Unmarshal),
	}

	s.HandleFunc("/ping", s.handlePing)
	s.HandleFunc("/current", s.handleCurrent)
	return s
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var body util.RequestBody
	if err = s.Decode(b, &body); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	s.lastIp = parseIp(r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleCurrent(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	apiKey := vars.Get("apiKey")
	if apiKey != s.cfg.APIKey {
		s.Write("expected", s.cfg.APIKey, "got", apiKey)
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
