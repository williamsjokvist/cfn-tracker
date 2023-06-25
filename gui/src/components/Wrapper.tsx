import React from "react";
import { Outlet } from "react-router-dom";
import { useTranslation } from "react-i18next";

import { CFNMachineContext } from '@/machine';
import { EventsOn, EventsOff } from "@@/runtime"
import { LoadingBar } from "@/ui/loading-bar";
import { Sidebar } from "./Sidebar";
import { ErrorMessage } from "./ErrorMessage";
import { NewVersionPrompt } from "@/ui/new-version";

export const PageWrapper: React.FC = () => {
  const { t } = useTranslation()
  const [_, send] = CFNMachineContext.useActor();
  const [loaded, setLoaded] = React.useState(0)
  const [errMsg, setErrMsg] = React.useState('')
  const [hasNewVersion, setNewVersion] = React.useState(false)

  React.useEffect(() => {
    EventsOn(`version-update`, hasNewVersion => setNewVersion(hasNewVersion))
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
      setErrMsg(t('cfnError'))
    })
  })

  return (
    <>
      <Sidebar>
        <NewVersionPrompt hasNewVersion={hasNewVersion} setNewVersion={setNewVersion}/>  
      </Sidebar>
      <main>
        <ErrorMessage message={errMsg}/>
        <LoadingBar percentage={loaded}/>
        <Outlet />
      </main>
    </>
  );
};