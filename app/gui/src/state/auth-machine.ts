import { setup, assign } from 'xstate'
import { createActorContext } from '@xstate/react'

import { SelectGame } from '@cmd/TrackingHandler'
import { model } from '@model'
import { EventsOff, EventsOn } from '@runtime'

import { TRACKING_MACHINE } from './tracking-machine'

type AuthMachineContextProps = {
  progress: number
  game?: model.GameType
  error: model.FGCTrackerError | null
  action: { localizationKey: string; secondsLeft: number } | null
}
export const AUTH_MACHINE = setup({
  types: {
    context: <AuthMachineContextProps>{}
  },
  actions: {
    selectGame: ({ context, self }) => {
      if (context.game) {
        SelectGame(context.game).catch(error => self.send({ type: 'error', error }))
      }
    },
    subscribeToProgressEvents: ({ self }) => {
      EventsOn('auth-progress', progress => {
        self.send({ type: 'loaded', progress })
        if (progress >= 100) {
          self.send({ type: 'finished' })
        }
      })
      EventsOn('auth-action-required', action => {
        self.send({ type: 'actionRequired', action })
      })
    },
    unsubscribeToProgressEvents: () => {
      EventsOff('auth-progress')
      EventsOff('auth-action-required')
    }
  },
  guards: {
    isLoaded: ({ context }) => context.progress >= 100
  }
}).createMachine({
  id: 'auth-machine',
  initial: 'gameForm',
  context: {
    progress: 0,
    error: null,
    action: null
  },
  states: {
    gameForm: {
      on: {
        submit: {
          actions: [
            assign({
              game: ({ event }) => event.game,
              error: null,
              action: null
            }),
            'selectGame',
            'subscribeToProgressEvents'
          ],
          target: 'loading'
        }
      }
    },
    loading: {
      on: {
        finished: {
          target: 'connected',
          guard: 'isLoaded',
          actions: [
            'unsubscribeToProgressEvents',
            assign({
              progress: 0,
              error: null,
              action: null
            })
          ]
        },
        loaded: {
          actions: [
            assign({
              progress: ({ event }) => event.progress
            })
          ]
        },
        actionRequired: {
          actions: [
            assign({
              action: ({ event }) => event.action
            })
          ]
        },
        error: {
          actions: [
            assign({
              error: ({ event }) => event.error,
              progress: 0,
              action: null
            }),
            'unsubscribeToProgressEvents'
          ],
          target: 'gameForm'
        }
      }
    },
    connected: {
      invoke: {
        id: 'cfn-tracker',
        src: TRACKING_MACHINE
      }
    }
  }
})

export const AuthMachineContext = createActorContext(AUTH_MACHINE)
