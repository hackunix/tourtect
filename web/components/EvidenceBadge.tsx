import styles from "./community.module.css";

const labels: Record<string, string> = { none: "Chưa có bằng chứng", metadata: "Có metadata", verified_receipt: "Đã kiểm tra biên nhận", verified_source: "Nguồn đã xác minh" };
export default function EvidenceBadge({ level }: { level: string }) { const verified = level.startsWith("verified"); return <span className={`${styles.badge} ${verified ? styles.evidence : ""}`} aria-label={`Evidence: ${level}`}><span aria-hidden>{verified ? "✓" : "i"}</span>{labels[level] || level}</span>; }
