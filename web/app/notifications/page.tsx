import { api, Notification, ProblemDetailError } from "@/lib/api";
import { getDictionary,getLocale } from "@/lib/i18n";
import AppShell from "@/components/AppShell";
import ProblemDetailsState from "@/components/ProblemDetailsState";
import styles from "@/components/community.module.css";
export const dynamic="force-dynamic";

export default async function NotificationsPage(){
  const dict=getDictionary(await getLocale());let items:Notification[]|null=null;let problem:ProblemDetailError|null=null;
  try{items=(await api.getNotifications()).items}catch(error){if(error instanceof ProblemDetailError)problem=error}
  return <AppShell dict={dict}><main className={styles.page}><header className={styles.pageHeader}><h1>{dict.notifications}</h1></header>
    {problem?<ProblemDetailsState detail={problem.message} requestId={problem.requestId}/>:items?.length?items.map(item=><article className={styles.post} key={item.notification_id}><div className={styles.avatar} aria-hidden>○</div><div><p>{item.message}</p><time className={styles.notice} dateTime={item.created_at}>{new Date(item.created_at).toLocaleString("vi-VN")}</time></div></article>):<div className={styles.empty}>Bạn chưa có thông báo mới.</div>}
  </main></AppShell>;
}
