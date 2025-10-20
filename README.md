# 4-In-A-Row-Assign
🎮 Four-in-a-Row

A full-stack multiplayer **Connect Four** game with real-time gameplay, leaderboard tracking, analytics, and bot support.

---

## 🧭 Table of Contents

- [Project Overview](#project-overview)
- [Tech Stack](#tech-stack)
- [Directory Structure](#directory-structure)
- [Backend Modules & Functionality](#backend-modules--functionality)
- [Frontend Modules & Functionality](#frontend-modules--functionality)
- [Database Schema](#database-schema)
- [Kafka Analytics](#kafka-analytics)
- [How to Run](#how-to-run)
- [API Endpoints](#api-endpoints)
- [License](#license)

---

## 🧩 Project Overview

This project implements a classic **Four-in-a-Row (Connect Four)** game with the following features:

- 🎯 Real-time multiplayer via **WebSockets**
- 🤖 **Bot opponent** if no player is available
- 🏆 Persistent **leaderboard** using PostgreSQL
- 📊 **Game analytics** emitted via Kafka
- 💫 **Modern React frontend** with animations

---

## ⚙️ Tech Stack

**Backend:**
- Go (Gin, Gorilla WebSocket)
- PostgreSQL (via Docker)
- Kafka (via Docker)
- kafka-go library

**Frontend:**
- React (with Vite)
- CSS Modules for styling

**DevOps:**
- Docker Compose for service orchestration

---

## 📂 Directory Structure

four-in-a-row/
backend/
cmd/server/ # Main server entrypoint
internal/
game/ # Game logic (board, moves, win check)
models/ # Data models (Player, SQL schema)
routes/ # API routes (leaderboard, etc.)
server/ # Server modules (WebSocket, matchmaking, PGStore, bot, manager)
services/ # Analytics, migrations, matchmaking, bot logic
websocket/ # (Legacy/experimental WebSocket handler)
static/ # Minimal static frontend (for testing)
frontend/
src/ # React components (App, AnimatedDisc, AnimatedBackground)
index.html # Frontend entrypoint
package.json # Frontend dependencies/scripts
vite.config.js # Vite dev server config
docker-compose.yml # Orchestrates Postgres, Kafka, Zookeeper
go.mod # Go dependencies

markdown
Copy code

---

## 🖥️ Backend Modules & Functionality

- **Game Logic** (`internal/game/game.go`)  
  Handles board state, move validation, and win/draw detection.

- **Matchmaking** (`internal/server/matchmaker.go`)  
  Pairs players or assigns a bot opponent if no match is found within a timeout.

- **Bot** (`internal/server/bot.go`)  
  Implements a simple AI: attempts to win, block, or pick the best column.

- **Manager** (`internal/server/manager.go`)  
  Tracks active games and player-to-game mapping.

- **WebSocket Handler** (`internal/server/ws.go`)  
  Manages real-time game state updates, moves, reconnections, and bot turns.

- **Postgres Store** (`internal/server/pgstore.go`)  
  Persists game data and provides leaderboard queries.

- **Kafka Producer** (`internal/server/kafka.go`)  
  Emits analytics events (`game_started`, `move_made`, `game_finished`).

- **API Routes** (`internal/routes/api.go`)  
  Defines REST endpoints such as `/leaderboard`.

---

## 🎨 Frontend Modules & Functionality

- **App.jsx** (`frontend/src/App.jsx`)  
  Main UI component: username input, connection management, game board, and leaderboard display.  
  Handles WebSocket communication and game updates.

- **AnimatedDisc.jsx** (`frontend/src/AnimatedDisc.jsx`)  
  Displays animated game discs (red/yellow).

- **AnimatedBackground.jsx** (`frontend/src/AnimatedBackground.jsx`)  
  Renders dynamic animated backgrounds for visual appeal.

- **CSS Modules** (`frontend/src/*.module.css`)  
  Handles all styling and animations for a smooth, responsive UI.

---

## 🗃️ Database Schema

- **Users Table:** Tracks username, wins, losses, and total games played.  
- **Games Table:** Records all games, players, winners, and timestamps.  
- **Leaderboard Table:** Stores player statistics for ranking.

See [`internal/models/schema.sql`](internal/models/schema.sql) for implementation details.

---

## 🔄 Kafka Analytics

- **Events Produced:**  
  - game_started 
  - move_made  
  - game_finished

- **Producer:** Emits analytics for each major game event.  
- **Consumer:** Can be extended for real-time dashboards or reporting.

---

## 🚀 How to Run

### 1. Start Backend Services
Make sure Docker is installed, then run:

docker-compose up
This starts PostgreSQL, Kafka, and Zookeeper.

2. Run the Go Backend
bash
Copy code
cd backend
go run cmd/server/main.go
Server runs on http://localhost:8080

3. Run the React Frontend
bash
Copy code
cd frontend
npm install
npm run dev
Frontend runs on http://localhost:5173
(Proxies /ws and /leaderboard requests to the backend.)

🌐 API Endpoints
Method	Endpoint	Description
GET	/ws?username=...	Opens a WebSocket for a game session
GET	/leaderboard	Returns top players (excluding bot)

📁 Key Source Files
internal/game/game.go

internal/server/ws.go

frontend/src/App.jsx

internal/routes/api.go

internal/server/pgstore.go
