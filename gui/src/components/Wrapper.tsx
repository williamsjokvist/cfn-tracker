import React from "react";
import { Sidebar } from "./Sidebar";

import { useStatStore } from "@/store/use-stat-store";
import { useAppStore } from "@/store/use-app-store";

import { core } from "@@/go/models";
import { EventsOn, EventsOff } from "@@/runtime"

export const PageWrapper: React.FC<React.PropsWithChildren> = ({ children }) => {
  const {
    setMatchHistory,
    setInitialized,
    setTracking,
    setLoading,
    setPaused,
  } = useStatStore();
  
  const { setNewVersionAvailable } = useAppStore();

  React.useEffect(() => {
    EventsOn(`cfn-data`, (mh: core.MatchHistory) => {
      setMatchHistory(mh);
      setTracking(true);
      setLoading(false)
    })

    EventsOn(`initialized`, (init: boolean) => {
      setInitialized(init);
      setLoading(!init)
    })

    EventsOn(`started-tracking`, () => {
      setPaused(false)
      setTracking(true)
    })

    EventsOn(`stopped-tracking`, () => setLoading(false))
    EventsOn(`version-update`, (hasNewVersion: boolean) => setNewVersionAvailable(hasNewVersion))

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