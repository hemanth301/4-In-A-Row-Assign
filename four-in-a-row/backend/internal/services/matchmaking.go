package services

import (
	"player/backend/internal/models"
	"sync"
	"time"
)

type Matchmaker struct {
	mu      sync.Mutex
	waiting map[string]*models.Player
}

type Session struct {
	GameID      string
	Players     []*models.Player
	CreatedAt   time.Time
	LastActive  time.Time
	BoardState  [6][7]int
	CurrentTurn int
	Status      string // waiting, active, finished
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{waiting: make(map[string]*models.Player)}
}

// AddPlayer adds a player to the matchmaking queue
func (m *Matchmaker) AddPlayer(p *models.Player) (matched *Session, withBot bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, other := range m.waiting {
		if other.Username != p.Username {
			// Found a match
			delete(m.waiting, other.Username)
			return &Session{
				GameID:     generateGameID(),
				Players:    []*models.Player{other, p},
				CreatedAt:  time.Now(),
				LastActive: time.Now(),
				Status:     "active",
			}, false
		}
	}
	// No match, add to waiting
	m.waiting[p.Username] = p
	return nil, false
}

func generateGameID() string {
	return time.Now().Format("20060102150405")
}
