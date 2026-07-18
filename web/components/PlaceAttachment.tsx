import Link from "next/link";
import type { PlaceAttachmentData } from "@/lib/api";
import styles from "./community.module.css";

export default function PlaceAttachment({ place }: { place: PlaceAttachmentData }) { return <Link className={styles.placeAttachment} href={`/places/${place.place_id}`}><strong>{place.name}</strong><span>{place.category} · {place.region_id} · cập nhật {new Date(place.freshness).toLocaleDateString("vi-VN")}</span></Link>; }
