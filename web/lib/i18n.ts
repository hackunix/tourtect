import "server-only";
import { cookies, headers } from "next/headers";
import type { Locale } from "./api";

const dictionaries = {
  "vi-VN": {
    home: "Trang chủ", search: "Tìm kiếm", create: "Đăng bài", saved: "Đã lưu", notifications: "Thông báo",
    profile: "Hồ sơ", sos: "Trợ giúp khẩn cấp", feedTitle: "Cộng đồng du lịch",
    feedSubtitle: "Kinh nghiệm địa phương, giá cả minh bạch và thông tin an toàn.",
    following: "Đang theo dõi", nearby: "Gần đây", latest: "Mới nhất", trending: "Nổi bật", safety: "An toàn",
    retry: "Thử lại", empty: "Chưa có bài viết phù hợp.", composerPrompt: "Bạn muốn du khách biết điều gì?",
  },
  en: {
    home: "Home", search: "Search", create: "Create", saved: "Saved", notifications: "Notifications",
    profile: "Profile", sos: "Emergency help", feedTitle: "Traveler community",
    feedSubtitle: "Local knowledge, transparent prices, and safety-aware guidance.",
    following: "Following", nearby: "Nearby", latest: "Latest", trending: "Trending", safety: "Safety",
    retry: "Retry", empty: "No matching posts yet.", composerPrompt: "What should travelers know?",
  },
} as const;

export async function getLocale(): Promise<Locale> {
  const cookieLocale = (await cookies()).get("tourtect_locale")?.value;
  if (cookieLocale === "en" || cookieLocale === "vi-VN") return cookieLocale;
  const preferred = (await headers()).get("accept-language") || "";
  return preferred.toLowerCase().startsWith("en") ? "en" : "vi-VN";
}

export type Dictionary = { [K in keyof typeof dictionaries["vi-VN"]]: string };
export const getDictionary = (locale: Locale): Dictionary => dictionaries[locale];
