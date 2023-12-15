import React from "react";
import { Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { AppSidebar } from "@/main/app-layout/app-sidebar";
import { ErrorMessage } from "@/main/error-message";

import { UpdatePrompt } from "@/ui/update-prompt";

import { BrowserOpenURL, EventsOn } from "@@/runtime";
import { LoadingBar } from "@/main/loading-bar";

export const AppWrapper: React.FC = () => {
  const { t } = useTranslation();
  const [hasNewVersion, setNewVersion] = React.useState(false);

  React.useEffect(() => {
    EventsOn("version-update", (hasNewVersion: boolean) => setNewVersion(hasNewVersion));
  }, []);

  return (
    <>
      <AppSidebar />
      <main>
        <ErrorMessage />
        <LoadingBar />
        {hasNewVersion && (
          <UpdatePrompt
            onDismiss={() => {
              BrowserOpenURL("https://cfn.williamsjokvist.se/");
              setNewVersion(false);
            }}
            text={t("newVersionAvailable")}
          />
        )}
        <Outlet />
      </main>
    </>
  );
};
