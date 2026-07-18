import type { PriceInsight as Insight } from "@/lib/api";
import styles from "./community.module.css";

const labels={typical:"Trong khoảng tham chiếu",elevated:"Cao hơn thông thường",high_risk:"Cao đáng kể",insufficient_data:"Chưa đủ dữ liệu"};
export default function PriceInsight({ insight }: { insight: Insight }) { return <section className={styles.priceInsight} data-level={insight.alert_level} aria-live="polite"><strong>{labels[insight.alert_level]}</strong><p>Đã nhập: {insight.observed.amount_minor} {insight.observed.currency}</p>{insight.reference&&<p>Khoảng tham chiếu: {insight.reference.p10_minor||"—"}–{insight.reference.p90_minor||"—"} {insight.reference.currency}</p>}<p>Độ tin cậy: {Math.round(insight.confidence*100)}% · Cỡ mẫu: {insight.reference?.effective_sample_size??"không có"}</p><p>{insight.reasons.join(" · ")||"Không có giải thích bổ sung"}</p><small>Dữ liệu {insight.dataset_version} · {new Date(insight.freshness).toLocaleDateString("vi-VN")}</small></section>; }
