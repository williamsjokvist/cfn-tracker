import Sidebar from "./Sidebar";
import { useEffect } from "react";
import { EventsOn, EventsOff } from "../../wailsjs/runtime"
import { useStatStore } from "../store/use-stat-store";
import { useAppStore } from "../store/use-app-store";

import { core } from "../../wailsjs/go/models";


const Wrapper = ({ children }: any) => {
  const { setMatchHistory, setTracking, setLoading, setInitialized, setPaused } = useStatStore();
  const { setNewVersionAvailable } = useAppStore();

  useEffect(() => {
    EventsOn(`cfn-data`, (mh: core.MatchHistory) => {
      setMatchHistory(mh);
      setTracking(true);
      setLoading(false)
    })

    EventsOn(`initialized`, (init) => {
      setInitialized(init);
      (init == true) && setLoading(false);
    })

    EventsOn(`stopped-tracking`, () => {
      setLoading(false)
    })

    EventsOn(`started-tracking`, () => {
      setPaused(false)
      setTracking(true)
    })

    EventsOn(`version-update`, (hasNewVersion) => {
      console.log('has new version', hasNewVersion)
      setNewVersionAvailable(hasNewVersion)
    })

    return () => {
      // EventsOff(`cfn-data`)
      // EventsOff(`initialized`)
      EventsOff(`version-update`)
    }
  }, [])

  return (
    <>
      <Sidebar />
      {children}
    </>
  );
};

export default Wrapper;
