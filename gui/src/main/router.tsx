import { Navigate, RouterProvider as ReactRouterProvider, createHashRouter } from 'react-router-dom'

import { OutputPage } from '@/pages/output/output-page'
import { MatchesListPage } from '@/pages/stats/matches-list-page'
import { SessionsListPage } from '@/pages/stats/sessions-list-page'
import { TrackingPage } from '@/pages/tracking/tracking-page'

import { AppWrapper } from './app-layout/app-wrapper'
import { AppErrorBoundary, PageErrorBoundary } from './app-layout/app-error'

import { GetUsers, GetSessions, GetMatches } from '@@/go/core/CommandHandler'

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
            path: '/tracking',
            loader: GetUsers
          },
          {
            element: <OutputPage />,
            path: '/output'
          },
          {
            element: <SessionsListPage />,
            path: '/sessions/:userId?',
            loader: ({ params }) => GetSessions(params.userId ?? '')
          },
          {
            element: <MatchesListPage />,
            path: '/sessions/:sessionId/matches/:page?/:limit?',
            loader: ({ params }) =>
              GetMatches(
                Number(params.sessionId),
                '',
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
