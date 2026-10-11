import { useTranslation } from 'react-i18next'

export function LoginPrompt() {
  const { t } = useTranslation()
  return (
    <div className='flex justify-center px-8 pt-12 text-center'>
      <p className='max-w-sm text-xl text-white/70'>{t('authNeedReloginHint')}</p>
    </div>
  )
}
