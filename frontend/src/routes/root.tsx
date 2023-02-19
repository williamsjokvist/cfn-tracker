import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { GiDeerTrack } from "react-icons/gi";
import { FaStop, FaPause, FaPlay } from "react-icons/fa";
import {
  StopTracking,
  OpenResultsDirectory,
  IsTracking,
  IsInitialized,
  GetAvailableLogs,
  StartTracking,
} from "../../wailsjs/go/backend/App";
import { PieChart } from "react-minimal-pie-chart";
import { AiFillFolderOpen } from "react-icons/ai";

import { useStatStore } from "../store/use-stat-store";

// TODO refactor this mess

const Root = () => {
  const { t } = useTranslation();
  const {
    matchHistory,
    resetMatchHistory,
    setTracking,
    isTracking,
    isLoading,
    setLoading,
    setInitialized,
    isInitialized,
    isPaused,
    setPaused,
  } = useStatStore();
  const [oldCfns, setOldCfns] = useState<string[] | null>(null);
  const [inputValue, setInputValue] = useState<string | null>(null);
  useEffect(() => {
    if (!isTracking) {
      GetAvailableLogs().then((logs) => {
        setOldCfns(logs);
      });

      IsTracking().then((isTracking) => {
        setTracking(isTracking);
        isTracking && setPaused(false);
      });
    }

    if (!isInitialized) {
      IsInitialized().then((isInitialized) => {
        setInitialized(isInitialized);
      });
    }
  }, []);

  return (
    <main className="grid grid-rows-[0fr_1fr] min-h-screen max-h-screen z-40 flex-1 text-white mx-auto">
      <header
        className="border-b border-slate-50 backdrop-blur border-opacity-10 select-none "
        style={
          {
            "--wails-draggable": "drag",
          } as React.CSSProperties
        }
      >
        <h2 className="pt-4 px-8 flex items-center justify-between gap-5 uppercase text-sm tracking-widest mb-4">
          {isTracking && !isLoading && !isPaused && t("tracking")}
          {(isLoading || !isInitialized) && t("loading")}
          {isPaused && !isLoading && (
            <>
              {t("pause")}
              <FaPause className="w-5 h-5 text-pink-600" />
            </>
          )}
          {!isTracking && isInitialized && !isLoading && t("startTracking")}
          {((isTracking || isLoading || !isInitialized) && !isPaused) && (
            <div
              className="animate-spin inline-block w-5 h-5 border-[3px] border-current border-t-transparent text-pink-600 rounded-full"
              role="status"
              aria-label="loading"
            ></div>
          )}
        </h2>
      </header>
      <div className="z-40 h-full flex justify-between items-center px-8 py-4">
        {matchHistory && isTracking && (
          <>
            <div className="relative w-full h-full grid grid-rows-[0fr_1fr] max-w-[320px]">
              <h3 className="whitespace-nowrap max-w-[145px] text-2xl">
                <span className="text-sm block">CFN</span>
                <span className="text-ellipsis block overflow-hidden">
                  {matchHistory.cfn}
                </span>
              </h3>
              <h4 className="text-2xl">
                <span className="text-sm block">LP</span>
                {matchHistory && matchHistory.lp && matchHistory.lp}
              </h4>
              <dl className="stat-grid-item w-full mt-2 relative text-center text-lg whitespace-nowrap">
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wider font-extralight">
                    {t("wins")}
                  </dt>
                  <dd className="text-4xl font-semibold">
                    {matchHistory.wins}
                  </dd>
                </div>
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">
                    {t("losses")}
                  </dt>
                  <dd className="text-4xl font-semibold">
                    {matchHistory.losses}
                  </dd>
                </div>
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">
                    {t("winRate")}
                  </dt>
                  <dd className="text-4xl font-semibold">
                    {matchHistory.winRate}%
                  </dd>
                </div>
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">
                    {t("winStreak")}
                  </dt>
                  <dd className="text-4xl font-semibold">
                    {matchHistory.winStreak}
                  </dd>
                </div>
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">
                    {t("lpGain")}
                  </dt>
                  <dd className="text-4xl font-semibold">
                    {matchHistory.lpGain > 0 && "+"}
                    {matchHistory.lpGain}
                  </dd>
                </div>
              </dl>
            </div>
            {matchHistory && isTracking && (
              <div className="relative mr-4 h-full grid content-between justify-items-center">
                <PieChart
                  className="pie-chart animate-enter max-w-[160px] max-h-[160px] backdrop-blur"
                  animate={true}
                  lineWidth={75}
                  paddingAngle={0}
                  animationDuration={10}
                  viewBoxSize={[60, 60]}
                  center={[30, 30]}
                  animationEasing={"ease-in-out"}
                  data={[
                    {
                      title: "Wins",
                      value: matchHistory.wins,
                      color: "rgba(0, 255, 77, .65)",
                    },
                    {
                      title: "Losses",
                      value: matchHistory.losses,
                      color: "rgba(251, 73, 73, 0.25)",
                    },
                  ]}
                >
                  <defs>
                    <linearGradient id="blue-gradient" direction={-65}>
                      <stop offset="0%" stopColor="#20BF55" />
                      <stop offset="100%" stopColor="#347fd0" />
                    </linearGradient>
                    <linearGradient id="red-gradient" direction={120}>
                      <stop offset="0%" stopColor="#EC9F05" />
                      <stop offset="100%" stopColor="#EE9617" />
                    </linearGradient>
                  </defs>
                </PieChart>

                <div className="relative bottom-[10px] grid justify-items-end right-[-10px]">
                  <button
                    disabled={isLoading}
                    onClick={() => {
                      StopTracking();
                      setLoading(true);
                      setTimeout(() => {
                        setLoading(false);
                        setPaused(false);
                        setTracking(false);
                      }, 30000)

                      setInitialized(false);
                    }}
                    type="button"
                    className="flex items-center mb-2 justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
                  >
                    <FaStop className="mr-3" /> {t("stop")}
                  </button>
                  <button
                    onClick={() => {
                      OpenResultsDirectory();
                    }}
                    style={{
                      filter: "hue-rotate(-120deg)",
                    }}
                    type="button"
                    className="whitespace-nowrap flex items-center justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
                  >
                    <AiFillFolderOpen className="w-4 h-4 mr-2" />
                    {t("openResultFolder")}
                  </button>

                </div>
              </div>
            )}
          </>
        )}
        {!(isTracking || isLoading || !isInitialized) && (
          <>
            <form
              className="mx-auto"
              onSubmit={(e) => {
                e.preventDefault();

                const cfn = (e.target as any).cfn.value;
                if (cfn == "") return;
                setLoading(true);
                const startTrack = async () => {
                  const isInitialized = await IsInitialized();
                  if (isInitialized == false) {
                    setLoading(false);
                    return;
                  }

                  StartTracking(cfn, true);
                };
                startTrack();
              }}
            >
              <h3 className="mb-2">{t("enterCfnName")}:</h3>
              <input
                disabled={isLoading}
                type="text"
                name="cfn"
                {...(inputValue && {
                  value: inputValue,
                })}
                onChange={(e) => {
                  setInputValue(e.target.value);
                }}
                className="bg-transparent border-b-2 border-0 focus-within:outline-none focus-visible:outline-none border-b-[rgba(255,255,255,0.275)] outline-none focus:outline-none hover:text-white transition-colors py-3 px-4 block w-full text-lg text-gray-300"
                placeholder={t("cfnName")!}
                autoCapitalize="off"
                autoComplete="off"
                autoCorrect="off"
                autoSave="off"
                onClick={() => {
                  setInputValue(null);
                }}
              />
              <div className={`flex items-center mt-4`}>
                <input type="checkbox"/>
                {`Restore previous session`}
              </div>
              <div className="flex justify-end">
                <button
                  disabled={isLoading}
                  type="submit"
                  style={{
                    filter: "hue-rotate(156deg)",
                  }}
                  className="mt-4 flex select-none items-center justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
                >
                  <GiDeerTrack className="mr-3" /> {t("start")}
                </button>
              </div>
            </form>
            {oldCfns && (
              <div className="border-l-[1px] border-[rgba(255,255,255,0.25)] pl-12 grid content-center items-center text-center h-[80vh] pr-3">
                {oldCfns.map((cfn, index) => {
                  return (
                    <button
                      disabled={isLoading}
                      onClick={() => setInputValue(cfn)}
                      className="bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-xl backdrop-blur mb-5 rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
                      type="button"
                      key={index}
                    >
                      {cfn}
                    </button>
                  );
                })}
              </div>
            )}
          </>
        )}
      </div>
    </main>
  );
};

export default Root;
