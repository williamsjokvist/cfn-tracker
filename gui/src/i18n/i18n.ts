import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

type AppLanguage = {
  code: string;
  nativeName: string;
  translation: Object;
}

export const APP_LANGUAGES: AppLanguage[] = [
  { code: 'en', nativeName: 'English', translation: await import('@/locales/en.json') },
  { code: 'fr', nativeName: 'Français', translation: await import('@/locales/fr.json') },
  { code: 'jp', nativeName: '日本', translation: await import('@/locales/jp.json') },
];

// https://www.i18next.com/overview/configuration-options
i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    debug: false,
    fallbackLng: APP_LANGUAGES[0].code,
    interpolation: {
      escapeValue: false,
    },
    lng: window.localStorage.getItem('lng'),
    resources: 
      APP_LANGUAGES.reduce((obj, item) => {
        return {
          ...obj,
          [item['code']]: { translation: item.translation }
        }
      }, {})
  });

export default i18n;