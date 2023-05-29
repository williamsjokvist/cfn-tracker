import { createMachine, assign } from 'xstate';
import {
  StopTracking,
  OpenResultsDirectory,
  IsTracking,
  IsInitialized,
  ResultsJSONExist,
  GetAvailableLogs,
  StartTracking,
} from "@@/go/core/App";

const startTracking = async (cfn: string, restoreData: boolean) => {
  StartTracking(cfn, restoreData)
}

export const createCfnMachine = (cfn: string) => {
  return createMachine({
    /** @xstate-layout N4IgpgJg5mDOIC5QGMBmA7AdAGwPYEMIBLdKAYgl3TE1gBd86a0s9CSoBtABgF1FQAB1ywidIlQEgAHogC0AJgBsAdkzKArBoAsARm5KNAZm2qANCACeiABy7MSvbqMmb2uxoXaAvt4stMEjEifGwiAC9ICioaekZmDED0YNCIyB5+JBBhUXFJLNkERRUjTB1uDRVuBTsaowULawQFI24HAE4Faq1uXW0K7QVff0S6ACd8ZABrDjIABQBBAFUAZQBRDKkc4PzQQuKFdRslVV1ddqUbL10lRsRtI3sWlTt+9pcDDWGQAPHJmdIZBWABUAPJzTZZbZ5dBSQoKdqYEwqJS6GwvRxGWp3ZoIzDnLoKfRnFQaJQ1b4BQT4ACusCiACU1islgBZDZ8LYiHawgryTqYbjcFRdTpC8mPXQ4oxKNrad7VckqC7CvqUxLUulREHgyFCbkwuGIZz2bg2DS6F4ioWdZw4hR4gnVYmWslGdVYVD4IjYRlrYEMgCaeuyBokvL2iG4OO4vj8IHQuAgcCkLC5uXDRqK9VKGhsRh0ug0nzJOLkRaUDk81Ramla7S+8YCbGIpHTPKzB0w2hewpU8vaLzcKjLpRr7XOFsq-SM7RsHqSKTCkQg7cNfKKCn76gMDeL73q-W0OL6iJUZzzAwRNnzC7+0w4a8zG+KGjK+cLxcMpasiElmAVfRPjOap5ybDVaXpVcoTDXYZHkc0OmFG8RQedo9BHX8EEGewHmQgtBgqd1wM9b1fWg-UMzg-ZnFKZQa0qA9TjLJRESLQlUXQ5QGzA3wgA */
    id: 'cfn',
    initial: 'loading',
    context: {
      cfn,
      matchHistory: null,
    },
    states: {
      loading: {
        onDone: {
          target: 'initialized'
        },
      },
      initialized: {
        onDone: {
          target: 'tracking'
        }
      },
      tracking: {
        on: {
          PAUSE: 'paused',
          STOP: 'initialized'
        }
      },
      paused: {
        on: {
          RESUME: 'tracking',
          STOP: 'initialized'
        }
      },
      failed: {
        on: {
          RETRY: 'loading'
        }
      }
    }
  });
}