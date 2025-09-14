package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vishalpatel08/bon-rewards-service/api"
	dbconnection "github.com/vishalpatel08/bon-rewards-service/dbConnection"
	"github.com/vishalpatel08/bon-rewards-service/service"
)

func main() {
	fmt.Println("Welcome to BON Reward Service")

	const dbPath = "rewards.db"
	const serverAddress = ":8080"

	repo, err := dbconnection.NewRepository(dbPath)
	if err != nil {
		log.Fatalf("FATAL: Could not initialize repository: %v", err)
	}
	log.Println("Successfully connected to the database.")

	rewardService := service.NewRewardService(repo)
	apiHandler := api.NewHandler(rewardService)
	router := api.SetupRouter(apiHandler)

	log.Printf("Starting server on http://localhost%s", serverAddress)
	if err := http.ListenAndServe(serverAddress, router); err != nil {
		log.Fatalf("FATAL: Could not start server: %v", err)
	}
}
