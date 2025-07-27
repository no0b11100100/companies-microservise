package main

import (
	configparser "companies/cmd/internal/configParser"
	"companies/cmd/internal/consts"
	"companies/cmd/internal/database"
	eventsender "companies/cmd/internal/eventSender"
	"companies/cmd/internal/server"
	"log"
)

// docker-compose down --remove-orphans

// @title Company API
// @version 1.0
// @description REST API for managing companies
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	log.Println(consts.ApplicationPrefix, "Starting app")

	config, err := configparser.LoadConfig("cmd/cfg/config.yml")

	if err != nil {
		log.Println(consts.ApplicationPrefix, "Failed to load config", err.Error())
	}

	db := database.NewMySQLDB(config.DB)

	eventSender := eventsender.NewEventSender(config.Kafka)

	restServer := server.NewRESTfulServer(config.HTTP, db, eventSender)

	restServer.Serve()

	restServer.Shutdown()

	// db.Close()
	// eventSender.Close()
}
