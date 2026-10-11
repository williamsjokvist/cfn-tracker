import React from 'react'
import { useTranslation } from 'react-i18next'

import { type model } from '@model'
import { Banner } from '@/ui/banner'

type ErrorContextType = [
  error: model.FGCTrackerError | null,
  setError: React.Dispatch<model.FGCTrackerError>
]
const ErrorContext = React.createContext<ErrorContextType | null>(null)
export const useErrorPopup = () => React.useContext(ErrorContext)![1]

export function ErrorPopupProvider(props: React.PropsWithChildren) {
  const { t } = useTranslation()
  const [error, setError] = React.useState<model.FGCTrackerError | null>(null)

  const dismiss = React.useCallback(() => setError(null), [])

  React.useEffect(() => {
    if (error === null) {
      return
    }
    const timer = window.setTimeout(dismiss, 15_000)
    return () => window.clearTimeout(timer)
  }, [error])

  return (
    <ErrorContext.Provider value={[error, setError]}>
      <Banner open={error !== null} variant='error' onClick={dismiss}>
        {error?.localizationKey && <span>{t(error.localizationKey)}</span>}
      </Banner>
      {props.children}
    </ErrorContext.Provider>
  )
}
