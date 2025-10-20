package main

import (
	"fmt"
	"net/http"
	"os"
	"player/backend/internal/routes"
	"player/backend/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	mgr := server.NewManager()
	mm := server.NewMatchmaker()

	// Postgres config
	pgdsn := os.Getenv("PG_DSN")
	if pgdsn == "" {
		pgdsn = "postgres://postgres:Hemanth@localhost:5432/fourinarow?sslmode=disable"
	}
	pgstore, err := server.NewPGStore(pgdsn)
	if err != nil {
		panic("Failed to connect to Postgres: " + err.Error())
	}
	if err := pgstore.InitSchema(); err != nil {
		panic("Failed to init Postgres schema: " + err.Error())
	}

	// Kafka config
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	if brokers[0] == "" {
		brokers = []string{"localhost:9092"}
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "analytics"
	}
	kprod := server.NewKafkaProducer(brokers, topic)

	ws := server.NewWSHandler(mgr, mm, pgstore, kprod)

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "4 in a Row backend is running 🚀")
	})

	router.GET("/ws", func(c *gin.Context) { ws.Handle(c) })
	// serve static frontend if present
	router.Static("/static", "./static")
	// Register correct leaderboard route
	routes.RegisterRoutes(router, pgstore)

	fmt.Println("Server running on http://localhost:8080")
	router.Run(":8080")
}
