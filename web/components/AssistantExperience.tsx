"use client";

import Link from "next/link";
import { FormEvent, useCallback, useEffect, useRef, useState } from "react";
import {
  api, AssistantConfirmation, AssistantEvidence, AssistantResponse, AssistantToolResult,
  Locale, PriceInsight as PriceInsightType, ProblemDetailError, SafetyAssessment, safeUUID,
} from "@/lib/api";
import PriceInsight from "./PriceInsight";
import SafetyAlert from "./SafetyAlert";
import styles from "./assistant.module.css";

const SESSION_KEY = "tourtect_assistant_session";
const DRAFT_KEY = "tourtect_assistant_draft";

type ConversationItem =
  | { kind: "user"; id: string; text: string }
  | { kind: "assistant"; id: string; response: AssistantResponse };

const quickActions = [
  { label: "Check a taxi price", prompt: "taxi price from Noi Bai airport to Hanoi Old Quarter is 350,000 VND" },
  { label: "Translate to Vietnamese", prompt: "translate 'Can you show me the way to the market?' to Vietnamese" },
  { label: "Report a taxi issue", prompt: "the driver is forcing me to pay extra and won't let me leave" },
  { label: "Ask about Hanoi Station", prompt: "tell me about Hanoi railway station" },
];

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function priceInsight(output: Record<string, unknown>): PriceInsightType | null {
  const candidate = isRecord(output.insight) ? output.insight : output;
  return typeof candidate.alert_level === "string" && isRecord(candidate.observed) && Array.isArray(candidate.reasons)
    ? candidate as unknown as PriceInsightType : null;
}

function safetyAssessment(output: Record<string, unknown>): SafetyAssessment | null {
  const candidate = isRecord(output.assessment) ? output.assessment : output;
  return typeof candidate.urgency === "string" && Array.isArray(candidate.safe_actions) && typeof candidate.safety_directory_version === "string"
    ? candidate as unknown as SafetyAssessment : null;
}

function displayValue(value: unknown): string {
  if (value === null || value === undefined) return "—";
  if (typeof value === "string" || typeof value === "number" || typeof value === "boolean") return String(value);
  if (Array.isArray(value)) return value.map(displayValue).join(" · ");
  return JSON.stringify(value);
}

