import { createMachine, assign, EventObject } from 'xstate';
import {
  StopTracking,
  IsTracking,
  IsInitialized,
  ResultsJSONExist,
  GetAvailableLogs,
  StartTracking,
} from "@@/go/core/App";
import { core } from "@@/go/models";

const startTracking = (cfn: string, restoreData: boolean) => {
  StartTracking(cfn, restoreData)
}

export type TrackEvent = {
  cfn: string; 
  restore: boolean;
} & EventObject


export type MatchEvent = {
  matchHistory: core.MatchHistory
} & EventObject


export const cfnMachine = createMachine({
  id: 'cfn-tracker',
  initial: 'initialized',
  context: {
    cfn: '',
    restore: false,
    matchHistory: <core.MatchHistory>null,
  },
  states: {
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
  }
});

const machine = createMachine({
  id: 'cfn-tracker',
  initial: 'initialized',
  context: {
    cfn: '',
    restore: false,
    matchHistory: <core.MatchHistory>null,
  },
  states: {
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
  }
});