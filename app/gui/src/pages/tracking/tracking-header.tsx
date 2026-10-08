import React from 'react'
import { useSelector } from '@xstate/react'
import { useTranslation } from 'react-i18next'

import { TrackingMachineContext } from '@/state/tracking-machine'
import { type LocalizationKey } from '@/main/i18n'
import * as Page from '@/ui/page'

export function TrackingHeader() {
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

  const retrying = state === 'retrying' && retry

  return (
    <Page.Header>
      <div className='flex items-baseline gap-4'>
        <Page.Title>{retrying ? t('retrying') : t('tracking')}</Page.Title>
        {retrying && (
          <span className='text-sm whitespace-nowrap text-amber-200' role='status'>
            {t(retry.reason as LocalizationKey)}
            {retry.nextRetryInMs > 0 && (
              <> · {t('retryingDetail', { attempt: retry.attempt, seconds })}</>
            )}
          </span>
        )}
      </div>
      <Page.LoadingIcon />
    </Page.Header>
  )
}