export function EvidenceList({ evidence }: { evidence: AssistantEvidence[] }) {
  if (!evidence.length) return null;
  return <section className={styles.evidenceList} aria-label="Evidence used"><h3>Why Tourtect is showing this</h3><ul>{evidence.map(item => <li key={item.evidence_id}><strong>{item.title}</strong><div>{item.summary}</div><div className={styles.evidenceMeta}><span>{item.evidence_level}</span><span>{item.freshness}</span><span>{item.source_type.replaceAll("_", " ")}</span>{item.observed_at&&<time dateTime={item.observed_at}>{new Date(item.observed_at).toLocaleDateString()}</time>}{item.source_url&&/^https?:\/\//.test(item.source_url)&&<a className={styles.evidenceLink} href={item.source_url} target="_blank" rel="noreferrer">Source</a>}</div></li>)}</ul></section>;
}

function ToolResultCard({ result }: { result: AssistantToolResult }) {
  const price = priceInsight(result.output);
  const safety = safetyAssessment(result.output);
  if (price) return <PriceInsight insight={price}/>;
  if (safety) return <SafetyAlert assessment={safety}/>;
  const visible = Object.entries(result.output).slice(0, 12);
  return <section className={styles.toolResult} data-status={result.status}><h3>{result.tool_name.replaceAll("_", " ")} · {result.status.replaceAll("_", " ")}</h3>{visible.length>0&&<dl className={styles.toolOutput}>{visible.map(([key,value])=><div key={key}><dt>{key.replaceAll("_", " ")}</dt><dd>{displayValue(value)}</dd></div>)}</dl>}{result.error_category&&<p>Could not complete: {result.error_category}</p>}</section>;
}

export function ConfirmationCard({ confirmation, busy, outcome, confirmedTarget, onDecision }: { confirmation: AssistantConfirmation; busy: boolean; outcome?: string; confirmedTarget?: string; onDecision: (decision: "confirmed" | "rejected") => void }) {
  const safeDialerTarget=confirmedTarget&&/^tel:[+0-9#*()-]+$/.test(confirmedTarget)?confirmedTarget:undefined;
  return <section className={styles.confirmation}><h3>{confirmation.title}</h3><p>{confirmation.description}</p><small>Expires {new Date(confirmation.expires_at).toLocaleString()}</small>{outcome?<><p role="status">{outcome}</p>{safeDialerTarget&&<a className={styles.primaryButton} href={safeDialerTarget}>Open verified dialer</a>}</>:<div className={styles.messageActions}><button className={styles.primaryButton} disabled={busy} onClick={()=>onDecision("confirmed")}>Confirm</button><button className={styles.secondaryButton} disabled={busy} onClick={()=>onDecision("rejected")}>Not now</button></div>}</section>;
}

export function ProviderDegradedState({ detail }: { detail?: string }) {
  return <section className={styles.degraded} role="alert"><h2>Automatic understanding is unavailable</h2><p>{detail||"Tourtect could not interpret this automatically right now."}</p><ul><li>Use the manual Price Check on Community</li><li>Use the rule-first Safety Assessment</li><li>Keep a private draft until service returns</li></ul><Link className={styles.secondaryButton} href="/community">Open reliable fallback tools</Link></section>;
}

function SuggestedActions({ response, onPrompt }: { response: AssistantResponse; onPrompt: (prompt: string) => void }) {
  if (!response.suggested_actions.length) return null;
  return <div className={styles.suggestedActions} aria-label="Suggested actions">{response.suggested_actions.map(action => {
    const safeTarget = action.target?.startsWith("/") ? action.target : undefined;
    if (action.action_type === "deep_link" && safeTarget && !action.requires_confirmation) return <Link className={styles.secondaryButton} href={safeTarget} key={action.action_id}>{action.label}</Link>;
    return <button className={styles.secondaryButton} key={action.action_id} onClick={()=>onPrompt(action.target||action.label)}>{action.label}</button>;
  })}</div>;
}

function AssistantResponseCard({ response, sessionId, onPrompt }: { response: AssistantResponse; sessionId: string | null; onPrompt: (prompt: string) => void }) {
  const [confirmationBusy,setConfirmationBusy]=useState(false);
  const [confirmationOutcome,setConfirmationOutcome]=useState<string>();
  const [confirmedTarget,setConfirmedTarget]=useState<string>();
  const [feedback,setFeedback]=useState<string>();
  const confirm=async(decision:"confirmed"|"rejected")=>{
    if(!sessionId||!response.requested_confirmation)return;
    setConfirmationBusy(true);
    try{const result=await api.confirmAssistantAction(sessionId,{confirmation_id:response.requested_confirmation.confirmation_id,decision});setConfirmedTarget(result.status==="confirmed"?result.target:undefined);setConfirmationOutcome(result.status==="confirmed"?"Action confirmed by the backend. Tourtect has not placed a call.":"Action declined.")}catch(error){setConfirmationOutcome(error instanceof ProblemDetailError?error.message:"Could not record this decision.")}finally{setConfirmationBusy(false)}
  };
  const sendFeedback=async(kind:"helpful"|"not_helpful")=>{
    if(!sessionId||feedback)return;
    try{await api.createAssistantFeedback(sessionId,{assistant_message_id:response.assistant_message_id,feedback_type:kind});setFeedback("Feedback quarantined for review.")}catch(error){setFeedback(error instanceof ProblemDetailError?error.message:"Could not save feedback.")}
  };
  return <article className={`${styles.assistantMessage} ${response.safety_state==="critical"?styles.critical:response.safety_state==="urgent"?styles.urgent:""}`}><div className={styles.assistantHeader}><strong>Tourtect</strong><span>{response.intent.replaceAll("_", " ")} · {Math.round(response.confidence*100)}%</span></div><p>{response.message}</p>{response.fallback_used&&<div className={styles.fallback}>A deterministic fallback was used. Automatic language generation did not change verified tool results.</div>}{response.tool_results.map(result=><ToolResultCard result={result} key={result.tool_result_id}/>)}{response.factors_considered.length>0&&<section className={styles.detailCard}><h3>Factors considered</h3><ul>{response.factors_considered.map((factor,index)=><li key={`${factor}-${index}`}>{factor}</li>)}</ul></section>}{response.missing_information.length>0&&<section className={styles.detailCard}><h3>Still needed</h3><ul>{response.missing_information.map((item,index)=><li key={`${item}-${index}`}>{item}</li>)}</ul></section>}<EvidenceList evidence={response.evidence}/>{response.requested_confirmation&&<ConfirmationCard confirmation={response.requested_confirmation} busy={confirmationBusy} outcome={confirmationOutcome} confirmedTarget={confirmedTarget} onDecision={confirm}/>}<SuggestedActions response={response} onPrompt={onPrompt}/><div className={styles.messageActions}><button className={styles.secondaryButton} onClick={()=>sendFeedback("helpful")} disabled={Boolean(feedback)}>Helpful</button><button className={styles.secondaryButton} onClick={()=>sendFeedback("not_helpful")} disabled={Boolean(feedback)}>Not helpful</button>{feedback&&<span className={styles.feedbackStatus} role="status">{feedback}</span>}</div><small className={styles.feedbackStatus}>Trace {response.trace_id}{response.dataset_version?` · Data ${response.dataset_version}`:""}{response.freshness?` · ${response.freshness}`:""}</small></article>;
}

export default function AssistantExperience({ locale, entryMode }: { locale: Locale; entryMode?: string }) {
  const [sessionId,setSessionId]=useState<string|null>(null);
  const [items,setItems]=useState<ConversationItem[]>([]);
  const [input,setInput]=useState("");
  const [processingConsent,setProcessingConsent]=useState(true);
  const [isStarting,setIsStarting]=useState(true);
  const [isSending,setIsSending]=useState(false);
  const [error,setError]=useState("");
  const [degraded,setDegraded]=useState(false);
  
  const activeController=useRef<AbortController|null>(null);
  const textareaRef=useRef<HTMLTextAreaElement>(null);
  const bottomRef=useRef<HTMLDivElement>(null);

  useEffect(()=>{
    const controller=new AbortController();activeController.current=controller;
    const timer=window.setTimeout(()=>{
      setInput(localStorage.getItem(DRAFT_KEY)||"");
      const saved=localStorage.getItem(SESSION_KEY);
      if(!saved){setIsStarting(false);return;}
      api.getAssistantSession(saved,controller.signal).then(session=>{setSessionId(session.session_id);setProcessingConsent(session.context.consent_state.processing);setItems((session.recent_responses||[]).map(response=>({kind:"assistant" as const,id:response.assistant_message_id,response}))) }).catch(error=>{if(error instanceof DOMException&&error.name==="AbortError")return;localStorage.removeItem(SESSION_KEY);if(!(error instanceof ProblemDetailError&&error.status===404)){setError(error instanceof ProblemDetailError?error.message:"Could not resume the assistance session.");setDegraded(true)}}).finally(()=>setIsStarting(false));
    },0);
    return()=>{window.clearTimeout(timer);controller.abort()};
  },[]);

  useEffect(()=>{localStorage.setItem(DRAFT_KEY,input)},[input]);

  useEffect(()=>{
    if(textareaRef.current){
      textareaRef.current.style.height="auto";
      textareaRef.current.style.height=`${textareaRef.current.scrollHeight}px`;
    }
  },[input]);

  useEffect(()=>{
    if(typeof bottomRef.current?.scrollIntoView === "function"){
      bottomRef.current.scrollIntoView({behavior:"smooth"});
    }
  },[items,isSending,degraded]);

  const ensureSession=useCallback(async(signal:AbortSignal)=>{
    if(sessionId)return sessionId;
    const session=await api.createAssistantSession({locale,interaction_mode:"text",processing_consent:processingConsent},signal);
    setSessionId(session.session_id);localStorage.setItem(SESSION_KEY,session.session_id);return session.session_id;
  },[locale,processingConsent,sessionId]);

  const send=async(event?:FormEvent)=>{
    event?.preventDefault();const text=input.trim();if(!text||isSending)return;
    if(!processingConsent){setError("Confirm processing consent before sending this message.");return;}
    setError("");setDegraded(false);setIsSending(true);
    const messageId=safeUUID();setItems(old=>[...old,{kind:"user",id:messageId,text}]);setInput("");
    const controller=new AbortController();activeController.current=controller;
    try{
      const id=await ensureSession(controller.signal);
      const response=await api.createAssistantMessage(id,{message_id:messageId,input_type:"text",text,locale,user_confirmed:false},controller.signal);
      setItems(old=>[...old,{kind:"assistant",id:response.assistant_message_id,response}])
    }catch(error){
      if(error instanceof DOMException&&error.name==="AbortError")return;
      setError(error instanceof ProblemDetailError?`${error.message}${error.requestId?` · ${error.requestId}`:""}`:"Tourtect could not interpret this automatically right now.");
      setDegraded(true)
    }finally{
      setIsSending(false);
      activeController.current=null
    }
  };



  const cancel=()=>{activeController.current?.abort();activeController.current=null;setIsSending(false)};
  const reset=async()=>{
    cancel();const old=sessionId;setSessionId(null);setItems([]);setError("");setDegraded(false);setProcessingConsent(true);localStorage.removeItem(SESSION_KEY);localStorage.removeItem(DRAFT_KEY);setInput("");
    if(old){try{await api.deleteAssistantSession(old)}catch{/* Local reset remains effective if the expired session is already gone. */}}
  };
  const choosePrompt=(prompt:string)=>{setInput(prompt);document.getElementById("assistant-input")?.focus()};
  const modeNotice=entryMode==="voice"?"Live voice is not connected on web yet. Type the utterance below; no simulated transcript will be created.":entryMode==="lens"?"Image capture is not available until the consent-backed capture API is implemented. Describe what you see below.":entryMode==="sos"?"Describe what is happening. Critical rules run before normal assistant routing, and Tourtect will never call or share location automatically.":null;

  return <main className={styles.assistantPage}><header className={styles.hero}><p className={styles.eyebrow}>AI-native travel safety companion</p><h1>What is happening?</h1><p className={styles.heroText}>You can type what is happening. Tourtect coordinates verified knowledge with deterministic Price and Safety engines, then shows the evidence.</p>{modeNotice&&<p className={styles.modeNotice} role="status">{modeNotice}</p>}</header><section className={styles.stream} aria-live="polite">{items.length===0&&<section className={styles.welcome}><h2>Ask naturally</h2><p>You do not need to choose a technical tool first.</p><div className={styles.quickActions}>{quickActions.map(action=><button className={styles.quickButton} onClick={()=>choosePrompt(action.prompt)} key={action.label}>{action.label}</button>)}</div></section>}{items.map(item=>item.kind==="user"?<div className={styles.userMessage} key={item.id}>{item.text}</div>:<AssistantResponseCard response={item.response} sessionId={sessionId} onPrompt={choosePrompt} key={item.id}/>)}{degraded&&<ProviderDegradedState detail={error}/>}<div ref={bottomRef} /></section><aside className={styles.composerDock} aria-label="Assistant composer"><form className={styles.composer} onSubmit={send}><label className="sr-only" htmlFor="assistant-input">Describe what is happening</label><textarea ref={textareaRef} rows={1} id="assistant-input" value={input} onChange={event=>setInput(event.target.value)} placeholder="Describe what is happening…" disabled={isStarting} style={{ overflowY: "hidden" }}/><div className={styles.composerControls}>{isSending?<button className={styles.dangerButton} type="button" onClick={cancel}>Cancel</button>:<button className={styles.primaryButton} disabled={!input.trim()||isStarting} type="submit">Send</button>}<button className={styles.secondaryButton} type="button" onClick={reset} disabled={isStarting}>Reset</button></div></form><label className={styles.consent}><input type="checkbox" checked={processingConsent} onChange={event=>setProcessingConsent(event.target.checked)} disabled={Boolean(sessionId)}/><span>Allow Tourtect to process messages in this temporary assistance session. Raw audio and camera frames are not uploaded by this web flow.</span></label><p className={`${styles.statusLine} ${error?styles.error:""}`} role={error?"alert":"status"}>{isStarting?"Resuming session…":isSending?"Tourtect is coordinating trusted tools…":error||sessionId?`Session active${sessionId?` · ${sessionId.slice(0,8)}`:""}`:"Ready when you are."}</p></aside></main>;
}
