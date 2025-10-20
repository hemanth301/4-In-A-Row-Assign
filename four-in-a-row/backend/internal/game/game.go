package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	Cols = 7
	Rows = 6
)

type Player int

type Game struct {
	ID       string          `json:"id"`
	Board    [Rows][Cols]int `json:"board"`
	Turn     int             `json:"turn"` // which player (1 or 2)
	Finished bool            `json:"finished"`
	Winner   int             `json:"winner"` // 0 none, 1 or 2
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGame() *Game {
	return &Game{ID: fmt.Sprintf("g-%d", rand.Int63()), Turn: 1}
}

var ErrInvalidColumn = errors.New("invalid column")
var ErrColumnFull = errors.New("column full")
var ErrNotYourTurn = errors.New("not your turn")

// Drop attempts to drop a disc for player into column. Returns row index.
func (g *Game) Drop(column int, player int) (int, error) {
	if column < 0 || column >= Cols {
		return -1, ErrInvalidColumn
	}
	if player != g.Turn {
		return -1, ErrNotYourTurn
	}
	for r := 0; r < Rows; r++ {
		if g.Board[r][column] == 0 {
			g.Board[r][column] = player
			// toggle turn
			if !g.Finished {
				if player == 1 {
					g.Turn = 2
				} else {
					g.Turn = 1
				}
			}
			return r, nil
		}
	}
	return -1, ErrColumnFull
}

func (g *Game) IsFull() bool {
	for c := 0; c < Cols; c++ {
		if g.Board[Rows-1][c] == 0 {
			return false
		}
	}
	return true
}

// CheckWin checks whether placing at (r,c) for player produced a connect4.
func (g *Game) CheckWin(r, c int, player int) bool {
	if player == 0 {
		return false
	}
	dirs := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	for _, d := range dirs {
		cnt := 1
		// forward
		rr, cc := r+d[0], c+d[1]
		for rr >= 0 && rr < Rows && cc >= 0 && cc < Cols && g.Board[rr][cc] == player {
			cnt++
			rr += d[0]
			cc += d[1]
		}
		// backward
		rr, cc = r-d[0], c-d[1]
		for rr >= 0 && rr < Rows && cc >= 0 && cc < Cols && g.Board[rr][cc] == player {
			cnt++
			rr -= d[0]
			cc -= d[1]
		}
		if cnt >= 4 {
			return true
		}
	}
	return false
}

func (g *Game) String() string {
	s := ""
	for r := Rows - 1; r >= 0; r-- {
		for c := 0; c < Cols; c++ {
			s += fmt.Sprintf("%d", g.Board[r][c])
		}
		s += "\n"
	}
	return s
}
