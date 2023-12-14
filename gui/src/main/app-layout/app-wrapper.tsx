import React from "react";
import { Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { TRACKING_MACHINE } from "@/main/machine";
import { AppSidebar } from "@/main/app-layout/app-sidebar";
import { ErrorMessage } from "@/main/error-message";

import { LoadingBar } from "@/ui/loading-bar";
import { UpdatePrompt } from "@/ui/update-prompt";

import { BrowserOpenURL, EventsOn } from "@@/runtime";
import type { errorsx } from "@@/go/models";
import { useMachine } from "@xstate/react";

export const AppWrapper: React.FC = () => {
  const { t } = useTranslation();
  const [loaded, setLoaded] = React.useState(0);
  const [hasNewVersion, setNewVersion] = React.useState(false);
  const [err, setErr] = React.useState<errorsx.FrontEndError | null>(null);

  const [_, send] = useMachine(TRACKING_MACHINE);

  React.useEffect(() => {
    EventsOn("stopped-tracking", () => send({ type: "stoppedTracking" }));
    EventsOn("auth-loaded", (percentage: number) => setLoaded(percentage));
    EventsOn("version-update", (hasNewVersion: boolean) => setNewVersion(hasNewVersion));
    EventsOn("initialized", () => {
      send({ type: "loadedGame" })
      setLoaded(100);
      setTimeout(() => setLoaded(0), 10);
    });
    EventsOn("cfn-data", (trackingState) => (
      send({
        type: "matchPlayed",
        trackingState,
      })
    ));
    EventsOn("error", (error: errorsx.FrontEndError) => {
      send({ type: "error" });
      setErr(error);
      console.error(error);
    });
  }, []);

  return (
    <>
      <AppSidebar />
      <main>
        <ErrorMessage
          error={err}
          onFadedOut={() => setErr(null)}
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
