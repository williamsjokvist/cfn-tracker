import React from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider, createHashRouter } from "react-router-dom";

import { initLocalization } from "@/i18n/i18n-config";

import { TrackingPage } from "@/pages/tracking/tracking-page";
import { OutputPage } from "@/pages/output-page";
import { AppWrapper } from "./app-layout/app-wrapper";
import { CFNMachineContext } from "./machine";

import "@/styles/globals.sass";
import { MatchesListPage } from "@/pages/stats/matches-list-page";
import { SessionsListPage } from "@/pages/stats/sessions-list-page";

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
      <CFNMachineContext.Provider>
        <RouterProvider router={router} />
      </CFNMachineContext.Provider>
    </React.StrictMode>
  );
};

createRoot(document.getElementById("root")!).render(<App />);
