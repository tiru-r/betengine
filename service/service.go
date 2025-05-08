package service

import (
	"fmt"
	"log"
	"time"

	"github.com/tiru-r/betengine/models"
)

// Storage interface defines storage operations
type Storage interface {
	SaveBet(bet models.Bet) error
	GetBetsByEventID(eventID string) []models.Bet
	SaveUser(user models.User) error
	GetUser(userID string) (models.User, error)
}

// handles bet operations
type BetService struct {
	store  Storage
	logger *log.Logger
}

// creates a new bet service
func NewBetService(store Storage, logger *log.Logger) *BetService {
	return &BetService{
		store:  store,
		logger: logger,
	}
}

// places a new bet
func (s *BetService) PlaceBet(userID, eventID string, odds, amount float64) (string, error) {
	// Validate inputs
	if amount <= 0 || odds <= 0 {
		return "", fmt.Errorf("invalid amount or odds")
	}

	// Check user exists and has sufficient balance
	user, err := s.store.GetUser(userID)
	if err != nil {
		user = models.User{ID: userID, Balance: 1000.0} // Initialize with default balance
	}
	if user.Balance < amount {
		return "", fmt.Errorf("insufficient balance: %.2f", user.Balance)
	}

	// Create bet
	betID := fmt.Sprintf("%d", time.Now().UnixNano())
	bet := models.Bet{
		ID:        betID,
		UserID:    userID,
		EventID:   eventID,
		Odds:      odds,
		Amount:    amount,
		Timestamp: time.Now(),
		Status:    "pending",
	}

	// Update user balance
	user.Balance -= amount
	if err := s.store.SaveUser(user); err != nil {
		return "", fmt.Errorf("failed to save user: %v", err)
	}
	if err := s.store.SaveBet(bet); err != nil {
		return "", fmt.Errorf("failed to save bet: %v", err)
	}

	s.logger.Printf("Bet placed: user=%s, event=%s, amount=%.2f, odds=%.2f", userID, eventID, amount, odds)
	return betID, nil
}

// SettleBet settles a bet
func (s *BetService) SettleBet(eventID, result string) error {
	// validate result
	if result != "win" && result != "lose" {
		return fmt.Errorf("invalid result: %s", result)
	}

	bets := s.store.GetBetsByEventID(eventID)
	for _, bet := range bets {
		if bet.Status != "pending" {
			continue
		}

		user, err := s.store.GetUser(bet.UserID)
		if err != nil {
			return fmt.Errorf("user not found: %s", bet.UserID)
		}

		bet.Status = result
		if result == "win" {
			payout := bet.Amount * bet.Odds
			user.Balance += payout
			s.logger.Printf("Bet won: user=%s, event=%s, payout=%.2f", bet.UserID, eventID, payout)
		} else {
			s.logger.Printf("Bet lost: user=%s, event=%s", bet.UserID, eventID)
		}

		if err := s.store.SaveUser(user); err != nil {
			return fmt.Errorf("failed to save user: %v", err)
		}
		if err := s.store.SaveBet(bet); err != nil {
			return fmt.Errorf("failed to save bet: %v", err)
		}
	}
	return nil
}

// GetBalance returns user balance
func (s *BetService) GetBalance(userID string) (float64, error) {
	user, err := s.store.GetUser(userID)
	if err != nil {
		return 0, err
	}
	return user.Balance, nil
}
