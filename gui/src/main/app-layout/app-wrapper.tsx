import React from "react";
import { Outlet } from "react-router-dom";

import { AppSidebar } from "@/main/app-layout/app-sidebar";
import { ErrorMessageProvider } from "@/main/app-layout/error-message";
import { UpdatePrompt } from "@/main/app-layout/update-prompt";
import { LoadingBar } from "@/main/app-layout/loading-bar";

export const AppWrapper: React.FC = () => (
  <>
    <AppSidebar />
    <main className="relative grid grid-rows-[0fr_1fr] h-screen z-40 flex-[1] text-white mx-auto">
      <LoadingBar />
      <ErrorMessageProvider>
        <UpdatePrompt />
        <React.StrictMode>
          <Outlet />
        </React.StrictMode>
      </ErrorMessageProvider>
    </main>
  </>
);
