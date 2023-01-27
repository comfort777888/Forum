package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/config"
)

type Server struct {
	Srv             *http.Server
	notify          chan error
	shutdownTimeOut time.Duration
	// db              *sql.DB
	// config     *config.Config
}

func NewServer(conf *config.Config, router *http.ServeMux) *Server {
	server := &Server{
		Srv: &http.Server{
			Addr:           ":" + conf.Port,
			Handler:        router,
			ReadTimeout:    conf.ReadTimeout,
			WriteTimeout:   conf.WriteTimeout,
			MaxHeaderBytes: conf.MaxHeaderBytes,
		},
		notify:          make(chan error, 1),
		shutdownTimeOut: conf.ShutdownTimeOut,
		// db:              db,
	}

	server.start()
	return server
}

func (s *Server) start() {
	log.Printf("server has been initiated on http://localhost%v/\n", s.Srv.Addr)
	go func() {
		s.notify <- s.Srv.ListenAndServe()
		close(s.notify)
		fmt.Println("notify chan")
	}()
}

// Notify figure out how it catches signal
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown gracefully shutdowns server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeOut)
	defer cancel()
	// s.db.Close()
	defer log.Println("graceful shutdown")
	return s.Srv.Shutdown(ctx)
}
