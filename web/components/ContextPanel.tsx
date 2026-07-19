"use client";

import { FormEvent, useState } from "react";
import Link from "next/link";
import { api, PriceInsight as PriceResult, ProblemDetailError, SafetyAssessment } from "@/lib/api";
import PriceInsight from "./PriceInsight";
import SafetyAlert from "./SafetyAlert";
import styles from "./community.module.css";

export default function ContextPanel() {
  const [amount,setAmount]=useState("380000");const [price,setPrice]=useState<PriceResult|null>(null);const [safety,setSafety]=useState<SafetyAssessment|null>(null);const [error,setError]=useState("");const [busy,setBusy]=useState(false);
  const checkPrice=async(e:FormEvent)=>{e.preventDefault();setBusy(true);setError("");try{setPrice(await api.checkPrice({vertical:"taxi",raw_item:"Airport taxi to Old Quarter",money:{amount_minor:amount,currency:"VND",exponent:0},unit:"trip",region_id:"hanoi-soc-son",service_segment:"standard",venue_type:"transport_vendor",transaction_context:"metered",observed_at:new Date().toISOString(),user_confirmed:true}))}catch(e){setError(e instanceof ProblemDetailError?`${e.message} · ${e.requestId}`:"Không thể kiểm tra giá")}finally{setBusy(false)}};
  const assess=async()=>{setBusy(true);setError("");try{setSafety(await api.assessSafety({observed_facts:["price_dispute"],confinement_indicators:[],threat_indicators:[],ability_to_leave:true}))}catch(e){setError(e instanceof ProblemDetailError?`${e.message} · ${e.requestId}`:"Không thể tải hướng dẫn")}finally{setBusy(false)}};
  return <>
    <section className={styles.contextCard}><h2>Kiểm tra giá nhanh</h2><p>So sánh với cohort đã phiên bản hóa; kết quả không phải bằng chứng gian lận.</p><form className={styles.toolForm} onSubmit={checkPrice}><label>Số tiền (VND)<input className={styles.field} inputMode="numeric" value={amount} onChange={e=>setAmount(e.target.value)}/></label><button className={styles.button} disabled={busy}>Kiểm tra</button></form>{price&&<PriceInsight insight={price}/>}</section>
    <section className={styles.contextCard} id="safety"><h2><span aria-hidden>⚠ </span>Hướng dẫn an toàn</h2><p>Mở đánh giá rule-first; không tự gọi hay chia sẻ vị trí.</p><button className={styles.button} type="button" onClick={assess} disabled={busy}>Đánh giá tình huống mẫu</button>{safety&&<SafetyAlert assessment={safety}/>}</section>
    {error&&<p className={styles.errorText} role="alert">{error}</p>}<section className={styles.contextCard}><h2>Khám phá địa điểm</h2><p>Chọn vùng thủ công trước khi xem nội dung gần đây.</p><Link className={styles.contextLink} href="/community?mode=nearby&region_id=hanoi-hoan-kiem">Xem quanh Hoàn Kiếm →</Link></section>
  </>;
}
