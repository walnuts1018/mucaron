import { getRequestConfig } from "next-intl/server";

export default getRequestConfig(async () => {
  // Provide a static locale, fetch a user setting,
  // read from `cookies()`, `headers()`, etc.
  const locale = "ja";

  return {
    locale,
    messages: (await import(`./locales/${locale}.json`)).default,
  };
});
