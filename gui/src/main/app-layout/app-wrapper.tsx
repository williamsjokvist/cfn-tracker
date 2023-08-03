import React from "react";
import { Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { CFNMachineContext } from "@/main/machine";
import { BrowserOpenURL, EventsOn } from "@@/runtime";

import { AppSidebar } from "@/main/app-layout/app-sidebar";
import { LoadingBar } from "@/ui/loading-bar";
import { ErrorMessage } from "@/ui/error-message";
import { UpdatePrompt } from "@/ui/update-prompt";

export const AppWrapper: React.FC = () => {
  const { t } = useTranslation();
  const [loaded, setLoaded] = React.useState(0);
  const [errMsg, setErrMsg] = React.useState<string | undefined>();
  const [hasNewVersion, setNewVersion] = React.useState(false);

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
      setErrMsg(t("cfnError"));
    });
  }, []);

  return (
    <>
      <AppSidebar />

      <main>
        <ErrorMessage message={errMsg} />
        <LoadingBar percentage={loaded} />
        <UpdatePrompt
          onDismiss={() => {
            BrowserOpenURL("https://williamsjokvist.github.io/cfn-tracker/");
            setNewVersion(false);
          }}
          text={t("newVersionAvailable")}
          show={hasNewVersion}
        />
        <Outlet />
      </main>
    </>
  );
};
