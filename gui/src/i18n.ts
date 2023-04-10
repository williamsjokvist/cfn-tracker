import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import eng from './locales/en.json';
import fr from './locales/fr.json';
import jp from './locales/jp.json';

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  // init i18next
  // for all options read: https://www.i18next.com/overview/configuration-options
  .init({
    debug: false,
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false,
    },
    resources: {
      en: { translation: eng },
      fr: { translation: fr },
      jp: { translation: jp }
    }
  });

export default i18n;