import { api, FeedResponse, PlaceDetail, ProblemDetailError } from "@/lib/api";
import { getDictionary,getLocale } from "@/lib/i18n";
import AppShell from "@/components/AppShell";
import PostItem from "@/components/PostItem";
import ProblemDetailsState from "@/components/ProblemDetailsState";
import styles from "@/components/community.module.css";
export const dynamic="force-dynamic";

export default async function PlacePage({params}:{params:Promise<{id:string}>}){
  const {id}=await params;const dict=getDictionary(await getLocale());let place:PlaceDetail|null=null;let feed:FeedResponse|null=null;let problem:ProblemDetailError|null=null;
  try{[place,feed]=await Promise.all([api.getPlace(id),api.getFeed({mode:"latest",limit:30})])}catch(error){if(error instanceof ProblemDetailError)problem=error}
  if(problem||!place||!feed)return <AppShell dict={dict}><main className={styles.page}><ProblemDetailsState detail={problem?.message||"Không thể tải địa điểm"} requestId={problem?.requestId}/></main></AppShell>;
  const posts=feed.items.filter(post=>post.places?.some(item=>item.place_id===id)||post.place_ids?.includes(id));
  return <AppShell dict={dict}><main className={styles.page}><header className={styles.placeHeader}><p>{place.category} · {place.region_id}</p><h1>{place.name}</h1>{place.address&&<p>{place.address}</p>}{place.description&&<p>{place.description}</p>}<div className={styles.badges}>{place.average_rating!==undefined&&<span className={styles.badge}>Đánh giá {place.average_rating.toFixed(1)}</span>}{place.freshness&&<span className={styles.badge}>Cập nhật {new Date(place.freshness).toLocaleDateString("vi-VN")}</span>}</div></header>{posts.length?posts.map(post=><PostItem post={post} key={post.post_id}/>):<div className={styles.empty}>Chưa có bài viết công khai cho địa điểm này.</div>}</main></AppShell>;
}
