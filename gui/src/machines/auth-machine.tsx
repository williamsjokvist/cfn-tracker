import { TRACKING_MACHINE } from "@/machines/tracking-machine";
import { SelectGame } from "@@/go/core/CommandHandler";
import { errorsx } from "@@/go/models";
import { EventsOn } from "@@/runtime/runtime";
import { createActorContext } from "@xstate/react";
import React from "react";
import { setup, assign } from "xstate";

type AuthMachineContextProps = {
  loaded: number
  game?: "sfv" | "sf6"
  error?: errorsx.FrontEndError;
}
export const AUTH_MACHINE = setup({
  types: {
    context: {} as AuthMachineContextProps
  },
  actions: {
    selectGame: ({ context, self }) => {
      SelectGame(context.game).then(() => {
        self.send({ type: "loadedGame" })
      }).catch(error => self.send({ type: "error", error }));
    },
  },
  guards: {
    isLoaded: ({ context }) => context.loaded >= 100
  }
}).createMachine({
  id: "auth-machine",
  initial: "gameForm",
  context: {
    loaded: 0
  },
  states: {
    gameForm: {
      on: {
        submit: {
          actions: [
            assign({
              game: ({ event }) => event.game,
            }),
            "selectGame",
          ],
          target: "loading",
        },
      },
    },
    loading: {
      on: {
        loaded: [
          {
            actions: [
              assign({
                loaded: ({ event }) => event.loaded
              }),
            ],
            target: "connected",
            guard: "isLoaded"
          },
          {
            actions: [
              assign({
                loaded: ({ event }) => event.loaded
              }),
            ],
          },
        ],
        error: {
          actions: [
            assign({
              error: ({ event }) => event.error,
            }),
          ],
          target: "gameForm"
        }
      },
    },
    connected: {
      always: {
         guard: "isLoaded",
         actions: [
           assign({
             loaded: 0
           }),
         ],
       },
      invoke: {
        id: "cfn-tracker",
        src: TRACKING_MACHINE
      }
    }
  },
})

export const AuthMachineContext = createActorContext(AUTH_MACHINE);

export const AuthMachineContextProvider = ({ children }) => {
  return (
    <AuthMachineContext.Provider>
      <AuthSubscriber>
        {children}
      </AuthSubscriber>
    </AuthMachineContext.Provider>
  )
}

const AuthSubscriber = ({ children }) => {
  const authActor = AuthMachineContext.useActorRef()
  React.useEffect(() => {
    EventsOn("auth-loaded", (loaded) =>  authActor.send({ type: "loaded", loaded }));
  }, [])
  return children
}
