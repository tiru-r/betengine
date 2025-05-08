package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tiru-r/betengine/service"
)

// manages API endpoints
type Handler struct {
	svc    *service.BetService
	logger *log.Logger
}

// creates a new API handler
func NewHandler(svc *service.BetService, logger *log.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

// Request structs
type placeBetRequest struct {
	UserID  string  `json:"user_id"`
	EventID string  `json:"event_id"`
	Odds    float64 `json:"odds"`
	Amount  float64 `json:"amount"`
}

type settleBetRequest struct {
	EventID string `json:"event_id"`
	Result  string `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// handles bet placement
func (h *Handler) HandlePlaceBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req placeBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	betID, err := h.svc.PlaceBet(req.UserID, req.EventID, req.Odds, req.Amount)
	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"bet_id": betID})
	w.WriteHeader(http.StatusCreated)
}

// handles bet settlement
func (h *Handler) HandleSettleBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req settleBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.SettleBet(req.EventID, req.Result); err != nil {
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Bet settled successfully"})
}

// handles balance queries
func (h *Handler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	balance, err := h.svc.GetBalance(userID)
	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}
