import { setup, assign } from 'xstate'
import { createActorContext } from '@xstate/react'

import { SelectGame } from '@cmd/TrackingHandler'
import type { errorsx } from '@model'
import { EventsOff, EventsOn } from '@runtime'

import { TRACKING_MACHINE } from './tracking-machine'

type AuthMachineContextProps = {
  progress: number
  game?: 'sfv' | 'sf6'
  error: errorsx.FormattedError | null
}
export const AUTH_MACHINE = setup({
  types: {
    context: <AuthMachineContextProps>{}
  },
  actions: {
    selectGame: ({ context, self }) => {
      SelectGame(context.game ?? 'sf6').catch(error => self.send({ type: 'error', error }))
    },
    subscribeToProgressEvents: ({ self }) => {
      EventsOn('auth-progress', progress => {
        self.send({ type: 'loaded', progress })
        if (progress >= 100) {
          self.send({ type: 'finished' })
        }
      })
    },
    unsubscribeToProgressEvents: () => {
      EventsOff('auth-progress')
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
    error: null
  },
  states: {
    gameForm: {
      on: {
        submit: {
          actions: [
            assign({
              game: ({ event }) => event.game,
              error: null
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
              error: null
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
        error: {
          actions: [
            assign({
              error: ({ event }) => event.error,
              progress: 0
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
