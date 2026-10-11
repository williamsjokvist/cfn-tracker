import { setup, assign } from 'xstate'
import { createActorContext } from '@xstate/react'

import { SelectGame } from '@cmd/TrackingHandler'
import { model } from '@model'
import { EventsOff, EventsOn } from '@runtime'

import { TRACKING_MACHINE } from './tracking-machine'

type AuthMachineContextProps = {
  game?: model.GameType
  error: model.FGCTrackerError | null
  action: { localizationKey: string; secondsLeft: number } | null
  step: string | null
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
    subscribeToAuthEvents: ({ self }) => {
      EventsOn('auth-success', () => self.send({ type: 'finished' }))
      EventsOn('auth-action-required', action => {
        self.send({ type: 'actionRequired', action })
      })
      EventsOn('auth-progress', step => self.send({ type: 'progress', step }))
    },
    unsubscribeToAuthEvents: () => {
      EventsOff('auth-success')
      EventsOff('auth-action-required')
      EventsOff('auth-progress')
    }
  }
}).createMachine({
  id: 'auth-machine',
  initial: 'gameForm',
  context: {
    error: null,
    action: null,
    step: null
  },
  states: {
    gameForm: {
      on: {
        submit: {
          actions: [
            assign({
              game: ({ event }) => event.game,
              error: null,
              action: null,
              step: null
            }),
            'selectGame',
            'subscribeToAuthEvents'
          ],
          target: 'loading'
        }
      }
    },
    loading: {
      on: {
        finished: {
          target: 'connected',
          actions: [
            'unsubscribeToAuthEvents',
            assign({
              error: null,
              action: null,
              step: null
            })
          ]
        },
        // A step after the login prompt means the user is done in the Chrome window.
        progress: {
          actions: [
            assign({
              step: ({ event }) => event.step,
              action: null
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
              action: null,
              step: null
            }),
            'unsubscribeToAuthEvents'
          ],
          target: 'gameForm'
        },
        // The login keeps running in the backend, but nothing listens for it anymore.
        pickGame: {
          actions: ['unsubscribeToAuthEvents', assign({ action: null, step: null })],
          target: 'gameForm'
        }
      }
    },
    connected: {
      invoke: {
        id: 'cfn-tracker',
        src: TRACKING_MACHINE
      },
      on: {
        pickGame: { target: 'gameForm' }
      }
    }
  }
})

export const AuthMachineContext = createActorContext(AUTH_MACHINE)
