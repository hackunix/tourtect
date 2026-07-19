import Link from "next/link";
import type { Dictionary } from "@/lib/i18n";
import styles from "./community.module.css";
import ThemeLocaleControls from "./ThemeLocaleControls";

const icons: Record<string, string> = { assistant: "✦", explore: "⌕", community: "◎", saved: "◇", notifications: "○", profile: "☺", live: "◉", lens: "▣", sos: "!" };

export default function AppNavigation({ dict }: { dict: Dictionary }) {
  const main = [
    ["assistant", dict.assistant, "/assistant"], ["explore", dict.explore, "/search"], ["community", dict.community, "/community"],
    ["saved", dict.saved, "/saved"], ["profile", dict.profile, "/profile"],
  ];
  return <>
    <aside className={styles.rail} aria-label="Tourtect app navigation">
      <Link className={styles.brand} href="/assistant" aria-label="Tourtect Assistant"><span className={styles.brandMark}>T</span><span className={styles.brandText}>Tourtect</span></Link>
      <nav className={styles.nav} aria-label="Primary navigation">{main.map(([key,label,href]) => <Link className={styles.navLink} href={href} key={key}><span className={styles.navIcon} aria-hidden>{icons[key]}</span><span className={styles.navText}>{label}</span></Link>)}</nav>
      <div className={styles.railUtility}>
        <Link className={styles.navLink} href="/assistant?mode=voice"><span className={styles.navIcon} aria-hidden>{icons.live}</span><span className={styles.navText}>{dict.live}</span></Link>
        <Link className={styles.navLink} href="/assistant?mode=lens"><span className={styles.navIcon} aria-hidden>{icons.lens}</span><span className={styles.navText}>{dict.lens}</span></Link>
        <Link className={styles.navLink} href="/notifications"><span className={styles.navIcon} aria-hidden>{icons.notifications}</span><span className={styles.navText}>{dict.notifications}</span></Link>
        <Link className={`${styles.navLink} ${styles.sos}`} href="/assistant?mode=sos"><span className={styles.navIcon} aria-hidden>{icons.sos}</span><span className={styles.navText}>{dict.sos}</span></Link>
        <ThemeLocaleControls compact />
      </div>
    </aside>
    <nav className={styles.bottomNav} aria-label="Mobile navigation">{main.map(([key,label,href]) => <Link className={styles.navLink} href={href} key={key}><span className={styles.navIcon} aria-hidden>{icons[key]}</span><span>{label}</span></Link>)}</nav>
  </>;
}
