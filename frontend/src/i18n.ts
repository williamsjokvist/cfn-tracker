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
          opponent: 'Opponent',
          character: 'Character',
          lpGain: 'LP Gain',
          delete: 'Clear Log',
          goBack: 'Go back',
          loading: 'Loading',
          wins: 'Wins',
          losses: 'Losses',
          winRate: 'Win Rate',
          stop: 'Stop',
          openResultFolder: 'Results Folder',
          enterCfnName: 'Enter your CFN name'
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
          opponent: 'Ennemi',
          character: 'Caractère',
          lpGain: 'Gain LP',
          delete: 'Supprimer le journal',
          goBack: 'Retourner',
          loading: 'Chargement, veuillez patienter',
          wins: 'Gagne',
          losses: 'Pertes',
          winRate: 'Taux de réussite',
          stop: 'Arrêter',
          openResultFolder: 'Dossier des résultats',
          enterCfnName: 'Entrez votre nom CFN'
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
          opponent: '敵',
          character: 'キャラクター',
          lpGain: 'LP得',
          delete: '日誌を削除',
          goBack: '戻る',
          loading: 'お待ちください',
          wins: '勝つ',
          losses: '損失',
          winRate: '勝率',
          stop: 'やめろ',
          openResultFolder: '結果フォルダ',
          enterCfnName: 'CFN名を入力してください'
        }
      }
    }
  });

export default i18n;