import React from 'react'
import { Icon } from '@iconify/react'
import { useAnimate } from 'framer-motion'
import { useTranslation } from 'react-i18next'

import { type model } from '@model'
import { cn } from '@/helpers/cn'

import type { LocalizationKey } from './i18n'

type ErrorContextType = [
  error: model.FormattedError | null,
  setError: React.Dispatch<model.FormattedError>
]
const ErrorContext = React.createContext<ErrorContextType | null>(null)
export const useErrorPopup = () => React.useContext(ErrorContext)![1]

export function ErrorPopupProvider(props: React.PropsWithChildren) {
  const { t } = useTranslation()
  const [scope, animate] = useAnimate()
  const [error, setError] = React.useState<model.FormattedError | null>(null)

  React.useEffect(() => {
    if (error === null) {
      return
    }
    animate('#error-message', { opacity: [0, 1] }).then(() => {
      animate('#error-message', { opacity: [1, 0] }, { delay: 3.5 }).then(() => setError(null))
    })
  }, [error])

  return (
    <ErrorContext.Provider value={[error, setError]}>
      <div ref={scope} className='absolute flex w-full justify-end'>
        <div
          id='error-message'
          className={cn(
            'flex items-center justify-around gap-6',
            'fixed z-50 mx-auto',
            'pointer-events-none rounded-bl-xl px-8 py-3 text-xl backdrop-blur-sm',
            'bg-gradient-to-r from-[#870e65] to-[#6c086d] text-white'
          )}
          style={{ opacity: 0 }}
        >
          <Icon
            icon='material-symbols:warning-outline'
            className='h-8 w-8 animate-blink text-[#ff6388]'
          />
          {error?.localizationKey && <span>{t(error.localizationKey)}</span>}
        </div>
      </div>
      {props.children}
    </ErrorContext.Provider>
  )
}
