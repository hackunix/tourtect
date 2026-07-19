import Link from "next/link";
import type { FeedMode } from "@/lib/api";
import type { Dictionary } from "@/lib/i18n";
import styles from "./community.module.css";

export default function FeedHeader({ mode, regionId, dict }: { mode: FeedMode; regionId?: string; dict: Dictionary }) {
  const tabs: [FeedMode,string][]=[["following",dict.following],["nearby",dict.nearby],["latest",dict.latest],["trending",dict.trending],["safety",dict.safety]];
  return <nav className={styles.tabs} aria-label="Feed filters">{tabs.map(([value,label])=><Link className={styles.tab} aria-current={mode===value?"page":undefined} key={value} href={`/community?mode=${value}${value==="nearby"?`&region_id=${regionId||"hanoi-hoan-kiem"}`:""}`}>{value==="safety"&&<span aria-hidden>⚠ </span>}{label}</Link>)}</nav>;
}
