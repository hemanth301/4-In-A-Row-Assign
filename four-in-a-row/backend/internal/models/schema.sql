-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    wins INT DEFAULT 0,
    losses INT DEFAULT 0,
    games_played INT DEFAULT 0
);

-- Games table
CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    player1 VARCHAR(50) NOT NULL,
    player2 VARCHAR(50) NOT NULL,
    winner VARCHAR(50),
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    duration INTERVAL
);
