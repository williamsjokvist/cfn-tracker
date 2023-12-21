import React from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createHashRouter } from "react-router-dom";

import { TrackingPage } from "@/pages/tracking/tracking-page";
import { OutputPage } from "@/pages/output/output-page";
import { MatchesListPage } from "@/pages/stats/matches-list-page";
import { SessionsListPage } from "@/pages/stats/sessions-list-page";

import { TrackingMachineContext } from "@/machines/tracking-machine";
import { AuthMachineContext } from "@/machines/auth-machine";

import { GetMatches, GetSessions, GetThemeList, GetUsers } from "@@/go/core/CommandHandler";

import { AppLoader } from "./app-layout/app-loader";
import { AppWrapper } from "./app-layout/app-wrapper";
import { AppErrorBoundary } from "./app-layout/app-error";
import { ErrorMessageProvider } from "./app-layout/error-message";
import { ConfigProvider } from "./config";

import { initI18n } from "./i18n";

import "./style.sass";

const router = createHashRouter([
  {
    element: <AppWrapper />,
    errorElement: <AppErrorBoundary outer />,
    children: [
      {
        errorElement: <AppErrorBoundary />,
        children: [
          {
            element: <TrackingPage />,
            path: "/",
            loader: GetUsers,
          },
          {
            element: <OutputPage />,
            path: "/output",
            loader: GetThemeList
          },
          {
            element: <SessionsListPage />,
            path: "/sessions/:userId?",
            loader: ({ params }) => GetSessions(params.userId ?? "")
          },
          {
            element: <MatchesListPage />,
            path: "/sessions/:sessionId/matches/:page?/:limit?",
            loader: ({ params }) => GetMatches(Number(params.sessionId), "", Number(params.page ?? 0), Number(params.limit ?? 0))
          }
        ]
      },
    ],
  },
]);

const App: React.FC = () => {
  initI18n();
  return (
    <AuthMachineContext.Provider>
      <TrackingMachineContext.Provider>
        <React.Suspense fallback={<AppLoader/>}>
          <ErrorMessageProvider>
            <ConfigProvider>
              <RouterProvider router={router} />
            </ConfigProvider>
          </ErrorMessageProvider>
        </React.Suspense>
      </TrackingMachineContext.Provider>
    </AuthMachineContext.Provider>
  );
}

createRoot(document.getElementById("root")!).render(<App />);
