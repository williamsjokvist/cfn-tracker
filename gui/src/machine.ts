import { createMachine, assign, EventObject, send } from 'xstate';
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
    matchHistory: <core.MatchHistory>{
      cfn: "",
      wins: 0,
      losses: 0,
      winRate: 0,
      lpGain: 0,
      lp: 0,
      opponent: "",
      opponentCharacter: "",
      opponentLeague: "",
      opponentLP: 0,
      totalLosses: 0,
      totalMatches: 0,
      totalWins: 0,
      result: false,
      winStreak: 0,
      timestamp: "",
      date: "",
    }
  },
  initial: 'idle',
  states: {
    loading: {
      on: {
        initialized: "idle",
      }
    },

    idle: {
      entry: 'stopTracking',
      on: {
        submit: {
          actions: assign({
            cfn: (ctx, e: any) => e.cfn,
            restore: (ctx, e: any) => e.restore,
          }),
          target: 'loadingCfn'
        },
      },
    },

    loadingCfn: {
      entry: "startTracking",
      on: {
        startedTracking: 'tracking'
      }
    },

    tracking: {
      on: {
        stoppedTracking: "idle",
        matchPlayed: {
          actions: assign({
            matchHistory: (_, e: any) => e.matchHistory
          })
        }
      },

      exit: "stopTracking",
    },

  },
}, {
  actions: {
    startTracking: ({cfn, restore, isTracking}) => {
      if (cfn && !isTracking) {
        StartTracking(cfn, restore)
        isTracking = true
      }
    },
    stopTracking: ({ isTracking }) => {
      if (!isTracking) return
      StopTracking().then(_ => {
        isTracking = false
      })
    }
  }
});

export const CFNMachineContext = createActorContext(cfnMachine);