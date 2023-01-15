import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  // init i18next
  // for all options read: https://www.i18next.com/overview/configuration-options
  .init({
    debug: true,
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false,
    },
    resources: {
      en: {
        translation: {
          tracking: 'Tracking',
          history: 'Match Log',
          language: 'Language',
          startTracking: 'Start Tracking',
          cfnName: 'CFN Name',
          start: 'Start',
        }
      },
      fr: {
        translation: {
          tracking: 'Suivie',
          history: 'Histoire',
          language: 'Langue',
          startTracking: 'Démarrer le suivi',
          cfnName: 'Nom CFN',
          start: 'Début',
        }
      },
      jp: {
        translation: {
          tracking: '追跡',
          history: 'マッチログ',
          language: '言語',
          startTracking: '追跡を開始',
          cfnName: 'CFN名',
          start: '始める',
        }
      }
    }
  });

export default i18n;