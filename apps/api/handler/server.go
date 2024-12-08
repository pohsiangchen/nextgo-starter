package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Server struct {
	router *chi.Mux
}

func NewServer() *Server {
	srv := &Server{
		router: chi.NewRouter(),
	}

	srv.initRoutes()

	return srv
}

func (s *Server) initRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}

func (s *Server) Start() {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      s.router,
		IdleTimeout:  time.Second * 30, // s.cfg.IdleTimeout,
		ReadTimeout:  time.Second * 10, // s.cfg.ReadTimeout,
		WriteTimeout: time.Minute,      // s.cfg.WriteTimeout,
	}

	go func() {
		log.Println("Server has started")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	// when a platform shutdowns your instance, it sends a SIGTERM or SIGINT signal
	// use `signal.Notify()` to relay incoming signals to the channel `sigChan`
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// set a buffer time for active connections to be processed instead of terminates all of them immediatelly
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	// when `Shutdown()` is called, `ListenAndServe()` immediately return `ErrServerClosed`.
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")

}
