package websocket

import (
	"net/http"
	"player/backend/internal/services"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Sessions map[string]*services.Session
	Conns    map[string]*websocket.Conn // username -> conn
	Mu       sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Sessions: make(map[string]*services.Session),
		Conns:    make(map[string]*websocket.Conn),
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (h *Hub) HandleWS(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	h.Mu.Lock()
	h.Conns[username] = conn
	h.Mu.Unlock()

	// TODO: matchmaking, session join/create, reconnection, gameplay loop
	// For now, just echo
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}
