package main

import (
	"fmt"
	"log"
	"net/http"
	"thirdparty-service/database/mongodb"
	"thirdparty-service/environment"

	srv "thirdparty-service/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load content of .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env ")
	}

	cfg := environment.LoadConfig()

	// Get mongodb instance
	store, _, err := mongodb.New(cfg.DatabaseURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal("failed to establish MongoDB connection")
	}

	addr := fmt.Sprintf(":%s", cfg.PORT)
	router := srv.MountServer(cfg, store)
	// start HTTP server
	fmt.Println(fmt.Sprintf("starting HTTP service running on port %v", addr))
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("error starting http server")
	}
}
