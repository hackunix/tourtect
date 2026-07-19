import type { Metadata } from "next";
import { getLocale } from "@/lib/i18n";
import "./globals.css";

export const metadata: Metadata = { title: "Tourtect — AI travel safety companion", description: "A place-aware travel companion grounded in verified community knowledge and deterministic safety and price guidance." };

export default async function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  const locale=await getLocale();
  return <html lang={locale} suppressHydrationWarning><head><script dangerouslySetInnerHTML={{__html:`try{var t=localStorage.getItem('tourtect_theme');if(t)document.documentElement.dataset.theme=t}catch(e){}`}}/></head><body>{children}</body></html>;
}
