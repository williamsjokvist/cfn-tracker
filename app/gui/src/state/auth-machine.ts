import { setup, assign } from 'xstate'
import { createActorContext } from '@xstate/react'

import { CancelAuthentication, SelectGame } from '@cmd/TrackingHandler'
import { model } from '@model'
import { EventsOff, EventsOn } from '@runtime'

import { TRACKING_MACHINE } from './tracking-machine'

type AuthMachineContextProps = {
  game?: model.GameType
  error: model.FGCTrackerError | null
  action: { localizationKey: string; secondsLeft: number } | null
  step: string | null
}
// Only the latest login may report an error, as a cancelled one fails after the
// user has already picked another game.
let latestLogin = 0

export const AUTH_MACHINE = setup({
  types: {
    context: <AuthMachineContextProps>{}
  },
  actions: {
    selectGame: ({ context, self }) => {
      if (context.game) {
        const login = ++latestLogin
        SelectGame(context.game).catch(
          error => login === latestLogin && self.send({ type: 'error', error })
        )
      }
    },
    cancelLogin: () => {
      latestLogin++
      CancelAuthentication()
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
        pickGame: {
          actions: ['cancelLogin', 'unsubscribeToAuthEvents', assign({ action: null, step: null })],
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
