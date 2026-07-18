import type { ReactNode } from "react";
import type { Dictionary } from "@/lib/i18n";
import AppNavigation from "./AppNavigation";
import styles from "./community.module.css";

export default function AppShell({ dict, children, context }: { dict: Dictionary; children: ReactNode; context?: ReactNode }) {
  return <div className={styles.shell}><AppNavigation dict={dict} />{children}<aside className={styles.context} aria-label="Traveler context">{context}</aside></div>;
}
