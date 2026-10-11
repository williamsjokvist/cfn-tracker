import React from 'react'
import { useSelector } from '@xstate/react'
import { useTranslation } from 'react-i18next'

import { AuthMachineContext } from '@/state/auth-machine'
import { TrackingMachineContext } from '@/state/tracking-machine'
import * as Page from '@/ui/page'
import * as Dialog from '@/ui/dialog'
import { Button } from '@/ui/button'

import { gameName } from './tracking-game-picker'

export function TrackingBreadcrumbs({ current }: { current?: string }) {
  const { t } = useTranslation()
  const authActor = AuthMachineContext.useActorRef()
  const trackingActor = TrackingMachineContext.useActorRef()
  const game = useSelector(authActor, ({ context }) => context.game)
  const isTracking = useSelector(trackingActor, ({ context }) => context.isTracking)
  const [confirming, setConfirming] = React.useState(false)

  const pickGame = () => {
    if (isTracking) {
      trackingActor.send({ type: 'cease' })
    }
    authActor.send({ type: 'pickGame' })
  }

  return (
    <>
      <Page.Breadcrumbs
        items={[
          { label: t('pickGame'), onClick: () => (isTracking ? setConfirming(true) : pickGame()) },
          { label: gameName(game) },
          ...(current ? [{ label: current }] : [])
        ]}
      />
      <Dialog.Root open={confirming} onOpenChange={setConfirming}>
        <Dialog.Content
          title='stopTrackingTitle'
          description='stopTrackingDescription'
          className='h-auto!'
        >
          <div className='mt-8 flex justify-end'>
            <Button
              onClick={() => {
                setConfirming(false)
                pickGame()
              }}
            >
              {t('stop')}
            </Button>
          </div>
        </Dialog.Content>
      </Dialog.Root>
    </>
  )
}
