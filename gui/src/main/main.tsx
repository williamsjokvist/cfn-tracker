import React from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createHashRouter } from "react-router-dom";

import { TrackingPage } from "@/pages/tracking/tracking-page";
import { OutputPage } from "@/pages/output/output-page";
import { MatchesListPage } from "@/pages/stats/matches-list-page";
import { SessionsListPage } from "@/pages/stats/sessions-list-page";

import { TrackingMachineContext } from "@/machines/tracking-machine";
import { AuthMachineContext } from "@/machines/auth-machine";

import { AppLoader } from "./app-layout/app-loader";
import { AppWrapper } from "./app-layout/app-wrapper";
import { AppErrorBoundary } from "./app-layout/app-error";
import { initI18n } from "./i18n";

import "./style.sass";

const router = createHashRouter([
  {
    element: <AppWrapper />,
    errorElement: <AppErrorBoundary />,
    children: [
      {
        element: <TrackingPage />,
        path: "/",
      },
      {
        element: <OutputPage />,
        path: "/output",
      },
      {
        element: <SessionsListPage />,
        path: "/sessions"
      },
      {
        element: <MatchesListPage />,
        path: "/sessions/:sessionId"
      }
    ],
  },
]);

const App: React.FC = () => {
  initI18n();
  return (
    <AuthMachineContext.Provider>
      <TrackingMachineContext.Provider>
        <React.Suspense fallback={<AppLoader/>}>
          <RouterProvider router={router} />
        </React.Suspense>
      </TrackingMachineContext.Provider>
    </AuthMachineContext.Provider>
  );
}

createRoot(document.getElementById("root")!).render(<App />);
