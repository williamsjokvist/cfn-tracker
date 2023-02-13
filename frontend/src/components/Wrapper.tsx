import Sidebar from "./Sidebar";
import { useEffect } from "react";
import { EventsOn, EventsOff } from "../../wailsjs/runtime"
import { useStatStore } from "../store/use-stat-store";
import { useAppStore } from "../store/use-app-store";

import { IMatchHistory } from '../types/match-history'

const Wrapper = ({ children }: any) => {
  const { setMatchHistory, setTracking, setLoading, setInitialized } = useStatStore();
  const { setNewVersionAvailable } = useAppStore();

  useEffect(() => {
    EventsOn(`cfn-data`, (mh: IMatchHistory) => {
      setMatchHistory(mh);
      setTracking(true);
      setLoading(false)
    })

    EventsOn(`initialized`, (init) => {
      setInitialized(init)
    })

    EventsOn(`version-update`, (hasNewVersion) => { 
      console.log('has new version', hasNewVersion)
      setNewVersionAvailable(hasNewVersion)
    })

    return () => {
      EventsOff(`cfn-data`)
      EventsOff(`initialized`)
      EventsOff(`version-update`)
    }
  }, [])

  return (
    <>
      <Sidebar />
      {children}
      <div className={`logo-pattern absolute filter-[grayscale(1)] bg-center`}/>
    </>
  );
};

export default Wrapper;
