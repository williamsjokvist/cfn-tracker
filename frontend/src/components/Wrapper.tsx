import Sidebar from "./Sidebar";
import { useEffect } from "react";
import { EventsOn, EventsOff } from "../../wailsjs/runtime"
import { useStatStore } from "../store/use-stat-store";
import { IMatchHistory } from '../types/match-history'

const Wrapper = ({ children }: any) => {
  const { setMatchHistory, setTracking, setLoading, setInitialized } = useStatStore();

  useEffect(() => {
    EventsOn(`cfn-data`, (mh: IMatchHistory) => {
      console.log(mh)
      setMatchHistory(mh);
      setTracking(true);
      setLoading(false)
    })

    EventsOn(`initialized`, (init) => {
      setInitialized(init)
    })

    return () => {
      EventsOff(`cfn-data`)
      EventsOff(`initialized`)
    }
  }, [])

  return (
    <>
      <Sidebar />
      {children}
      <div className='logo-pattern absolute filter-[grayscale(1)] bg-[url(src/assets/logo.png)] bg-center'/>
    </>
  );
};

export default Wrapper;
