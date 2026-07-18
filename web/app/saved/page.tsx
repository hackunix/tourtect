import { api, FeedResponse, ProblemDetailError } from "@/lib/api";
import { getDictionary,getLocale } from "@/lib/i18n";
import AppShell from "@/components/AppShell";
import PostItem from "@/components/PostItem";
import ProblemDetailsState from "@/components/ProblemDetailsState";
import styles from "@/components/community.module.css";
export const dynamic="force-dynamic";

export default async function SavedPage(){
  const dict=getDictionary(await getLocale());let data:FeedResponse|null=null;let problem:ProblemDetailError|null=null;
  try{data=await api.getSaved()}catch(error){if(error instanceof ProblemDetailError)problem=error}
  return <AppShell dict={dict}><main className={styles.page}><header className={styles.pageHeader}><h1>{dict.saved}</h1></header>
    {problem?<ProblemDetailsState detail={problem.message} requestId={problem.requestId}/>:data?.items.length?data.items.map(post=><PostItem post={{...post,viewer_saved:true}} key={post.post_id}/>):<div className={styles.empty}>Các bài bạn lưu sẽ xuất hiện ở đây.</div>}
  </main></AppShell>;
}
