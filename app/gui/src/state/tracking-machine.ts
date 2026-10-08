import { assign, setup } from 'xstate'
import { createActorContext } from '@xstate/react'

import { ForcePoll, StartTracking, StopTracking } from '@cmd/TrackingHandler'
import type { model } from '@model'
import { EventsOff, EventsOn } from '@runtime'

type TrackingMachineContextProps = {
  user: model.User | null
  restore: boolean
  isTracking: boolean
  match: model.Match
  error: model.FGCTrackerError | null
  retry: { attempt: number; nextRetryInMs: number; reason: string } | null
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
    forcePoll: ForcePoll,
    subscribeToTrackingEvents: ({ self }) => {
      EventsOn('match', match => self.send({ type: 'matchPlayed', match }))
      EventsOn('stopped-tracking', () => self.send({ type: 'cease' }))
      EventsOn('tracking-retrying', retry => self.send({ type: 'retrying', retry }))
      EventsOn('tracking-recovered', () => self.send({ type: 'recovered' }))
      EventsOn('tracking-error', error => self.send({ type: 'error', error }))
    },
    unsubscribeToTrackingEvents: ({ self }) => {
      EventsOff('match')
      EventsOff('stopped-tracking')
      EventsOff('tracking-retrying')
      EventsOff('tracking-recovered')
      EventsOff('tracking-error')
    }
  }
}).createMachine({
  id: 'cfn-tracker',
  context: {
    user: null,
    error: null,
    restore: false,
    isTracking: false,
    match: <model.Match>{},
    retry: null
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
        retrying: {
          actions: assign({ retry: ({ event }) => event.retry }),
          target: 'retrying'
        },
        matchPlayed: {
          actions: assign({
            match: ({ event }) => event.match
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
        retrying: {
          actions: assign({ retry: ({ event }) => event.retry }),
          target: 'retrying'
        },
        error: {
          actions: [
            assign({ error: ({ event }) => event.error, isTracking: false, retry: null }),
            'unsubscribeToTrackingEvents'
          ],
          target: 'cfnForm'
        },
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
            match: ({ event }) => event.match
          })
        }
      }
    },
    retrying: {
      on: {
        retrying: { actions: assign({ retry: ({ event }) => event.retry }) },
        recovered: { actions: assign({ retry: null }), target: 'tracking' },
        matchPlayed: {
          actions: assign({ match: ({ event }) => event.match, retry: null }),
          target: 'tracking'
        },
        forcePoll: { actions: ['forcePoll'] },
        error: {
          actions: [
            assign({ error: ({ event }) => event.error, isTracking: false, retry: null }),
            'unsubscribeToTrackingEvents'
          ],
          target: 'cfnForm'
        },
        cease: {
          actions: [
            'stopTracking',
            'unsubscribeToTrackingEvents',
            assign({ isTracking: false, retry: null })
          ],
          target: 'cfnForm'
        }
      }
    }
  }
})

export const TrackingMachineContext = createActorContext(TRACKING_MACHINE)
