import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { GiDeerTrack } from "react-icons/gi";
import { FaStop } from "react-icons/fa";
import {
  Track,
  IsTracking,
  StopTracking,
  GetMatchHistory,
} from "../../wailsjs/go/main/App";
import { PieChart } from 'react-minimal-pie-chart';
import { useStatStore } from "../store/use-stat-store";

const Root = () => {
  const { t } = useTranslation();
  const [isLoading, setLoading] = useState(false);
  const [isCurrentlyTracking, setTracking] = useState(false);
  const {matchHistory, setMatchHistory, resetMatchHistory} = useStatStore();

  useEffect(() => {
    if (isCurrentlyTracking == true) return;

    const fetchIsTracking = async () => {
      const getIsTracking = await IsTracking();
      setTracking(getIsTracking);
    };

    fetchIsTracking();
  }, []);

  useEffect(() => {
    let interval: number | null = null;
    if (isCurrentlyTracking) {
      interval = setInterval(async () => {
        const mh = await GetMatchHistory();
        console.log(mh)
        setMatchHistory({
          cfn: mh.cfn,
          wins: mh.wins,
          losses: mh.losses,
          winRate: mh.winRate,
          lpGain: mh.lpGain,
          lp: mh.lp
        });
      }, 3000);
    }

    return () => {
      interval && clearInterval(interval);
    };
  }, [isCurrentlyTracking]);

  return (
    <main className="grid grid-rows-[0fr_1fr] min-h-screen max-h-screen z-40 flex-1 text-white mx-auto">
      <header className="border-b border-slate-50 backdrop-blur border-opacity-10 select-none " style={{
        '--wails-draggable': 'drag'
      } as React.CSSProperties}>
        <h2 className="pt-4 px-8 pl-12 flex items-center justify-between gap-5 uppercase text-sm tracking-widest mb-4">
          {isCurrentlyTracking || isLoading ? "Tracking" : t("startTracking")}
          {(isCurrentlyTracking || isLoading) && (
            <div
              className="animate-spin inline-block w-5 h-5 border-[3px] border-current border-t-transparent text-pink-600 rounded-full"
              role="status"
              aria-label="loading"
            ></div>
          )}
        </h2>
      </header>
      <div className="z-40 h-full flex justify-between items-center px-8 py-4">
        {matchHistory && isCurrentlyTracking && (
          <>
            <div className="relative grid grid-rows-[0fr_1fr]">
              <h3 className="text-3xl mr-8 pr-8 border-r border-slate-50 border-opacity-10">
                <span className="text-sm block">CFN</span>
                {matchHistory.cfn}
              </h3>
              <h4 className="text-3xl">
                <span className="text-sm block">LP</span>
                {matchHistory && matchHistory.lp && matchHistory.lp}
              </h4>
              <dl className="stat-grid-item w-full mt-4  relative text-center text-xl max-w-[265px] whitespace-nowrap">
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wider font-extralight">Wins</dt>
                  <dd className="text-5xl font-semibold">
                    {matchHistory.wins}
                  </dd>
                </div>
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">Losses</dt>
                  <dd className="text-5xl font-semibold">
                    {matchHistory.losses}
                  </dd>
                </div>
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">Win Ratio</dt>
                  <dd className="text-5xl font-semibold">
                    {matchHistory.winRate}%
                  </dd>
                </div>
                <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">LP Gain</dt>
                  <dd className="text-5xl font-semibold">
                    {matchHistory.lpGain}
                  </dd>
                </div>
              </dl>
            </div>
            {matchHistory && isCurrentlyTracking && (
              <div className='relative mr-16 h-full'>
                <PieChart
                  className='pie-chart animate-enter max-w-[160px] max-h-[160px] mt-12 backdrop-blur'
                  animate={true}
                  lineWidth={75}
                  paddingAngle={0}
                  animationDuration={10}
                  viewBoxSize={[60, 60]}
                  center={[30, 30]}
                  animationEasing={'ease-in-out'}
                  data={[
                    { title: 'Wins', value: matchHistory.wins, color: 'rgba(0, 255, 77, .65)' },
                    { title: 'Losses', value: matchHistory.losses, color: 'rgba(251, 73, 73, 0.25)' },
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

                <button
                  onClick={() => {
                    StopTracking();
                    setTracking(false);
                    setLoading(false);
                  }}
                  type="button"
                  className="backdrop-blur absolute bottom-8 right-8 flex items-center justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
                >
                  <FaStop className="mr-3" /> Stop
                </button>
              </div>
            )}
          </>
        )}
        {!(isCurrentlyTracking || isLoading) && (
          <form
            className="mx-auto"
            onSubmit={(e) => {
              e.preventDefault();
              const cfn = (e.target as any).cfn.value;
              if (cfn == "") return;
              setLoading(true);

              const x = async () => {
                const isTracking = await Track(cfn);
                setTracking(isTracking);
                if (isTracking == false) {
                  alert("Failed to track CFN");
                } else {
                  console.log("is Tracking");
                  setTracking(true);
                  resetMatchHistory()
                }
                setLoading(false);
              };
              x();
            }}
          >
            <h3 className="mb-2">Enter your CFN name</h3>
            <input
              disabled={isLoading}
              type="text"
              name="cfn"
              className="py-3 px-4 block w-full border-gray-200 rounded-md text-lg focus:border-orange-500 focus:ring-orange-500 dark:bg-gray-900 dark:border-gray-800 dark:text-gray-300"
              placeholder={t("cfnName")!}
            />
            <div className="flex justify-end">
              <button
                disabled={isLoading}
                type="submit"
                className="mt-4 flex select-none items-center justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
              >
                <GiDeerTrack className="mr-3" /> {t("start")}
              </button>
            </div>
          </form>
        )}
      </div>
    </main>
  );
};

export default Root;
