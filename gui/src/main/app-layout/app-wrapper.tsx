import React from "react";
import { Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { TrackingMachineContext } from "@/machines/tracking-machine";
import { AuthMachineContext } from "@/machines/auth-machine";

import { AppSidebar } from "@/main/app-layout/app-sidebar";
import { ErrorMessage } from "@/main/error-message";

import { LoadingBar } from "@/ui/loading-bar";
import { UpdatePrompt } from "@/ui/update-prompt";

import { BrowserOpenURL, EventsOn } from "@@/runtime";

export const AppWrapper: React.FC = () => {
  const { t } = useTranslation();
  const [loaded, setLoaded] = React.useState(0);
  const [hasNewVersion, setNewVersion] = React.useState(false);

  const trackingActor = TrackingMachineContext.useActorRef()
  const authActor = AuthMachineContext.useActorRef()

  React.useEffect(() => {
    EventsOn("stopped-tracking", () => trackingActor.send({ type: "cease" }));
    EventsOn("initialized", () => {
      authActor.send({ type: "loadedGame" })
      setTimeout(() => setLoaded(0), 10)
    });
    EventsOn("cfn-data", (trackingState) => trackingActor.send({ type: "matchPlayed", trackingState }));
    EventsOn("auth-loaded", (percentage: number) =>  setLoaded(percentage));
    EventsOn("version-update", (hasNewVersion: boolean) => setNewVersion(hasNewVersion));
  }, []);

  return (
    <>
      <AppSidebar />
      <main>
        <ErrorMessage />
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
