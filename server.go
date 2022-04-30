package github_webhook

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"strings"
)

type HandlerFunc func(payload map[string]interface{}) error

type Server struct {
	bindAddr         string
	urlPath          string
	secretEnvVarName string
	maxPayloadSize   int
	handlerFunc      HandlerFunc
}

func NewServer(bindAddr, urlPath, secretEnvVarName string, maxPayloadSize int, handlerFunc HandlerFunc) (*Server, error) {
	if secretEnvVarName != "" {
		val, ok := os.LookupEnv(secretEnvVarName)
		if !ok {
			return nil, fmt.Errorf("environment variable [%s] not found", secretEnvVarName)
		}
		if val == "" {
			return nil, fmt.Errorf("empty secret")
		}
	}
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	return &Server{
		bindAddr:         bindAddr,
		urlPath:          urlPath,
		secretEnvVarName: secretEnvVarName,
		maxPayloadSize:   maxPayloadSize,
		handlerFunc:      handlerFunc,
	}, nil
}

func (s *Server) Serve() error {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post(s.urlPath, newRequestHandler(s.handlerFunc, s.secretEnvVarName, s.maxPayloadSize))
	log.Printf("starting github webhook server, listening on: %s, url path: %s\n", s.bindAddr, s.urlPath)
	if err := http.ListenAndServe(s.bindAddr, router); err != nil {
		return err
	}
	log.Printf("exiting github webhook server\n")
	return nil
}
