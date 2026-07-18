import type { SafetyAssessment } from "@/lib/api";
import styles from "./community.module.css";

const labels={critical:"Nguy cấp",urgent:"Cần chú ý ngay",non_emergency:"Không khẩn cấp",information:"Thông tin an toàn"};
export default function SafetyAlert({ assessment }: { assessment: SafetyAssessment }) { return <section className={styles.safetyAlert} data-urgency={assessment.urgency} role={assessment.urgency==="critical"?"alert":"status"}><h3><span aria-hidden>⚠ </span>{labels[assessment.urgency]}</h3>{assessment.silent_mode_recommended&&<p><strong>Khuyến nghị chế độ im lặng.</strong></p>}<ul>{assessment.safe_actions.map((action,index)=><li key={`${action}-${index}`}>{action}</li>)}</ul><p>Chỉ dẫn đã duyệt · phiên bản {assessment.safety_directory_version}</p>{assessment.surface_emergency_options&&<p><strong>Thiết bị có thể hiển thị lựa chọn gọi khẩn cấp; Tourtect sẽ không tự gọi hoặc tự chia sẻ vị trí.</strong></p>}</section>; }
