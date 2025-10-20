package routes

import (
	"player/backend/internal/server"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, pg *server.PGStore) {
	r.GET("/leaderboard", func(c *gin.Context) {
		leaders, err := pg.GetLeaderboard()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		// Exclude bot from leaderboard
		arr := []interface{}{}
		for _, l := range leaders {
			if l.Username == "bot" {
				continue
			}
			arr = append(arr, l)
		}
		c.JSON(200, arr)
	})
	// ...other routes
}
