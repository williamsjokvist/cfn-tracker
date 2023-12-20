import React from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createHashRouter } from "react-router-dom";
import { initLocalization } from "@/main/i18n-config";

import { TrackingPage } from "@/pages/tracking/tracking-page";
import { OutputPage } from "@/pages/output/output-page";
import { AppWrapper } from "./app-layout/app-wrapper";

import { MatchesListPage } from "@/pages/stats/matches-list-page";
import { SessionsListPage } from "@/pages/stats/sessions-list-page";
import { TrackingMachineContext } from "@/machines/tracking-machine";
import { AuthMachineContext } from "@/machines/auth-machine";

import "./style.sass";

const router = createHashRouter([
  {
    element: <AppWrapper />,
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
  initLocalization();
  return (
    <React.StrictMode>
      <AuthMachineContext.Provider>
        <TrackingMachineContext.Provider>
          <RouterProvider router={router} />
        </TrackingMachineContext.Provider>
      </AuthMachineContext.Provider>
    </React.StrictMode>
  );
};

createRoot(document.getElementById("root")!).render(<App />);
