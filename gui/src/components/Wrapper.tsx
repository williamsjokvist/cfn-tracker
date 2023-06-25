import React from "react";
import { Outlet } from "react-router-dom";

import { useAppStore } from "@/store/use-app-store";
import { CFNMachineContext } from '@/machine';
import { EventsOn, EventsOff } from "@@/runtime"
import { LoadingBar } from "@/ui/loading-bar";
import { Sidebar } from "./Sidebar";
import { ErrorMessage } from "./ErrorMessage";

export const PageWrapper: React.FC = () => {
  const [_, send] = CFNMachineContext.useActor();
  const { setNewVersionAvailable } = useAppStore();
  const [loaded, setLoaded] = React.useState(0)
  const [errMsg, setErrMsg] = React.useState('')

  React.useEffect(() => {
    EventsOn(`version-update`, updated => setNewVersionAvailable(updated))
    EventsOn(`initialized`, () => {
      send('initialized')
      setLoaded(100)
      setTimeout(() => setLoaded(0),10)
    })
    EventsOn(`started-tracking`, () => send('startedTracking'))
    EventsOn(`stopped-tracking`, () => send('stoppedTracking'))
    EventsOn(`auth-loaded`, percentage => setLoaded(percentage))
    EventsOn(`cfn-data`, matchHistory => send({
      type: 'matchPlayed',
      matchHistory
    }))

    EventsOn(`error-cfn`, error => {
      send({
        type: 'errorCfn',
        error
      })
      setErrMsg(error)
    })
  })

  return (
    <>
      <Sidebar />
      <main>
        <ErrorMessage message={errMsg}/>
        <LoadingBar percentage={loaded}/>
        <Outlet />
      </main>
    </>
  );
};