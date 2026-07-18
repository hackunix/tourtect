import { Suspense } from "react";
import { api, FeedMode, PlaceSummary, ProblemDetailError } from "@/lib/api";
import { getDictionary, getLocale } from "@/lib/i18n";
import AppShell from "@/components/AppShell";
import ContextPanel from "@/components/ContextPanel";
import FeedHeader from "@/components/FeedHeader";
import FeedSkeleton from "@/components/FeedSkeleton";
import PostComposer from "@/components/PostComposer";
import PostItem from "@/components/PostItem";
import ProblemDetailsState from "@/components/ProblemDetailsState";
import styles from "@/components/community.module.css";

export const dynamic = "force-dynamic";
const validModes: FeedMode[]=["following","nearby","latest","trending","safety"];

async function FeedItems({mode,regionId,empty}:{mode:FeedMode;regionId?:string;empty:string}){
  let feed; let problem: ProblemDetailError | null = null;
  try { feed=await api.getFeed({mode,region_id:regionId,limit:20}); } catch(error) { if(error instanceof ProblemDetailError) problem=error; }
  if(problem)return <ProblemDetailsState detail={problem.message} requestId={problem.requestId}/>;
  if(!feed)return <ProblemDetailsState detail="Không thể kết nối tới Tourtect API. Không có dữ liệu giả được hiển thị."/>;
  if(feed.items.length===0)return <div className={styles.empty}>{empty}</div>;
  return <>{feed.items.map(post=><PostItem post={post} key={post.post_id}/>)}</>;
}

export default async function Home({searchParams}:{searchParams:Promise<{mode?:string;region_id?:string}>}){
  const params=await searchParams;const mode=validModes.includes(params.mode as FeedMode)?params.mode as FeedMode:"latest";const regionId=params.region_id;const locale=await getLocale();const dict=getDictionary(locale);
  let places: PlaceSummary[]=[];try{places=(await api.getPlaces({limit:20})).items}catch{/* composer remains usable without a place attachment */}
  return <AppShell dict={dict} context={<ContextPanel/>}><main className={styles.feed}><header className={styles.feedIntro}><h1>{dict.feedTitle}</h1><p>{dict.feedSubtitle}</p></header><FeedHeader mode={mode} regionId={regionId} dict={dict}/><PostComposer locale={locale} places={places} prompt={dict.composerPrompt}/><Suspense fallback={<FeedSkeleton/>}><FeedItems mode={mode} regionId={regionId} empty={dict.empty}/></Suspense></main></AppShell>;
}
