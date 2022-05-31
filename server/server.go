package server

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mrod502/ddns/config"
	gocache "github.com/mrod502/go-cache"
)

type logger struct{}

func (l logger) Write(v ...string) error {
	fmt.Println(v)
	return nil
}

var (
	ErrInvalid = errors.New("invalid credentials")
)

type Server struct {
	cfg    config.Config
	pubKey *ecdsa.PublicKey
	*mux.Router
	log logger
	ip  *gocache.Container[string]
}

func NewServer(cfg config.Config) (s *Server, err error) {
	s = &Server{
		cfg:    cfg,
		Router: mux.NewRouter(),
		ip:     gocache.NewContainer("", time.Now()),
	}

	s.HandleFunc("/ping", s.authenticated(s.handlePing)).Methods(http.MethodPost)
	s.HandleFunc("/ip", s.authenticated(s.handleGetIp)).Methods(http.MethodGet)
	return
}

func (s *Server) Start(pubkey *ecdsa.PublicKey) error {
	s.pubKey = pubkey

	s.log = logger{}

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
	s.ip.Store(r.RemoteAddr)

	w.Write([]byte("OK\n"))
}

func (s *Server) handleGetIp(w http.ResponseWriter, r *http.Request)

func (s *Server) authenticated(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.authenticate(r) != nil {
			b, _ := json.Marshal(r.Header)
			s.log.Write("FAILED REQUEST", string(b), r.RemoteAddr)
			http.Error(w, "fuck off", http.StatusUnauthorized)
			return
		}
		f(w, r)
	}
}

func (s *Server) authenticate(r *http.Request) error {
	un, pw, ok := r.BasicAuth()
	if !ok || !s.validateTimestamp(r) {
		return ErrInvalid
	}
	authString := buildAuthString(r.Header.Get("request-timestamp"), un, pw)
	hash := sha256.Sum256([]byte(authString))
	signature, _ := base64.StdEncoding.DecodeString(r.Header.Get("request-signature"))
	if !ecdsa.VerifyASN1(s.pubKey, hash[:], []byte(signature)) {
		return ErrInvalid
	}
	return nil
}

func buildAuthString(timestamp, username, password string) string {
	return fmt.Sprintf("%s$%s$%s\n", timestamp, username, password)
}

func (s *Server) validateTimestamp(r *http.Request) bool {
	ts, _ := strconv.Atoi(r.Header.Get("request-timestamp"))

	return abs(ts-int(time.Now().UnixMilli())) < 500
}

func abs[T int | int8 | int16 | int32 | int64 | float32 | float64](v T) T {
	if v < 0 {
		return -v
	}
	return v
}
