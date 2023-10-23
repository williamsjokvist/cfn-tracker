import React from "react";
import { Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { CFNMachineContext } from "@/main/machine";
import { AppSidebar } from "@/main/app-layout/app-sidebar";
import { ErrorMessage } from "@/main/error-message";

import { LoadingBar } from "@/ui/loading-bar";
import { UpdatePrompt } from "@/ui/update-prompt";

import { BrowserOpenURL, EventsOn } from "@@/runtime";

export const AppWrapper: React.FC = () => {
  const { t } = useTranslation();
  const [loaded, setLoaded] = React.useState(0);
  const [hasNewVersion, setNewVersion] = React.useState(false);
  const [errorMessage, setErrorMessage] = React.useState(null);

  const [_, send] = CFNMachineContext.useActor();

  React.useEffect(() => {
    EventsOn("started-tracking", () => send("startedTracking"));
    EventsOn("stopped-tracking", () => send("stoppedTracking"));
    EventsOn("auth-loaded", (percentage) => setLoaded(percentage));
    EventsOn("version-update", (hasNewVersion) => setNewVersion(hasNewVersion));

    EventsOn("initialized", () => {
      send("initialized");
      setLoaded(100);
      setTimeout(() => setLoaded(0), 10);
    });

    EventsOn("cfn-data", (matchHistory) =>
      send({
        type: "matchPlayed",
        matchHistory,
      })
    );

    EventsOn("error-cfn", (error) => {
      send({
        type: "errorCfn",
        error,
      });
      setErrorMessage(t("cfnError"));
    });
  }, []);

  return (
    <>
      <AppSidebar />
      <main>
        <ErrorMessage
          errorMessage={errorMessage}
          onFadedOut={() => setErrorMessage(null)}
        />
        <LoadingBar percentage={loaded} />
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
