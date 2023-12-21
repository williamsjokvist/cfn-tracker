import i18n from "i18next";
import LanguageDetector from 'i18next-browser-languagedetector';
import HttpBackend, { type HttpBackendOptions } from "i18next-http-backend";
import { initReactI18next } from 'react-i18next';

import type { locales } from "@@/go/models";
import { GetTranslation } from "@@/go/core/CommandHandler";

export type LocalizationKey = keyof locales.Localization

// https://www.i18next.com/overview/configuration-options
export const initI18n = () => {
  return i18n
    .use(LanguageDetector)
    .use(HttpBackend)
    .use(initReactI18next)
    .init<HttpBackendOptions>({
      fallbackLng: "en-GB",
      load: "currentOnly",
      lng: "en-GB",
      react: {
        useSuspense: true,
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
