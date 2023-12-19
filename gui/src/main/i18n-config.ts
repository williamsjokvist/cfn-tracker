import i18n from "i18next";
import LanguageDetector from 'i18next-browser-languagedetector';
import HttpBackend, { type HttpBackendOptions } from "i18next-http-backend";
import { initReactI18next } from 'react-i18next';

import type { locales } from "@@/go/models";
import { GetTranslation } from "@@/go/core/CommandHandler";

type AppLanguage = {
  code: string;
  nativeName: string;
}

export type LocalizationKey = keyof locales.Localization

export const APP_LANGUAGES: AppLanguage[] = [
  { code: 'en-GB', nativeName: 'English' },
  { code: 'fr-FR', nativeName: 'Français' },
  { code: 'ja-JP', nativeName: '日本' },
];

// https://www.i18next.com/overview/configuration-options
export const initLocalization = () => {
  return i18n
    .use(LanguageDetector)
    .use(HttpBackend)
    .use(initReactI18next)
    .init<HttpBackendOptions>({
      fallbackLng: APP_LANGUAGES[0].code,
      load: "currentOnly",
      lng: window.localStorage.getItem("lng") ?? APP_LANGUAGES[0].code,
      react: {
        useSuspense: false,
      },
      backend: {
        loadPath: '{{lng}}',
        request: (options, url, payload, callback) => {
          GetTranslation(url).then((data) => {
            callback(null, {
              status: 200,
              data
            })
          })
        }
      }
    });
}
