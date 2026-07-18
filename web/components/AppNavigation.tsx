import Link from "next/link";
import type { Dictionary } from "@/lib/i18n";
import styles from "./community.module.css";
import ThemeLocaleControls from "./ThemeLocaleControls";

const icons: Record<string, string> = { home: "⌂", search: "⌕", create: "+", saved: "◇", notifications: "○", profile: "☺", sos: "!" };

export default function AppNavigation({ dict }: { dict: Dictionary }) {
  const main = [
    ["home", dict.home, "/"], ["search", dict.search, "/search"], ["create", dict.create, "/#composer"],
    ["saved", dict.saved, "/saved"], ["notifications", dict.notifications, "/notifications"], ["profile", dict.profile, "/profile"],
  ];
  return <>
    <aside className={styles.rail} aria-label="Tourtect app navigation">
      <Link className={styles.brand} href="/" aria-label="Tourtect home"><span className={styles.brandMark}>T</span><span className={styles.brandText}>Tourtect</span></Link>
      <nav className={styles.nav} aria-label="Primary navigation">{main.map(([key,label,href]) => <Link className={styles.navLink} href={href} key={key}><span className={styles.navIcon} aria-hidden>{icons[key]}</span><span className={styles.navText}>{label}</span></Link>)}</nav>
      <div className={styles.railUtility}>
        <Link className={`${styles.navLink} ${styles.sos}`} href="/#safety"><span className={styles.navIcon} aria-hidden>{icons.sos}</span><span className={styles.navText}>{dict.sos}</span></Link>
        <ThemeLocaleControls compact />
      </div>
    </aside>
    <nav className={styles.bottomNav} aria-label="Mobile navigation">{main.slice(0,5).map(([key,label,href]) => <Link className={styles.navLink} href={href} key={key}><span className={styles.navIcon} aria-hidden>{icons[key]}</span><span>{label}</span></Link>)}</nav>
  </>;
}
