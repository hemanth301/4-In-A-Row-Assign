import React, { useState, useRef, useEffect } from "react";
import styles from "./App.module.css";
import AnimatedBackground from "./AnimatedBackground";
import AnimatedDisc from "./AnimatedDisc";

const COLS = 7;
const ROWS = 6;

export default function App() {
  const [username, setUsername] = useState("");
  const [connected, setConnected] = useState(false);
  const [game, setGame] = useState(null);
  const [error, setError] = useState("");
  const [leaderboard, setLeaderboard] = useState(null);
  const wsRef = useRef(null);
  const [animDrop, setAnimDrop] = useState({});

  // Cookie helpers
  function setCookie(name, value, days) {
    let expires = "";
    if (days) {
      const date = new Date();
      date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
      expires = "; expires=" + date.toUTCString();
    }
    document.cookie = `${name}=${value || ""}${expires}; path=/`;
  }

  function getCookie(name) {
    const nameEQ = name + "=";
    const ca = document.cookie.split(";");
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i].trim();
      if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
  }

  useEffect(() => {
    const saved = getCookie("username");
    if (saved) setUsername(saved);
  }, []);

  function connect() {
    setError("");
    setCookie("username", username, 30);

    const ws = new WebSocket(
      `ws://${window.location.hostname}:8080/ws?username=${encodeURIComponent(
        username
      )}`
    );

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onmessage = (ev) => {
      try {
        const data = JSON.parse(ev.data);
        if (data.error) setError(data.error);
        else {
          setGame(data);
          setAnimDrop({});
        }
      } catch (e) {
        console.error("WebSocket parse error:", e);
      }
    };

    wsRef.current = ws;
  }

  function drop(col) {
    if (!wsRef.current || !game || game.finished) return;

    wsRef.current.send(JSON.stringify({ action: "drop", column: col }));

    // Animate disc drop
    const row = game.board.findIndex((r) => r[col] === 0);
    if (row !== -1) {
      setAnimDrop({ col, row });
      setTimeout(() => setAnimDrop({}), 500);
    }
  }

  function fetchLeaderboard() {
    fetch("/leaderboard")
      .then(async (r) => {
        const text = await r.text();
        try {
          const data = JSON.parse(text);
          if (Array.isArray(data)) {
            setLeaderboard(data);
          } else {
            setError("Leaderboard data is not valid.");
          }
        } catch (err) {
          setError("Leaderboard fetch error: Invalid server response.");
          console.error("Leaderboard fetch error:", err, text);
        }
      })
      .catch((err) => {
        setError("Leaderboard fetch error: " + err.message);
        console.error("Leaderboard fetch error:", err);
      });
  }

  return (
    <>
      {/* Animated background with extra moving objects */}
      <AnimatedBackground />
      {/* Main game UI */}
      <div className={styles.container}>
        <div className={styles.header}>4 in a Row</div>

        {/* Username Input and Connect Button */}
        <div className={styles.usernameBar}>
          <input
            className={styles.input}
            placeholder="Enter username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            disabled={connected}
          />
          {!connected ? (
            <button
              className={styles.button}
              onClick={connect}
              disabled={!username}
            >
              Connect
            </button>
          ) : (
            <button
              className={styles.button}
              onClick={() => wsRef.current && wsRef.current.close()}
            >
              Disconnect
            </button>
          )}
        </div>

        {/* Error Message */}
        {error && <div className={styles.error}>{error}</div>}

        {/* Connection Status */}
        {connected && (
          <div className={styles.status}>
            Connected as <b>{username}</b>
          </div>
        )}

        {/* Game Board */}
        {game && (
          <>
            <div className={styles.gameInfo}>
              Game ID: <b>{game.id}</b>
            </div>
            <div className={styles.board}>
              {game.board
                .slice()
                .reverse()
                .map((row, rIdx) =>
                  row.map((cell, cIdx) => {
                    const boardRow = ROWS - 1 - rIdx;
                    const isAnim =
                      animDrop.col === cIdx && animDrop.row === boardRow;

                    return (
                      <div
                        key={`${rIdx}-${cIdx}`}
                        className={styles.cell}
                        onClick={() =>
                          !game.finished && cell === 0 && drop(cIdx)
                        }
                      >
                        {cell !== 0 && (
                          <AnimatedDisc color={cell} animate={isAnim} />
                        )}
                      </div>
                    );
                  })
                )}
            </div>

            {/* Game Status */}
            {game.finished && (
              <div className={styles.status}>
                {game.winner === 0
                  ? "Draw!"
                  : game.winner === 1
                  ? "Red wins!"
                  : "Yellow wins!"}
              </div>
            )}
          </>
        )}

        {/* Leaderboard */}
        <div className={styles.leaderboard}>
          <button className={styles.button} onClick={fetchLeaderboard}>
            Show Leaderboard
          </button>
          {leaderboard && Array.isArray(leaderboard) && (
            <table className={styles.leaderboardTable}>
              <thead>
                <tr>
                  <th>Rank</th>
                  <th>Player</th>
                  <th>Wins</th>
                </tr>
              </thead>
              <tbody>
                {leaderboard.map((entry, idx) => (
                  <tr key={entry.username}>
                    <td>{idx + 1}</td>
                    <td>{entry.username}</td>
                    <td>{entry.wins}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </>
  );
}
