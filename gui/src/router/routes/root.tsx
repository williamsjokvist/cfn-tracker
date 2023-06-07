import React from "react";
import { useTranslation } from "react-i18next";
import { CFNMachineContext } from '@/machine'
import { PageHeader } from "@/ui/header";

import { GamePicker } from "./root/game-picker";
import { Tracking } from "./root/tracking";
import { CFNForm } from "./root/cfn-form";

export const RootPage: React.FC = () => {
  const { t } = useTranslation();
  const [state] = CFNMachineContext.useActor();

  return (
    <>
      <PageHeader 
        {...(state.matches('tracking') && { text: t("tracking") })}
        {...(state.matches('idle') && { text: t("startTracking") })}
        {...(state.matches('gamePicking') && { text: t("pickGame") })}
        {...((state.matches('loading') || state.matches('loadingCfn')) && { 
          text: t("loading"), 
        })}
        {...(!((state.matches('idle') || state.matches('gamePicking'))) && { showSpinner: true })}
      />
      <div className="z-40 h-full w-full justify-self-center flex justify-between items-center px-8 py-4">
        {state.matches('gamePicking') && <GamePicker/>}
        {state.matches('tracking') && <Tracking/>}
        {state.matches('idle') && <CFNForm/>}
      </div>
    </>
  );
};

