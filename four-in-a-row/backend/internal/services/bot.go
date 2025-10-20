package services

import "player/backend/internal/game"

// BotMove returns the best column for the bot to play
func BotMove(board *[game.Rows][game.Cols]int, botPlayer int) int {
	// 1. Try to win
	for c := 0; c < game.Cols; c++ {
		copy := *board
		g := &game.Game{Board: copy}
		row, err := g.Drop(c, botPlayer)
		if err == nil && g.CheckWin(row, c, botPlayer) {
			return c
		}
	}
	// 2. Block opponent's win
	opponent := 1
	if botPlayer == 1 {
		opponent = 2
	}
	for c := 0; c < game.Cols; c++ {
		copy := *board
		g := &game.Game{Board: copy}
		row, err := g.Drop(c, opponent)
		if err == nil && g.CheckWin(row, c, opponent) {
			return c
		}
	}
	// 3. Prefer center columns
	order := []int{3, 2, 4, 1, 5, 0, 6}
	for _, c := range order {
		copy := *board
		g := &game.Game{Board: copy}
		if _, err := g.Drop(c, botPlayer); err == nil {
			return c
		}
	}
	// 4. Fallback: first available
	for c := 0; c < game.Cols; c++ {
		copy := *board
		g := &game.Game{Board: copy}
		if _, err := g.Drop(c, botPlayer); err == nil {
			return c
		}
	}
	return 0
}
