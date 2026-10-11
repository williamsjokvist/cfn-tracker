import React from 'react'
import { useTranslation } from 'react-i18next'

import { EventsOff, EventsOn } from '@runtime'
import type { LocalizationKey } from '@/main/i18n'
import { Banner } from '@/ui/banner'

// Shown app-wide, as the backend reports some things after the user has left the
// page they're about, like a comparison stopped by leaving it.
export function InfoPopup() {
  const { t } = useTranslation()
  const [localizationKey, setLocalizationKey] = React.useState<LocalizationKey | null>(null)

  React.useEffect(() => {
    EventsOn('info', setLocalizationKey)
    return () => EventsOff('info')
  }, [])

  React.useEffect(() => {
    if (localizationKey === null) {
      return
    }
    const timer = window.setTimeout(() => setLocalizationKey(null), 5_000)
    return () => window.clearTimeout(timer)
  }, [localizationKey])

  return (
    <Banner open={localizationKey !== null} variant='info' onClick={() => setLocalizationKey(null)}>
      {localizationKey && <span>{t(localizationKey)}</span>}
    </Banner>
  )
}
