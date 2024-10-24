import type { Metadata } from "next";
import { Nunito, Noto_Sans_JP } from "next/font/google";
import "./globals.css";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";

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

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const locale = await getLocale();
  const messages = await getMessages();

  return (
    <html lang={locale}>
      <head>
        <meta name="msapplication-TileColor" content="#feb4c1" />
        <meta name="theme-color" content="#feb4c1" />
      </head>
      <body className={`${NunitoFont.variable} ${NotoFont.variable} font-Noto`}>
        <NextIntlClientProvider messages={messages}>
          {children}
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
