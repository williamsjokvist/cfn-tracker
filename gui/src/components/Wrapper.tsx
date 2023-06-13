import React from "react";
import { Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Sidebar } from "./Sidebar";

import { useAppStore } from "@/store/use-app-store";
import { CFNMachineContext } from '@/machine';

import { EventsOn, EventsOff } from "@@/runtime"
import { APP_LANGUAGES } from "@/i18n";

export const PageWrapper: React.FC = () => {
  const { i18n } = useTranslation();
  const [_, send] = CFNMachineContext.useActor();
  const { setNewVersionAvailable } = useAppStore();

  React.useEffect(() => {
    EventsOn(`version-update`, updated => setNewVersionAvailable(updated))
    EventsOn(`initialized`, () => send('initialized'))
    EventsOn(`started-tracking`, () => send('startedTracking'))
    EventsOn(`stopped-tracking`, () => send('stoppedTracking'))
    EventsOn(`cfn-data`, matchHistory => send({
      type: 'matchPlayed',
      matchHistory
    }))

    const storedLang = window.localStorage.getItem('lng')
    if (APP_LANGUAGES.map(lng => lng.code).includes(storedLang))
      i18n.changeLanguage(storedLang)
  })

  return (
    <>
      <Sidebar />
      <main>
        <Outlet />
      </main>
    </>
  );
};