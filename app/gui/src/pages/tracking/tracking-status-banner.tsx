import React from 'react'
import { useSelector } from '@xstate/react'
import { useTranslation } from 'react-i18next'

import { TrackingMachineContext } from '@/state/tracking-machine'
import { type LocalizationKey } from '@/main/i18n'

export function TrackingStatusBanner() {
  const { t } = useTranslation()
  const actor = TrackingMachineContext.useActorRef()
  const state = useSelector(actor, snapshot => snapshot.value)
  const retry = useSelector(actor, snapshot => snapshot.context.retry)
  const [seconds, setSeconds] = React.useState(0)

  React.useEffect(() => {
    if (!retry) return
    const deadline = Date.now() + retry.nextRetryInMs
    const update = () => setSeconds(Math.max(0, Math.ceil((deadline - Date.now()) / 1000)))
    update()
    const timer = window.setInterval(update, 1000)
    return () => window.clearInterval(timer)
  }, [retry])

  if (state === 'retrying' && retry) {
    return (
      <div className='mx-6 mt-3 rounded-lg bg-amber-500/20 px-4 py-2 text-amber-100' role='status'>
        {retry.nextRetryInMs > 0 && (
          <>{t('retryingDetail', { attempt: retry.attempt, seconds })} — </>
        )}
        {t(retry.reason as LocalizationKey)}
      </div>
    )
  }
  if (state === 'tracking') {
    return <div className='mx-6 mt-3 rounded-lg bg-emerald-500/15 px-4 py-2'>{t('tracking')}</div>
  }
  return null
}
