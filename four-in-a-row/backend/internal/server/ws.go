package server

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"player/backend/internal/game"
	"player/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type WSHandler struct {
	mgr   *Manager
	mm    *Matchmaker
	pg    *PGStore
	kafka *KafkaProducer
	// conns: gameID -> username -> conn
	conns map[string]map[string]*websocket.Conn
	// disconnect timers: gameID -> username -> timer
	timers map[string]map[string]*time.Timer
	mu     sync.Mutex
}

func NewWSHandler(mgr *Manager, mm *Matchmaker, pg *PGStore, kafka *KafkaProducer) *WSHandler {
	return &WSHandler{
		mgr:    mgr,
		mm:     mm,
		pg:     pg,
		kafka:  kafka,
		conns:  make(map[string]map[string]*websocket.Conn),
		timers: make(map[string]map[string]*time.Timer),
	}
}

func (h *WSHandler) Handle(c *gin.Context) {
	username := c.Query("username")
	gameID := c.Query("gameID")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Try to find existing game for this user
	var g *game.Game
	var gid string
	var found bool
	if gameID != "" {
		g, found = h.mgr.Get(gameID)
		gid = gameID
	} else {
		g, gid, found = h.mgr.GetGameByPlayer(username)
	}

	if !found || g == nil || g.Finished {
		// Not found or finished, matchmake
		g, createdWithBot, other := h.mm.AddWaiting(username, 10*time.Second)
		if createdWithBot {
			h.mgr.Add(g, username, "bot")
			gid = g.ID
			// Emit game started event
			if h.kafka != nil {
				payload, _ := json.Marshal(map[string]interface{}{
					"game_id":   gid,
					"players":   h.mgr.GetPlayers(gid),
					"timestamp": time.Now().UTC(),
				})
				h.kafka.Emit(services.EventGameStarted, string(payload))
			}
			// Immediately send new game state to client after bot joins
			conn.WriteJSON(g)
			// If bot is player 2 and it's bot's turn, make bot move after 1s delay
			players := h.mgr.GetPlayers(gid)
			if len(players) == 2 && players[1] == "bot" && g.Turn == 2 && !g.Finished {
				go func() {
					time.Sleep(1 * time.Second)
					importBotMove := func() int {
						// Import BotMove from server/bot.go
						return BotMove(g, 2)
					}
					col := importBotMove()
					r, err := g.Drop(col, 2)
					if err == nil {
						// check win
						if g.CheckWin(r, col, 2) {
							g.Finished = true
							g.Winner = 2
						} else if g.IsFull() {
							g.Finished = true
							g.Winner = 0
						}
						h.mgr.Add(g, players...)
						// Broadcast bot move to all clients
						h.mu.Lock()
						for _, c := range h.conns[gid] {
							c.WriteJSON(g)
						}
						h.mu.Unlock()
					}
				}()
			}
		} else if other != "" {
			h.mgr.Add(g, username, other)
			gid = g.ID
			// Emit game started event
			if h.kafka != nil {
				payload, _ := json.Marshal(map[string]interface{}{
					"game_id":   gid,
					"players":   h.mgr.GetPlayers(gid),
					"timestamp": time.Now().UTC(),
				})
				h.kafka.Emit(services.EventGameStarted, string(payload))
			}
		}
	}

	// Register connection
	h.mu.Lock()
	if h.conns[gid] == nil {
		h.conns[gid] = make(map[string]*websocket.Conn)
	}
	h.conns[gid][username] = conn
	// Cancel disconnect timer if present
	if h.timers[gid] != nil && h.timers[gid][username] != nil {
		h.timers[gid][username].Stop()
		delete(h.timers[gid], username)
	}
	h.mu.Unlock()

	// Send initial state
	conn.WriteJSON(g)

	// Broadcast helper
	broadcast := func(g *game.Game) {
		h.mu.Lock()
		defer h.mu.Unlock()
		for _, c := range h.conns[gid] {
			c.WriteJSON(g)
		}
	}

	// read loop
	for {
		var msg struct {
			Action string `json:"action"`
			Column int    `json:"column"`
		}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		if msg.Action == "drop" && g != nil && !g.Finished {
			// Determine player number
			pnum := 1
			players := h.mgr.GetPlayers(gid)
			if len(players) == 2 && players[1] == username {
				pnum = 2
			}
			r, err := g.Drop(msg.Column, pnum)
			if err == nil {
				// Emit move event
				if h.kafka != nil {
					payload, _ := json.Marshal(map[string]interface{}{
						"game_id":   gid,
						"player":    username,
						"column":    msg.Column,
						"row":       r,
						"timestamp": time.Now().UTC(),
					})
					h.kafka.Emit(services.EventMoveMade, string(payload))
				}
				// check win
				if g.CheckWin(r, msg.Column, pnum) {
					g.Finished = true
					g.Winner = pnum
				} else if g.IsFull() {
					g.Finished = true
					g.Winner = 0
				}
				h.mgr.Add(g, players...)
				broadcast(g)
				// Persist completed game and update leaderboard
				if g.Finished {
					moves := "[]" // TODO: store moves if needed
					p1, p2 := players[0], ""
					if len(players) > 1 {
						p2 = players[1]
					}
					if h.pg != nil {
						h.pg.SaveGame(g.ID, p1, p2, g.Winner, moves)
						if g.Winner == 1 {
							h.pg.AddWin(p1)
						} else if g.Winner == 2 {
							h.pg.AddWin(p2)
						}
					}
					// Emit game finished event
					if h.kafka != nil {
						payload, _ := json.Marshal(map[string]interface{}{
							"game_id":   gid,
							"winner":    g.Winner,
							"players":   players,
							"timestamp": time.Now().UTC(),
						})
						h.kafka.Emit(services.EventGameFinished, string(payload))
					}
				} else {
					// If bot is player 2 and it's now bot's turn, make bot move after 1s delay
					if len(players) == 2 && players[1] == "bot" && g.Turn == 2 && !g.Finished {
						go func() {
							time.Sleep(1 * time.Second)
							col := BotMove(g, 2)
							r, err := g.Drop(col, 2)
							if err == nil {
								// check win
								if g.CheckWin(r, col, 2) {
									g.Finished = true
									g.Winner = 2
								} else if g.IsFull() {
									g.Finished = true
									g.Winner = 0
								}
								h.mgr.Add(g, players...)
								broadcast(g)
								// Persist completed game and update leaderboard
								if g.Finished {
									moves := "[]" // TODO: store moves if needed
									p1, p2 := players[0], ""
									if len(players) > 1 {
										p2 = players[1]
									}
									if h.pg != nil {
										h.pg.SaveGame(g.ID, p1, p2, g.Winner, moves)
										if g.Winner == 1 {
											h.pg.AddWin(p1)
										} else if g.Winner == 2 {
											h.pg.AddWin(p2)
										}
									}
									// Emit game finished event
									if h.kafka != nil {
										payload, _ := json.Marshal(map[string]interface{}{
											"game_id":   gid,
											"winner":    g.Winner,
											"players":   players,
											"timestamp": time.Now().UTC(),
										})
										h.kafka.Emit(services.EventGameFinished, string(payload))
									}
								}
							}
						}()
					}
				}
			} else {
				conn.WriteJSON(gin.H{"error": err.Error()})
			}
		}
	}

	// On disconnect, start 30s timer
	h.mu.Lock()
	if h.timers[gid] == nil {
		h.timers[gid] = make(map[string]*time.Timer)
	}
	h.timers[gid][username] = time.AfterFunc(30*time.Second, func() {
		// Forfeit if not reconnected
		h.mu.Lock()
		defer h.mu.Unlock()
		if h.conns[gid][username] == nil && g != nil && !g.Finished {
			g.Finished = true
			// Opponent wins
			players := h.mgr.GetPlayers(gid)
			winner := 1
			if len(players) == 2 && players[1] == username {
				winner = 1
			} else {
				winner = 2
			}
			g.Winner = winner
			h.mgr.Add(g, players...)
			for _, c := range h.conns[gid] {
				c.WriteJSON(g)
			}
		}
	})
	// Remove connection
	delete(h.conns[gid], username)
	h.mu.Unlock()
}
