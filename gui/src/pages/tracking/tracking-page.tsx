import React from "react";
import { useTranslation } from "react-i18next";
import { TRACKING_MACHINE } from "@/main/machine";
import { PageHeader } from "@/ui/page-header";

import { TrackingGamePicker } from "./tracking-game-picker";
import { TrackingLiveUpdater } from "./tracking-live-updater";
import { TrackingForm } from "./tracking-form";
import { useMachine } from "@xstate/react";

export const TrackingPage: React.FC = () => {
  const { t } = useTranslation();
  const [state] = useMachine(TRACKING_MACHINE);

  switch (state.value) {
    case "formGame":
      return <TrackingGamePicker />
    case "formCfn":
      return <TrackingForm />
    case "tracking":
      return <TrackingLiveUpdater />
    case "loadingCfn":
    case "loadingGame":
      return <PageHeader text={t("loading")} showSpinner />
  }
};
