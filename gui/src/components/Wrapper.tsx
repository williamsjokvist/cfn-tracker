import React from "react";
import { Outlet } from "react-router-dom";
import { Sidebar } from "./Sidebar";

import { useAppStore } from "@/store/use-app-store";
import { CFNMachineContext } from '@/machine';

import { EventsOn, EventsOff } from "@@/runtime"

export const PageWrapper: React.FC = () => {
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