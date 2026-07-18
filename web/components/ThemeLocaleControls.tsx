"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import styles from "./community.module.css";

export default function ThemeLocaleControls({ compact = false }: { compact?: boolean }) {
  const [dark, setDark] = useState(false); const router = useRouter();
  const toggleTheme = () => { const current = document.documentElement.dataset.theme === "dark" || (!document.documentElement.dataset.theme && matchMedia("(prefers-color-scheme: dark)").matches); const next = !current; setDark(next); document.documentElement.dataset.theme = next ? "dark" : "light"; localStorage.setItem("tourtect_theme", next ? "dark" : "light"); };
  const toggleLocale = () => { const current = document.documentElement.lang; const next = current.startsWith("vi") ? "en" : "vi-VN"; document.cookie = `tourtect_locale=${next}; path=/; max-age=31536000; samesite=lax`; router.refresh(); };
  return <div className={styles.nav}>
    <button className={styles.navButton} type="button" onClick={toggleTheme} aria-label={dark ? "Use light theme" : "Use dark theme"}><span className={styles.navIcon} aria-hidden>{dark ? "☀" : "◐"}</span>{!compact && <span>{dark ? "Light" : "Dark"}</span>}</button>
    <button className={styles.navButton} type="button" onClick={toggleLocale} aria-label="Change interface language"><span className={styles.navIcon} aria-hidden>文</span>{!compact && <span>VI / EN</span>}</button>
  </div>;
}
