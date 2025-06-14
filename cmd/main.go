package main

import (
	"github.com/golates/api-gateway/internal/api"
	"github.com/golates/api-gateway/pkg/config"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load app configuration
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	cfg := config.LoadConfig()

	// Setup chi API
	server := api.NewAPI(cfg)
	server.SetupMiddlewares()
	server.SetupRoutes()
	err = server.RunServer()
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
