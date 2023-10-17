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

  switch (state.value) {
    case "gamePicking":
      return <TrackingGamePicker />
    case "tracking":
      return <TrackingLiveUpdater />
    case "idle":
      return <TrackingForm />
    default:
      return <PageHeader text={t("loading")} showSpinner />
  }
};
