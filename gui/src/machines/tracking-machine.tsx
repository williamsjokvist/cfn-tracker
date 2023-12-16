import { assign, setup } from "xstate";
import { createActorContext } from "@xstate/react";

import {
  StartTracking,
  StopTracking,
} from "@@/go/core/CommandHandler";
import type { errorsx, model } from "@@/go/models";
import React from "react";
import { EventsOn } from "@@/runtime/runtime";

type TrackingMachineContextProps = {
  user?: model.User
  restore: boolean;
  isTracking: boolean;
  trackingState: model.TrackingState;
  error?: errorsx.FrontEndError;
};

export const TRACKING_MACHINE = setup({
  types: {
    context: {} as TrackingMachineContextProps,
  },
  actions: {
    startTracking: async ({ context, self }) => {
      if (!context.user || context.isTracking) return
      try {
        await StartTracking(context.user.code, context.restore)
        context.isTracking = true;
      } catch (error) {
        self.send({ type: "error", error })
      }
    },
    stopTracking: async ({ context, self }) => {
      try {
        await StopTracking();
        context.isTracking = false;
        self.send({ type: "cease" })
      } catch (error) {
        self.send({ type: "error", error })
      }
    },
  },
}).createMachine(
  {
    id: "cfn-tracker",
    context: {
      restore: false,
      isTracking: false,
      trackingState: {} as model.TrackingState,
    },
    initial: "cfnForm",
    states: {
      cfnForm: {
        on: {
          submit: {
            actions: [
              assign({
                user: ({ event }) => event.user,
                restore: ({ event }) => event.restore,
              }),
              "startTracking",
            ],
            target: "loading",
          },
        },
      },
      loading: {
        on: {
          matchPlayed: "tracking",
          error: {
            actions: [
              assign({
                error: ({ event }) => event.error,
              }),
            ],
            target: "cfnForm"
          }
        },
      },
      tracking: {
        on: {
          cease: {
            actions: "stopTracking",
            target: "cfnForm"
          },
          matchPlayed: {
            actions: assign({
              trackingState: ({ event }) => event.trackingState,
            }),
          },
        },
      },
    },
  },
);

export const TrackingMachineContext = createActorContext(TRACKING_MACHINE);  

export const TrackingMachineContextProvider = ({ children }) => {
  return (
    <TrackingMachineContext.Provider>
      <TrackingSubscriber>
        {children}
      </TrackingSubscriber>
    </TrackingMachineContext.Provider>
  )
}

const TrackingSubscriber = ({ children }) => {
  const trackingActor = TrackingMachineContext.useActorRef()
  React.useEffect(() => {
    EventsOn("cfn-data", (trackingState) => trackingActor.send({ type: "matchPlayed", trackingState }));
    EventsOn("stopped-tracking", () => trackingActor.send({ type: "cease" }));
  }, []);
  return children
}
