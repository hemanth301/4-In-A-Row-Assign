package server

import (
	"player/backend/internal/game"
)

// Bot tries to win, else block, else pick first available column
func BotMove(g *game.Game, botPlayer int) int {
	// 1) try winning move
	for c := 0; c < game.Cols; c++ {
		r := findRowForCol(&g.Board, c)
		if r >= 0 {
			// simulate
			g.Board[r][c] = botPlayer
			if g.CheckWin(r, c, botPlayer) {
				g.Board[r][c] = 0
				return c
			}
			g.Board[r][c] = 0
		}
	}
	// 2) block opponent
	opponent := 1
	if botPlayer == 1 {
		opponent = 2
	}
	for c := 0; c < game.Cols; c++ {
		r := findRowForCol(&g.Board, c)
		if r >= 0 {
			g.Board[r][c] = opponent
			if g.CheckWin(r, c, opponent) {
				g.Board[r][c] = 0
				return c
			}
			g.Board[r][c] = 0
		}
	}
	// 3) prefer center columns
	order := []int{3, 2, 4, 1, 5, 0, 6}
	for _, c := range order {
		if findRowForCol(&g.Board, c) >= 0 {
			return c
		}
	}
	return 0
}

func findRowForCol(b *[game.Rows][game.Cols]int, col int) int {
	for r := 0; r < game.Rows; r++ {
		if b[r][col] == 0 {
			return r
		}
	}
	return -1
}
