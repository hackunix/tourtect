import type { Metadata } from "next";
import { Noto_Sans } from "next/font/google";
import { getLocale } from "@/lib/i18n";
import "./globals.css";

const noto = Noto_Sans({ subsets: ["latin", "vietnamese", "cyrillic"], display: "swap", variable: "--font-noto" });
export const metadata: Metadata = { title: "Tourtect — Cộng đồng du lịch an toàn", description: "Kinh nghiệm địa phương, giá cả minh bạch và hướng dẫn an toàn cho du khách." };

export default async function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  const locale=await getLocale();
  return <html lang={locale} suppressHydrationWarning className={noto.variable}><head><script dangerouslySetInnerHTML={{__html:`try{var t=localStorage.getItem('tourtect_theme');if(t)document.documentElement.dataset.theme=t}catch(e){}`}}/></head><body>{children}</body></html>;
}
