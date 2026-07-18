import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Tourtect — Travel Price Transparency & Community Safety Shield",
  description: "Tourtect helps travelers scan pricing anomalies, report taxi scams, and obtain rule-first instant safety playbooks during local disputes.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  );
}
