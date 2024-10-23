import type { Metadata } from "next";
import { Nunito, Noto_Sans_JP } from "next/font/google";
import "./globals.css";

const NunitoFont = Nunito({
  subsets: ["latin"],
  variable: "--font-Nunito",
});

const NotoFont = Noto_Sans_JP({
  subsets: ["latin"],
  variable: "--font-Noto",
});
const title = "mucaron";
const description = "mucaron is music streaming app";
const url = "https://mucaron.walnuts.dev";

export const metadata: Metadata = {
  metadataBase: new URL(url),
  title: title,
  description: description,
  openGraph: {
    title: title,
    description,
    url,
    siteName: title,
    locale: "ja_JP",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <head>
        <meta name="msapplication-TileColor" content="#feb4c1" />
        <meta name="theme-color" content="#feb4c1" />
      </head>
      <body className={`${NunitoFont.variable} ${NotoFont.variable} font-Noto`}>
        {children}
      </body>
    </html>
  );
}
