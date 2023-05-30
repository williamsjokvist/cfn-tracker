import React from "react";
import { Sidebar } from "./Sidebar";

import { useAppStore } from "@/store/use-app-store";
import { CFNMachineContext } from '@/state-machine/machine';

import { core } from "@@/go/models";
import { EventsOn, EventsOff } from "@@/runtime"

export const PageWrapper: React.FC<React.PropsWithChildren> = ({ children }) => {
  const [state, send] = CFNMachineContext.useActor();
  
  const { setNewVersionAvailable } = useAppStore();

  React.useEffect(() => {
    EventsOn(`cfn-data`, (mh: core.MatchHistory) => {
      state.context.matchHistory = mh
    })

    EventsOn(`initialized`, () => send('initialized'))
    EventsOn(`started-tracking`, () => send('success'))
    EventsOn(`stopped-tracking`, () => send('stop'))
    EventsOn(`version-update`, (hasUpdated: boolean) => setNewVersionAvailable(hasUpdated))

    return () => {
      EventsOff(`version-update`)
      EventsOff(`stopped-tracking`)
      EventsOff(`started-tracking`)
      EventsOff(`initialized`)
      EventsOff(`cfn-data`)
    }
  })

  return (
    <>
      <Sidebar />
      <main className="grid grid-rows-[0fr_1fr] min-h-screen max-h-screen z-40 flex-1 text-white mx-auto">
        {children}
      </main>
    </>
  );
};