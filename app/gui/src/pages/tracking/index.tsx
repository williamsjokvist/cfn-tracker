import React from 'react'
import { useSelector } from '@xstate/react'
import { useTranslation } from 'react-i18next'

import { TrackingMachineContext } from '@/state/tracking-machine'
import { AuthMachineContext } from '@/state/auth-machine'
import { useErrorPopup } from '@/main/error-popup'
import * as Page from '@/ui/page'

import { TrackingForm } from './tracking-form'
import { TrackingGamePicker } from './tracking-game-picker'
import { TrackingLiveUpdater } from './tracking-live-updater'

export function TrackingPage() {
  const { t } = useTranslation()

  const trackingActor = TrackingMachineContext.useActorRef()
  const authActor = AuthMachineContext.useActorRef()

  const authState = useSelector(authActor, ({ value }) => value)
  const trackingState = useSelector(trackingActor, ({ value }) => value)

  const authError = useSelector(authActor, ({ context }) => context.error)
  const authAction = useSelector(authActor, ({ context }) => context.action)
  const trackingError = useSelector(trackingActor, ({ context }) => context.error)

  const setError = useErrorPopup()

  React.useEffect(() => {
    authError && setError(authError)
  }, [authError])

  React.useEffect(() => {
    trackingError && setError(trackingError)
  }, [trackingError])

  switch (authState) {
    case 'gameForm':
      return <TrackingGamePicker onSubmit={game => authActor.send({ type: 'submit', game })} />
    case 'loading':
      return (
        <Page.Root>
          <Page.Header>
            <Page.Title>
              {authAction
                ? t(authAction.localizationKey as 'authNeedRelogin', {
                    seconds: authAction.secondsLeft
                  })
                : t('loading')}
            </Page.Title>
            {!authAction && <Page.LoadingIcon />}
          </Page.Header>
          {authAction && (
            <div className='flex flex-col items-center justify-center gap-6 px-8 text-center'>
              <i
                aria-label='loading'
                className='text-highlight inline-block h-12 w-12 animate-spin rounded-full border-[4px] border-current border-t-transparent'
                role='status'
              />
              <p className='max-w-sm text-white/70'>{t('authNeedReloginHint')}</p>
            </div>
          )}
        </Page.Root>
      )
  }

  switch (trackingState) {
    case 'cfnForm':
      return (
        <>
          {trackingError && (
            <div className='mx-6 mt-3 rounded-lg bg-red-500/20 px-4 py-2 text-red-100'>
              {t(trackingError.localizationKey)}
            </div>
          )}
          <TrackingForm />
        </>
      )
    case 'tracking':
    case 'retrying':
      return <TrackingLiveUpdater />
    case 'loading':
    default:
      return (
        <Page.Root>
          <Page.Header>
            <Page.Title>{t('loading')}</Page.Title>
            <Page.LoadingIcon />
          </Page.Header>
        </Page.Root>
      )
  }
}
