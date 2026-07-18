"use client";

import { useState } from "react";
import type { Post } from "@/lib/api";
import { api, ProblemDetailError } from "@/lib/api";
import EvidenceBadge from "./EvidenceBadge";
import PlaceAttachment from "./PlaceAttachment";
import ThreadConnector from "./ThreadConnector";
import styles from "./community.module.css";

const typeLabels: Record<string,string>={discussion:"Thảo luận",question:"Câu hỏi",review:"Đánh giá",price_report:"Báo giá",scam_report:"Mẫu rủi ro được báo cáo",tip:"Mẹo",official_alert:"Cảnh báo chính thức"};
export default function PostItem({ post }: { post: Post }) {
  const [useful,setUseful]=useState(Boolean(post.viewer_useful));const [saved,setSaved]=useState(Boolean(post.viewer_saved));const [usefulCount,setUsefulCount]=useState(post.useful_count||0);const [error,setError]=useState("");
  const toggleUseful=async()=>{const next=!useful;setUseful(next);setUsefulCount(c=>Math.max(0,c+(next?1:-1)));try{await api.setUseful(post.post_id,next)}catch(e){setUseful(!next);setUsefulCount(c=>Math.max(0,c+(next?-1:1)));setError(e instanceof ProblemDetailError?e.message:"Không thể cập nhật")}};
  const toggleSaved=async()=>{const next=!saved;setSaved(next);try{await api.setSaved(post.post_id,next)}catch(e){setSaved(!next);setError(e instanceof ProblemDetailError?e.message:"Không thể lưu bài")}};
  const share=async()=>{const url=`${location.origin}/?post=${post.post_id}`;try{if(navigator.share)await navigator.share({title:post.title,url});else await navigator.clipboard.writeText(url)}catch{/* user cancelled */}};
  const initials=(post.author?.display_name||"T").slice(0,1).toUpperCase();
  return <article className={styles.post} id={`post-${post.post_id}`}>
    <div className={styles.avatar} aria-hidden>{initials}</div><div className={styles.postMain}>
      <div className={styles.postMeta}><span className={styles.author}>{post.author?.display_name||"Tourtect traveler"}</span><span aria-hidden>·</span><time dateTime={post.created_at}>{new Date(post.created_at).toLocaleDateString("vi-VN",{day:"2-digit",month:"short"})}</time></div>
      <h2 className={styles.postTitle}>{post.title}</h2><p className={styles.postBody}>{post.body}</p>
      <div className={styles.badges}><span className={styles.badge}>{typeLabels[post.post_type]||post.post_type}</span><EvidenceBadge level={post.evidence_level}/>{post.original_locale!=="vi-VN"&&<span className={styles.badge}>Bản gốc: {post.original_locale}</span>}</div>
      {post.places?.map(place=><PlaceAttachment place={place} key={place.place_id}/>)}
      <div className={styles.postActions}><ThreadConnector postId={post.post_id} initialCount={post.comment_count}/><button className={styles.action} data-active={useful} type="button" onClick={toggleUseful} aria-pressed={useful}>✓ Hữu ích {usefulCount||""}</button><button className={styles.action} data-active={saved} type="button" onClick={toggleSaved} aria-pressed={saved}>◇ {saved?"Đã lưu":"Lưu"}</button><button className={styles.action} type="button" onClick={share}>↗ Chia sẻ</button></div>
      {error&&<p className={styles.errorText} role="alert">{error}</p>}
    </div>
  </article>;
}
