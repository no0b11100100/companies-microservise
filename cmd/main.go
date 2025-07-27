package main

// @title Company API
// @version 1.0
// @description REST API for managing companies
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	app := NewApp("./cmd/cfg/config.yml")
	app.Run()
}
