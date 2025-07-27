package server

import (
	"companies/cmd/internal/auth"
	configparser "companies/cmd/internal/configParser"
	"companies/cmd/internal/consts"
	"companies/cmd/internal/database"
	eventsender "companies/cmd/internal/eventSender"
	"companies/cmd/internal/server/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type RESTfulServer struct {
	router *chi.Mux
	addr   string
	port   string
	srv    *http.Server
}

//go:generate mockgen -source=server.go -destination=../../tests/mocks/mock_rest_server.go -package=mocks
type RESTServer interface {
	Serve()
	Shutdown() error
}

func NewRESTfulServer(config configparser.HTTP, db database.Database, eventSender eventsender.EventSender) RESTServer {
	addr := configparser.GetCfgValue("HTTP_HOST", config.Addr)
	port := configparser.GetCfgValue("HTTP_PORT", config.Port)

	server := &RESTfulServer{addr: addr, port: port}

	server.router = chi.NewRouter()

	// metrics.Init()

	server.router.Use(middleware.Logger)
	// server.router.Use(metrics.MetricsMiddleware)

	server.srv = &http.Server{
		Addr:              fmt.Sprintf("%v:%v", addr, port),
		Handler:           server.router,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1024 * 1024,
	}

	server.initHandlers(db, eventSender)

	return server
}

func (s *RESTfulServer) initHandlers(db database.Database, eventSender eventsender.EventSender) {
	create := handlers.NewCreateRecordHandler(db, eventSender)
	update := handlers.NewUpdateRecordHandler(db, eventSender)
	get := handlers.NewGetRecordHandler(db)
	delete := handlers.NewDeleteRecordHandler(db, eventSender)

	// just for test
	s.router.Post("/api/v1/token", auth.HandleFunc)

	s.router.Get("/swagger/*", httpSwagger.WrapHandler)

	s.router.Handle("/metrics", promhttp.Handler())

	s.router.Route("/api/v1/companies", func(r chi.Router) {
		r.With(auth.JWTMiddleware).Post("/", create)
		r.With(auth.JWTMiddleware).Patch("/{id}", update)
		r.With(auth.JWTMiddleware).Delete("/{id}", delete)
		r.Get("/{id}", get)
	})
}

func (s *RESTfulServer) Serve() {
	log.Println(consts.ApplicationPrefix, "Starting RESTful server")
	s.srv.ListenAndServe()
}

func (s *RESTfulServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
