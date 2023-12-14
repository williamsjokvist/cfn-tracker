import { assign, EventObject, setup } from "xstate";
import { createActorContext } from "@xstate/react";

import {
  StartTracking,
  StopTracking,
  SelectGame,
} from "@@/go/core/CommandHandler";
import type { model } from "@@/go/models";

export type MatchEvent = {
  trackingState: model.TrackingState;
} & EventObject;

type CFNMachineContext = {
  user?: model.User;
  game?: "sfv" | "sf6";
  restore: boolean;
  isTracking: boolean;
  trackingState: model.TrackingState;
};

export const TRACKING_MACHINE = setup({
  types: {
    context: <CFNMachineContext>{},
  },
  actions: {
    initialize: ({ context, self }) => {
      SelectGame(context.game).catch(err => self.send({ type: "error", err }));
    },
    startTracking: async ({ context, self }) => {
      if (!context.user || context.isTracking) return
      try {
        await StartTracking(context.user.code, context.restore)
        context.isTracking = true;
        self.send({ type: "startedTracking" })
      } catch (err) {
        self.send({ type: "error", err })
      }
    },
    stopTracking: async ({ context, self }) => {
      if (!context.isTracking) return;
      try {
        await StopTracking();
        context.isTracking = false;
        self.send({ type: "stoppedTracking" })
      } catch (err) {
        self.send({ type: "error", err })
      }
    },
  },
}).createMachine(
  {
    id: "cfn-tracker",
    context: {
      restore: false,
      isTracking: false,
      trackingState: <model.TrackingState>{},
    },
    initial: "formGame",
    states: {
      "formGame": {
        on: {
          submit: {
            actions: assign({
              game: ({ event }) => event.game,
            }),
            target: "loadingGame",
          },
        },
      },
      "loadingGame": {
        entry: "initialize",
        on: {
          loadedGame: "formCfn",
          error: "formGame"
        },
      },
      "formCfn": {
        on: {
          submit: {
            actions: assign({
              user: ({ event }) => event.user,
              restore: ({ event }) => event.restore,
            }),
            target: "loadingCfn",
          },
        },
      },
      "loadingCfn": {
        entry: "startTracking",
        on: {
          startedTracking: "tracking",
          error: "formCfn"
        },
      },
      tracking: {
        on: {
          stoppedTracking: "formCfn",
          matchPlayed: {
            actions: assign({
              trackingState: ({ event }) => event.trackingState,
            }),
          },
        },
        exit: "stopTracking",
      },
    },
  },
);

