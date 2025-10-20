package server

import (
	"player/backend/internal/game"
	"sync"
	"time"
)

type Matchmaker struct {
	mu      sync.Mutex
	waiting map[string]*Session
}

type Session struct {
	Username string
	JoinedAt time.Time
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{waiting: make(map[string]*Session)}
}

// AddWaiting adds a user to waiting pool and returns after timeout a match decision.
func (m *Matchmaker) AddWaiting(username string, wait time.Duration) (*game.Game, bool, string) {
	m.mu.Lock()
	// if someone else waiting, match
	for other := range m.waiting {
		if other != username {
			delete(m.waiting, other)
			m.mu.Unlock()
			g := game.NewGame()
			return g, false, other
		}
	}
	// otherwise add self and wait
	m.waiting[username] = &Session{Username: username, JoinedAt: time.Now()}
	m.mu.Unlock()

	timer := time.NewTimer(wait)
	<-timer.C

	m.mu.Lock()
	defer m.mu.Unlock()
	// if still waiting -> remove and return bot game
	if _, ok := m.waiting[username]; ok {
		delete(m.waiting, username)
		g := game.NewGame()
		return g, true, ""
	}
	return nil, false, ""
}
