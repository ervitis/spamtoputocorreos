package spamtoputocorreos

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

type (
	Server struct {
		port string
		mux  *http.ServeMux
	}
)

func NewServer(port string) *Server {
	return &Server{
		port: port,
		mux:  http.NewServeMux(),
	}
}

func (s *Server) Start() {
	const timeout = time.Second * 15

	s.mux.HandleFunc("/health", handlerHealthCheck)

	server := &http.Server{
		Addr:              ":" + s.port,
		Handler:           s.mux,
		ReadTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       50 * time.Second,
	}

	go func(server *http.Server) {
		log.Panic(server.ListenAndServe())
	}(server)

	<-GlobalSignalHandler

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
	log.Println("Shutting down server")
	os.Exit(0)
}

func handlerHealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
