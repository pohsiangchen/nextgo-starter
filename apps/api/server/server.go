package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"

	"apps/api/config"
	"apps/api/database"
	sqlcStore "apps/api/database/sqlc"
	"apps/api/domain/health"
	"apps/api/domain/user"
	"apps/api/middleware"
)

type Server struct {
	cfg       *config.Config
	db        *sql.DB
	store     *sqlcStore.Queries
	validator *validator.Validate

	router *chi.Mux
}

func NewServer() *Server {
	srv := &Server{
		cfg:    config.Get(),
		router: chi.NewRouter(),
	}

	// initialization
	srv.newDatabase()
	srv.newValidator()
	srv.initRoutes()

	return srv
}

func (s *Server) newDatabase() {
	if s.cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file or set in environment variable")
	}

	dsn := database.DataSourceName(int(s.cfg.Database.Port), s.cfg.Database.Host, s.cfg.Database.User, s.cfg.Database.Password, s.cfg.Database.Name, s.cfg.Database.SslMode)
	s.db = database.NewDB(s.cfg.Database.Driver, dsn)
	s.db.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
	s.db.SetMaxIdleConns(s.cfg.Database.MaxIdleConnections)
	s.db.SetConnMaxLifetime(s.cfg.Database.ConnectionsMaxLifeTime)

	s.store = sqlcStore.New(s.db)
}

func (s *Server) newValidator() {
	s.validator = validator.New(validator.WithRequiredStructEnabled())
}

func (s *Server) initRoutes() {
	s.router.Use(chiMiddleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RequestID("req_id", "X-Request-Id"))
	s.router.Use(middleware.Recovery)

	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Get("/livez", health.Get)

	// initialize user routes
	s.router.Route("/api/v1", func(r chi.Router) {
		s.initUser(r)
	})
}

func (s *Server) initUser(r chi.Router) {
	userService := user.NewUserServiceImpl(s.store)
	userCtrl, err := user.NewUserController(userService, s.validator)
	if err != nil {
		log.Fatalf("Error initializing user controller: %v", err)
	}
	user.RegisterRoutes(r, userCtrl)
}

func (s *Server) Start() {
	server := http.Server{
		Addr:        s.cfg.API.Host + ":" + s.cfg.API.Port,
		Handler:     s.router,
		IdleTimeout: s.cfg.IdleTimeout,
		// to avoid clients hold up a connection by being slow to write or read
		// see https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	go func() {
		log.Printf("Server has started - http://%v\n", server.Addr)
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
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), s.cfg.API.GracefulTimeout)
	defer shutdownRelease()

	// when `Shutdown()` is called, `ListenAndServe()` immediately return `ErrServerClosed`.
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
		// force the server to close if the graceful shutdown is unable to complete within the specified timeout
		server.Close()
	}
	s.closeResources()
	log.Println("Graceful shutdown complete.")
}

func (s *Server) closeResources() {
	_ = s.db.Close()
}
