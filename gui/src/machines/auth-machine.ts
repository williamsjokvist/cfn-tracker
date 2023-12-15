import { TRACKING_MACHINE } from "@/machines/tracking-machine";
import { SelectGame } from "@@/go/core/CommandHandler";
import { createActorContext } from "@xstate/react";
import { setup, assign } from "xstate";

type AuthMachineContext = {
  game?: "sfv" | "sf6"
}
export const AUTH_MACHINE = setup({
  types: {
    context: <AuthMachineContext>{}
  },
  actions: {
    selectGame: ({ context, self }) => {
      SelectGame(context.game).then(() => {
        self.send({ type: "loadedGame" })
      }).catch(err => self.send({ type: "error", err }));
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
        error: "gameForm"
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
