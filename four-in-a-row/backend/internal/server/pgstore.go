package server

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PGStore struct {
	db *sql.DB
}

func NewPGStore(dsn string) (*PGStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PGStore{db: db}, nil
}

func (s *PGStore) InitSchema() error {
	schema := `
CREATE TABLE IF NOT EXISTS games (
	id TEXT PRIMARY KEY,
	player1 TEXT,
	player2 TEXT,
	winner INT,
	moves JSONB,
	created_at TIMESTAMP DEFAULT now()
);
CREATE TABLE IF NOT EXISTS leaderboard (
	username TEXT PRIMARY KEY,
	wins INT
);
`
	_, err := s.db.Exec(schema)
	return err
}

func (s *PGStore) SaveGame(id, p1, p2 string, winner int, moves string) error {
	_, err := s.db.Exec(`INSERT INTO games (id, player1, player2, winner, moves) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (id) DO NOTHING`, id, p1, p2, winner, moves)
	return err
}

func (s *PGStore) AddWin(username string) error {
	_, err := s.db.Exec(`INSERT INTO leaderboard (username, wins) VALUES ($1,1) ON CONFLICT (username) DO UPDATE SET wins=leaderboard.wins+1`, username)
	return err
}

type Leader struct {
	Username string `json:"username"`
	Wins     int    `json:"wins"`
}

func (s *PGStore) GetLeaderboard() ([]Leader, error) {
	rows, err := s.db.Query(`SELECT username, wins FROM leaderboard ORDER BY wins DESC LIMIT 20`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []Leader{}
	for rows.Next() {
		var l Leader
		if err := rows.Scan(&l.Username, &l.Wins); err != nil {
			// Log the error and continue to the next row
			// This prevents one bad row from failing the whole operation
			continue
		}
		res = append(res, l)
	}
	return res, nil
}
