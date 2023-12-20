import React from "react";
import { Outlet } from "react-router-dom";

import { AppSidebar } from "@/main/app-layout/app-sidebar";
import { ErrorMessage } from "@/main/app-layout/error-message";

import { UpdatePrompt } from "@/main/app-layout/update-prompt";

import { LoadingBar } from "@/main/app-layout/loading-bar";

export const AppWrapper: React.FC = () => (
  <>
    <AppSidebar />
    <main>
      <ErrorMessage />
      <LoadingBar />
      <UpdatePrompt />
      <Outlet />
    </main>
  </>
);
