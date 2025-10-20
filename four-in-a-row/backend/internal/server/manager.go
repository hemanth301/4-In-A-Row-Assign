package server

import (
	"player/backend/internal/game"
	"sync"
)

type Manager struct {
	mu           sync.Mutex
	games        map[string]*game.Game
	playerToGame map[string]string   // username -> gameID
	gamePlayers  map[string][]string // gameID -> usernames
}

func NewManager() *Manager {
	return &Manager{
		games:        make(map[string]*game.Game),
		playerToGame: make(map[string]string),
		gamePlayers:  make(map[string][]string),
	}
}

func (m *Manager) Add(g *game.Game, players ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.games[g.ID] = g
	if len(players) > 0 {
		m.gamePlayers[g.ID] = players
		for _, p := range players {
			m.playerToGame[p] = g.ID
		}
	}
}

func (m *Manager) Get(id string) (*game.Game, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	g, ok := m.games[id]
	return g, ok
}

// GetGameByPlayer returns the game and gameID for a player if active
func (m *Manager) GetGameByPlayer(username string) (*game.Game, string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	gid, ok := m.playerToGame[username]
	if !ok {
		return nil, "", false
	}
	g, ok := m.games[gid]
	return g, gid, ok
}

// GetPlayers returns the usernames for a gameID
func (m *Manager) GetPlayers(gameID string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.gamePlayers[gameID]
}

func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.games, id)
	delete(m.gamePlayers, id)
	// remove playerToGame entries
	for p, gid := range m.playerToGame {
		if gid == id {
			delete(m.playerToGame, p)
		}
	}
}
