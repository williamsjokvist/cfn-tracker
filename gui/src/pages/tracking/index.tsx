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
            <Page.Title>{t('loading')}</Page.Title>
            <Page.LoadingIcon />
          </Page.Header>
        </Page.Root>
      )
  }

  switch (trackingState) {
    case 'cfnForm':
      return <TrackingForm />
    case 'tracking':
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
