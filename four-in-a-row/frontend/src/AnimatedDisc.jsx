import React from "react";
import styles from "./AnimatedDisc.module.css";

export default function AnimatedDisc({ color, animate }) {
  // Color: 1 = red, 2 = yellow
  const discColor = color === 1 ? styles.red : styles.yellow;
  return (
    <div
      className={
        styles.disc + " " + discColor + (animate ? " " + styles.animate : "")
      }
    />
  );
}
