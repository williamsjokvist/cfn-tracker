import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

import eng from './locales/en.json';
import fr from './locales/fr.json';
import jp from './locales/jp.json';

type AppLanguage = {
  code: string;
  nativeName: string;
}

export const APP_LANGUAGES: AppLanguage[] = [
  { code: 'en', nativeName: 'English' },
  { code: 'fr', nativeName: 'Français' },
  { code: 'jp', nativeName: '日本' },
];

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