import { assign, setup } from 'xstate'
import { createActorContext } from '@xstate/react'

import { ForcePoll, StartTracking, StopTracking } from '@cmd'
import type { errorsx, model } from '@model'
import { EventsOff, EventsOn } from '@runtime'

type TrackingMachineContextProps = {
  user: model.User | null
  restore: boolean
  isTracking: boolean
  trackingState: model.TrackingState
  error: errorsx.FormattedError | null
}

export const TRACKING_MACHINE = setup({
  types: {
    context: <TrackingMachineContextProps>{}
  },
  guards: {
    notTracking: ({ context }) => !context.isTracking
  },
  actions: {
    startTracking: ({ context, self }) => {
      context.user &&
        StartTracking(context.user.code, context.restore).catch(error =>
          self.send({ type: 'error', error })
        )
    },
    stopTracking: ({ self }) => {
      StopTracking().catch(error => self.send({ type: 'error', error }))
    },
    forcePoll: () => {
      ForcePoll()
    },
    subscribeToTrackingEvents: ({ self }) => {
      EventsOn('cfn-data', trackingState => self.send({ type: 'matchPlayed', trackingState }))
      EventsOn('stopped-tracking', () => self.send({ type: 'cease' }))
    },
    unsubscribeToTrackingEvents: ({ self }) => {
      EventsOff('cfn-data')
      EventsOff('stopped-tracking')
    }
  }
}).createMachine({
  id: 'cfn-tracker',
  context: {
    user: null,
    error: null,
    restore: false,
    isTracking: false,
    trackingState: <model.TrackingState>{}
  },
  initial: 'cfnForm',
  states: {
    cfnForm: {
      on: {
        submit: {
          guard: 'notTracking',
          actions: [
            assign({
              user: ({ event }) => event.user,
              restore: ({ event }) => event.restore,
              isTracking: true,
              error: null
            }),
            'startTracking',
            'subscribeToTrackingEvents'
          ],
          target: 'loading'
        }
      }
    },
    loading: {
      on: {
        matchPlayed: {
          actions: assign({
            trackingState: ({ event }) => event.trackingState
          }),
          target: 'tracking'
        },
        error: {
          actions: [
            assign({
              error: ({ event }) => event.error,
              isTracking: false
            }),
            'unsubscribeToTrackingEvents'
          ],
          target: 'cfnForm'
        }
      }
    },
    tracking: {
      on: {
        forcePoll: {
          actions: ['forcePoll']
        },
        cease: {
          actions: [
            'stopTracking',
            'unsubscribeToTrackingEvents',
            assign({
              isTracking: false
            })
          ],
          target: 'cfnForm'
        },
        matchPlayed: {
          actions: assign({
            trackingState: ({ event }) => event.trackingState
          })
        }
      }
    }
  }
})

export const TrackingMachineContext = createActorContext(TRACKING_MACHINE)
