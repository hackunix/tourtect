import AssistantExperience from "@/components/AssistantExperience";
import AppShell from "@/components/AppShell";
import { getDictionary, getLocale } from "@/lib/i18n";

export const dynamic = "force-dynamic";

export default async function AssistantPage({ searchParams }: { searchParams: Promise<{ mode?: string }> }) {
  const locale=await getLocale();
  const dict=getDictionary(locale);
  const {mode}=await searchParams;
  return <AppShell dict={dict}><AssistantExperience locale={locale} entryMode={mode}/></AppShell>;
}
