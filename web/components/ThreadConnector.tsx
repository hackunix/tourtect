"use client";

import { FormEvent, useState } from "react";
import { api, Comment, ProblemDetailError } from "@/lib/api";
import styles from "./community.module.css";

export default function ThreadConnector({ postId, initialCount = 0 }: { postId: string; initialCount?: number }) {
  const [open,setOpen]=useState(false); const [comments,setComments]=useState<Comment[]>([]); const [loading,setLoading]=useState(false); const [body,setBody]=useState(""); const [error,setError]=useState("");
  const load=async()=>{setOpen(true);if(comments.length||loading)return;setLoading(true);setError("");try{setComments((await api.getComments(postId)).items)}catch(e){setError(e instanceof ProblemDetailError?`${e.message} · ${e.requestId}`:"Không thể tải phản hồi") }finally{setLoading(false)}};
  const submit=async(e:FormEvent)=>{e.preventDefault();if(!body.trim())return;setLoading(true);setError("");try{const item=await api.addComment(postId,body.trim());setComments(old=>[...old,item]);setBody("")}catch(e){setError(e instanceof ProblemDetailError?e.message:"Không thể gửi phản hồi")}finally{setLoading(false)}};
  const roots=comments.filter(item=>!item.parent_comment_id); const children=(id:string)=>comments.filter(item=>item.parent_comment_id===id);
  return <>
    <button className={styles.action} type="button" onClick={load} aria-expanded={open}>↩ Phản hồi {initialCount > 0 ? initialCount : ""}</button>
    {open&&<div className={styles.thread}>{loading&&comments.length===0&&<p className={styles.notice} role="status">Đang tải phản hồi…</p>}{roots.map(root=><div className={styles.comment} key={root.comment_id}><div className={styles.commentMeta}>{root.author?.display_name||"Traveler"}</div><p>{root.body}</p>{children(root.comment_id).slice(0,3).map(reply=><div className={styles.thread} key={reply.comment_id}><div className={styles.comment}><div className={styles.commentMeta}>{reply.author?.display_name||"Traveler"}</div><p>{reply.body}</p></div></div>)}</div>)}
      <form className={styles.replyForm} onSubmit={submit}><label className="sr-only" htmlFor={`reply-${postId}`}>Viết phản hồi</label><input className={styles.field} id={`reply-${postId}`} value={body} onChange={e=>setBody(e.target.value)} placeholder="Thêm phản hồi…"/><button className={styles.button} disabled={loading||!body.trim()}>Gửi</button></form>{error&&<p className={styles.errorText} role="alert">{error}</p>}</div>}
  </>;
}
