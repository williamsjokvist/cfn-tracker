import Footer from "./Footer";
import Sidebar from "./Sidebar";
import { useState, useEffect } from "react";
import { EventsOn, EventsOff } from "../../wailsjs/runtime"
import { useStatStore } from "../store/use-stat-store";
import { IMatchHistory } from '../types/match-history'

const Wrapper = ({ children }: any) => {
  const { setMatchHistory, setTracking, setLoading } = useStatStore();

  useEffect(() => {
    /*
    if (!isTracking) {
      const getIsTracking = async () => {
        const trackStatus = await IsTracking()
        setTracking(trackStatus)
      }
      getIsTracking()
    }*/

    EventsOn(`cfn-data`, (mh: IMatchHistory) => {
      console.log(mh)
      setMatchHistory(mh);
      setTracking(true);
      setLoading(false)
    })

    return () => {
      EventsOff(`cfn-data`)
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
