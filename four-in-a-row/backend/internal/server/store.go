package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type Store struct {
	mu     sync.Mutex
	Path   string
	Leader map[string]int `json:"leader"`
}

func NewStore(path string) *Store {
	s := &Store{Path: path, Leader: make(map[string]int)}
	s.load()
	return s
}

func (s *Store) load() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		return
	}
	b, _ := ioutil.ReadFile(s.Path)
	_ = json.Unmarshal(b, s)
}

func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, _ := json.MarshalIndent(s, "", "  ")
	return ioutil.WriteFile(s.Path, b, 0644)
}

func (s *Store) AddWin(username string) error {
	s.mu.Lock()
	s.Leader[username]++
	s.mu.Unlock()
	return s.Save()
}

func (s *Store) Leaderboard() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make(map[string]int)
	for k, v := range s.Leader {
		copy[k] = v
	}
	return copy
}
