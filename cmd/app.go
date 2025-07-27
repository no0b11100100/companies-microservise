package main

import (
	configparser "companies/cmd/internal/configParser"
	"companies/cmd/internal/consts"
	"companies/cmd/internal/database"
	eventsender "companies/cmd/internal/eventSender"
	"companies/cmd/internal/server"
	"errors"
	"io"
	"log"
)

type app struct {
	db          io.Closer
	restServer  server.RESTServer
	eventSender io.Closer
}

func NewApp(configPath string) *app {
	log.Println(consts.ApplicationPrefix, "Starting app")

	config, err := configparser.LoadConfig(configPath)

	if err != nil {
		log.Println(consts.ApplicationPrefix, "Failed to load config", err.Error())
	}

	db := database.NewMySQLDB(config.DB)

	eventSender := eventsender.NewEventSender(config.Kafka)

	restServer := server.NewRESTfulServer(config.HTTP, db, eventSender)

	return &app{db, restServer, nil}
}

func (a *app) Run() {
	a.restServer.Serve()
}

func (a *app) Close() error {
	log.Println(consts.ApplicationPrefix, "Shutting down application")
	errorString := ""
	errorString += a.db.Close().Error()
	errorString += a.restServer.Shutdown().Error()
	errorString += a.eventSender.Close().Error()

	return errors.New(errorString)
}
