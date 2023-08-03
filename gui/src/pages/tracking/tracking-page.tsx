import React from "react";
import { useTranslation } from "react-i18next";
import { CFNMachineContext } from "@/main/machine";
import { PageHeader } from "@/ui/page-header";

import { TrackingGamePicker } from "./tracking-game-picker";
import { TrackingLiveUpdater } from "./tracking-live-updater";
import { TrackingForm } from "./tracking-form";

export const TrackingPage: React.FC = () => {
  const { t } = useTranslation();
  const [state] = CFNMachineContext.useActor();

  return (
    <>
      <PageHeader
        {...(state.matches("tracking") && { text: t("tracking") })}
        {...(state.matches("idle") && { text: t("startTracking") })}
        {...(state.matches("gamePicking") && { text: t("pickGame") })}
        {...((state.matches("loading") || state.matches("loadingCfn")) && {
          text: t("loading"),
        })}
        {...(!(state.matches("idle") || state.matches("gamePicking")) && {
          showSpinner: true,
        })}
      />
      <div className="z-40 h-full w-full justify-self-center flex justify-between items-center px-8 py-4">
        {state.matches("gamePicking") && <TrackingGamePicker />}
        {state.matches("tracking") && <TrackingLiveUpdater />}
        {state.matches("idle") && <TrackingForm />}
      </div>
    </>
  );
};
