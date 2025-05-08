package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tiru-r/betengine/api"
	"github.com/tiru-r/betengine/service"
	"github.com/tiru-r/betengine/storage"
)

func main() {
	logger := log.New(os.Stdout, "bet-engine: ", log.LstdFlags)
	store := storage.NewStore(logger)
	svc := service.NewBetService(store, logger)
	apiHandler := api.NewHandler(svc, logger)

	// Set up routes
	http.HandleFunc("/api/bets/place", apiHandler.HandlePlaceBet)
	http.HandleFunc("/api/bets/settle", apiHandler.HandleSettleBet)
	http.HandleFunc("/api/balance", apiHandler.HandleGetBalance)

	// Start server
	logger.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
