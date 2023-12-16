import { TRACKING_MACHINE } from "@/machines/tracking-machine";
import { SelectGame } from "@@/go/core/CommandHandler";
import { errorsx } from "@@/go/models";
import { createActorContext } from "@xstate/react";
import { setup, assign } from "xstate";

type AuthMachineContext = {
  game?: "sfv" | "sf6";
  error?: errorsx.FrontEndError;
}
export const AUTH_MACHINE = setup({
  types: {
    context: <AuthMachineContext>{}
  },
  actions: {
    selectGame: ({ context, self }) => {
      SelectGame(context.game).then(() => {
        self.send({ type: "loadedGame" })
      }).catch(error => self.send({ type: "error", error }));
    },
  }
}).createMachine({
  id: "auth-machine",
  initial: "gameForm",
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
        loadedGame: "connected",
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
      invoke: {
        id: "cfn-tracker",
        src: TRACKING_MACHINE
      }
    }
  }
})

export const AuthMachineContext = createActorContext(AUTH_MACHINE);  
