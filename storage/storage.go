package storage

import (
	"fmt"
	"log"
	"sync"

	"github.com/tiru-r/betengine/models"
)

// Storage manages bets and users in memory
type Store struct {
	bets   map[string]models.Bet
	users  map[string]models.User
	mutex  sync.RWMutex
	logger *log.Logger
}

// NewStore creates a new in-memory storage
func NewStore(logger *log.Logger) *Store {
	return &Store{
		bets:   make(map[string]models.Bet),
		users:  make(map[string]models.User),
		logger: logger,
	}
}

// Storing a bet
func (s *Store) SaveBet(bet models.Bet) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.bets[bet.ID] = bet
	return nil
}

// Retrieving bets by event ID
func (s *Store) GetBetsByEventID(eventID string) []models.Bet {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var bets []models.Bet
	for _, bet := range s.bets {
		if bet.EventID == eventID {
			bets = append(bets, bet)
		}
	}
	return bets
}

// Storing a user
func (s *Store) SaveUser(user models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.users[user.ID] = user
	return nil
}

// Retrieving a user by ID
func (s *Store) GetUser(userID string) (models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return models.User{}, fmt.Errorf("user not found: %s", userID)
	}
	return user, nil
}
