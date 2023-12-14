import React from "react";
import { useTranslation } from "react-i18next";
import { TrackingMachineContext } from "@/machines/tracking-machine";
import { AuthMachineContext } from "@/machines/auth-machine";
import { PageHeader } from "@/ui/page-header";

import { TrackingGamePicker } from "./tracking-game-picker";
import { TrackingLiveUpdater } from "./tracking-live-updater";
import { TrackingForm } from "./tracking-form";
import { useSelector } from "@xstate/react";

export const TrackingPage: React.FC = () => {
  const { t } = useTranslation();

  const trackingActor = TrackingMachineContext.useActorRef()
  const authActor = AuthMachineContext.useActorRef()

  const authState = useSelector(authActor, (snapshot) => snapshot.value)
  const trackingState = useSelector(trackingActor, (snapshot) => snapshot.value)

  switch (authState) {
    case "gameForm":
      return <TrackingGamePicker onSubmit={(game: string) => authActor.send({ type: "submit", game })} />
    case "loading":
      return <PageHeader text={t("loading")} showSpinner />      
  }

  switch (trackingState) {
    case "cfnForm":
      return <TrackingForm />
    case "tracking":
      return <TrackingLiveUpdater />
    case "loading":
    default:
      return <PageHeader text={t("loading")} showSpinner />
  }
};
