import React, { useEffect, useState } from "react";
import styles from "./AnimatedBackground.module.css";

// Helper to generate random properties for stars
function randomStarProps() {
  const size = Math.random() * 3 + 2;
  const orbit = Math.random() * 40 + 30;
  const duration = Math.random() * 10 + 10;
  const delay = Math.random() * 8;
  return {
    left: Math.random() * 100 + "vw",
    top: Math.random() * 100 + "vh",
    size,
    orbit,
    duration,
    delay
  };
}

export default function AnimatedBackground() {
  const [stars, setStars] = useState([]);

  useEffect(() => {
    // Generate 60 stars
    const arr = Array.from({ length: 60 }, () => randomStarProps());
    setStars(arr);
  }, []);

  return (
    <div className={styles["animated-bg"]}>
      {stars.map((star, i) => (
        <div
          key={i}
          className={styles.star}
          style={{
            left: star.left,
            top: star.top,
            width: star.size,
            height: star.size,
            animationDuration: star.duration + "s",
            animationDelay: star.delay + "s"
          }}
        />
      ))}
    </div>
  );
}
