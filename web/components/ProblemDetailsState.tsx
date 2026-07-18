import Link from "next/link";
import styles from "./community.module.css";
export default function ProblemDetailsState({ title="Dịch vụ tạm thời không khả dụng", detail, requestId }: { title?: string; detail: string; requestId?: string }) { return <section className={styles.problem} role="alert"><h2>{title}</h2><p>{detail}</p>{requestId&&<details><summary>Chi tiết kỹ thuật</summary><code>Request ID: {requestId}</code></details>}<Link className={styles.button} href="/">Thử lại</Link></section>; }
