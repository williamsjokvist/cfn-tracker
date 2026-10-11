import { Navigate, RouterProvider as ReactRouterProvider, createHashRouter } from 'react-router-dom'

import { SettingsPage } from '@/pages/settings'
import { OutputPage } from '@/pages/output'
import { MatchesListPage } from '@/pages/matches'
import { SessionsListPage } from '@/pages/sessions'
import { SessionOverviewPage, type SessionData } from '@/pages/session'
import { TrackingPage } from '@/pages/tracking'

import { AppWrapper } from './app-wrapper'
import { AppErrorBoundary, PageErrorBoundary } from './app-error'

import { GetSession, GetMatches } from '@cmd/CommandHandler'

async function loadSession(sessionId: number, page = 0, limit = 0): Promise<SessionData> {
  const [session, matches] = await Promise.all([
    GetSession(sessionId),
    GetMatches(sessionId, '', page, limit)
  ])
  return { session, matches: matches ?? [] }
}

const router = createHashRouter([
  {
    element: <AppWrapper />,
    errorElement: <AppErrorBoundary />,
    children: [
      {
        errorElement: <PageErrorBoundary />,
        children: [
          {
            element: <Navigate to='tracking' />,
            index: true
          },
          {
            element: <TrackingPage />,
            path: '/tracking'
          },
          {
            element: <OutputPage />,
            path: '/output'
          },
          {
            element: <SettingsPage />,
            path: '/settings'
          },
          {
            element: <SessionsListPage />,
            path: '/sessions'
          },
          {
            element: <SessionOverviewPage />,
            path: '/sessions/:sessionId',
            loader: ({ params }) => loadSession(Number(params.sessionId))
          },
          {
            element: <MatchesListPage />,
            path: '/sessions/:sessionId/matches/:page?/:limit?',
            loader: ({ params }) =>
              loadSession(
                Number(params.sessionId),
                Number(params.page ?? 0),
                Number(params.limit ?? 0)
              )
          }
        ]
      }
    ]
  }
])

export const RouterProvider = () => <ReactRouterProvider router={router} />
