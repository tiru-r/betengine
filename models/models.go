package models

import "time"

// Bet represents a single bet
type Bet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"`
	Odds      float64   `json:"odds"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"` // "pending", "won", "lost"
}

// User represents a user with balance
type User struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}
