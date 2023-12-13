import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

type AppLanguage = {
  code: string;
  nativeName: string;
  translation: Object;
}

export const APP_LANGUAGES: AppLanguage[] = [
  { code: 'en-GB', nativeName: 'English', translation: await import('./locales/en.json') },
  { code: 'fr', nativeName: 'Français', translation: await import('./locales/fr.json') },
  { code: 'ja-JP', nativeName: '日本', translation: await import('./locales/jp.json') },
];

// https://www.i18next.com/overview/configuration-options
export const initLocalization = () => {
  i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: APP_LANGUAGES[0].code,
    lng: window.localStorage.getItem('lng'),
    resources:
      APP_LANGUAGES.reduce((obj, item) => {
        return {
          ...obj,
          [item['code']]: { translation: item.translation }
        }
      }, {})
  });
}
