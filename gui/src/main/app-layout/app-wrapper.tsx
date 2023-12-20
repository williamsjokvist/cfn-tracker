import React from "react";
import { Outlet } from "react-router-dom";

import { AppSidebar } from "@/main/app-layout/app-sidebar";
import { ErrorMessage } from "@/main/error-message";

import { LoadingBar } from "@/main/loading-bar";

export const AppWrapper: React.FC = () => {
  return (
    <>
      <AppSidebar />
      <main>
        <ErrorMessage />
        <LoadingBar />
        <Outlet />
      </main>
    </>
  );
};
