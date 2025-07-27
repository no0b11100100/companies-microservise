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

func main() {
	log.Println(consts.ApplicationPrefix, "Starting app")

	config, err := configparser.LoadConfig("cmd/cfg/config.yml")

	if err != nil {
		log.Println(consts.ApplicationPrefix, "Failed to load config", err.Error())
	}

	db := database.NewMySQLDB(config.DB /*, "root", "password", "db", "3306", "companiesdb"*/)

	eventSender := eventsender.NewEventSender(config.Kafka /*, "kafka:9092"*/)

	restServer := server.NewRESTfulServer(config.HTTP /*"0.0.0.0", "8080",*/, db, eventSender)

	restServer.Serve()

	restServer.Shutdown()

	// db.Close()
	// eventSender.Close()
}
