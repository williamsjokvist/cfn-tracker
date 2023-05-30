import { createMachine, assign, EventObject } from 'xstate';
import { createActorContext } from '@xstate/react';

import {
  StartTracking,
  StopTracking
} from "@@/go/core/App";
import { core } from "@@/go/models";

export type TrackEvent = {
  cfn: string; 
  restore: boolean;
} & EventObject


export type MatchEvent = {
  matchHistory: core.MatchHistory
} & EventObject


export const cfnMachine = createMachine({
  /** @xstate-layout N4IgpgJg5mDOIC5QGMBmA7AtAFwE4ENkBrMXAOgBsB7fCAS3SgGIG7s78K6AvSAbQAMAXUSgADlVhs6VdKJAAPRJgBMZAMzqVA9QEYVAVgA0IAJ7KDa9QHZdAgJw2DAX2cm0WPIRLk6ECmBMsACuAEYAtmyCIkggElLssvJKCNZkACz2AGwOTibmqQZkWenpBnqGru4YOATEpGRexAzMsNhUYtHy8dJJsSn66RnWKlkGWYb5iDlk4yUCEy5uIB613g1NRC1M4fjYyAAWAAoU+Kb8wt2SvXL9iAAcZPp2jtbGZogq97pkKurpTlcy3QVAgcHkq02pCuCRkt1AKUwBgEGi0On07wKqgEVlsAnu1nU9xUX3SOiqKxqUPI1FoLRhN2SiB+93s9wBbymCEwul0QwJL0By0hdR8ZD8AQZiXhimZ9mK9j5C0mHwQdg0+IMryW1U8oo2ovpsR60qZCHUv3uBl04xVBV06nlmRUugJOspevW5GCYgge0gAFk9ocABJ0NpUXAFcTXU13BB-NTs+4ExzaN72Tmq3RvMj2ewutl-G34+xA5xAA */
  id: 'cfn-tracker',
  predictableActionArguments: true,
  
  context: {
    cfn: '',
    restore: false,
    isTracking: false,
    matchHistory: <core.MatchHistory>null
  },

  states: {
    loading: {
      on: {
        initialized: "idle"
      }
    },

    idle: {
      on: {
        submit: "fetchingCfn",
      }
    },

    fetchingCfn: {
      on: {
        success: 'tracking',
        // todo:  add fail state
      }
    },

    tracking: {
      on: {
        stop: "idle",
      },

      entry: "startTracking",
      exit: "stopTracking",
    },
  },

  initial: "loading"
}, {
  actions: {
    startTracking: ({cfn, restore, isTracking}) => {
      if (!isTracking) 
        StartTracking(cfn, restore)
    },
    stopTracking: () => StopTracking()
  }
});

export const CFNMachineContext = createActorContext(cfnMachine);

/*

    loading: { 
      on: {
        INITIALIZED: {
          target: 'initialized'
        }
      }
    },
    initialized: {
      on: {
        TRACK: {
          target: 'tracking',
          actions: assign({
            cfn: (ctx, event: TrackEvent) => event.cfn,
            restore: (ctx, event: TrackEvent) => event.restore
          })
        }
      }
    },
    tracking: {
      on: {
        STOP: 'initialized',
        MATCH: {
          actions: assign({
            matchHistory: (ctx, event: MatchEvent) => event.matchHistory,
          })
        }
      }
    },
*/