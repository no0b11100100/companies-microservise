package server

import (
	"companies/cmd/internal/auth"
	configparser "companies/cmd/internal/configParser"
	"companies/cmd/internal/consts"
	"companies/cmd/internal/database"
	eventsender "companies/cmd/internal/eventSender"
	"companies/cmd/internal/server/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type RESTfulServer struct {
	router *chi.Mux
	addr   string
	port   string
}

type RESTServer interface {
	Serve()
	Shutdown()
}

func NewRESTfulServer(config configparser.HTTP, db database.Database, eventSender eventsender.EventSender) RESTServer {
	addr := configparser.GetCfgValue("HTTP_HOST", config.Addr)
	port := configparser.GetCfgValue("HTTP_PORT", config.Port)

	server := &RESTfulServer{addr: addr, port: port}

	server.router = chi.NewRouter()
	server.router.Use(middleware.Logger)

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

	s.router.Route("/api/v1/companies", func(r chi.Router) {
		r.With(auth.JWTMiddleware).Post("/", create)
		r.With(auth.JWTMiddleware).Patch("/{id}", update)
		r.With(auth.JWTMiddleware).Delete("/{id}", delete)
		r.Get("/{id}", get)
	})
}

func (s *RESTfulServer) Serve() {
	log.Println(consts.ApplicationPrefix, "Starting RESTful server")
	http.ListenAndServe(fmt.Sprintf("%v:%v", s.addr, s.port), s.router)
}

func (s *RESTfulServer) Shutdown() {}
