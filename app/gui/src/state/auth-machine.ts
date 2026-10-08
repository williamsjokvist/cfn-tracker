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
    },
    unsubscribeToAuthEvents: () => {
      EventsOff('auth-success')
      EventsOff('auth-action-required')
    }
  }
}).createMachine({
  id: 'auth-machine',
  initial: 'gameForm',
  context: {
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
              action: null
            }),
            'unsubscribeToAuthEvents'
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
