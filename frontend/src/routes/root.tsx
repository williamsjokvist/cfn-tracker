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

const Root = () => {
  const { t } = useTranslation();
  const [isLoading, setLoading] = useState(false);
  const [isCurrentlyTracking, setTracking] = useState(false);
  const [cfn, setCFN] = useState('');
  const [matchHistory, setMatchHistory] = useState<{
    wins: number;
    losses: number;
    winRate: any;
    lpGain: any;
    lp: number;
  } | null>(null);

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
      console.log("gooo");
      interval = setInterval(async () => {
        const mh = await GetMatchHistory();
        console.log(mh)
        setMatchHistory({
          wins: mh.wins,
          losses: mh.losses,
          winRate: mh.winRate,
          lpGain: mh.lpGain,
          lp: mh.lp
        });
        setCFN(mh.cfn)
      }, 3000);
    }

    return () => {
      interval && clearInterval(interval);
    };
  }, [isCurrentlyTracking]);

  return (
    <main className="grid grid-rows-[0fr_1fr] min-h-screen max-h-screen z-40 flex-1 text-white mx-auto">
      <header className="border-b border-slate-50 backdrop-blur border-opacity-10 --wails-draggable select-none " style={{
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
      <div className="z-40 h-full flex justify-between px-12 py-4">
        {isCurrentlyTracking && (
          <>
            <div className="relative">
              <h3 className="text-3xl">
                <span className="text-sm block">CFN</span>
                {cfn && cfn}
              </h3>
              <h4 className="text-3xl">
                <span className="text-sm block">LP</span>
                {matchHistory && matchHistory.lp && matchHistory.lp}
              </h4>
              <button
                onClick={() => {
                  StopTracking();
                  setTracking(false);
                  setLoading(false);
                }}
                type="button"
                className="absolute bottom-8 left-8 flex items-center justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
              >
                <FaStop className="mr-3" /> Stop
              </button>
            </div>
            {matchHistory && (
              <dl className="relative text-center grid gap-2 text-xl max-w-[250px] whitespace-nowrap">
                <div className="flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 rounded-xl backdrop-blur">
                  <dt className="tracking-wider font-extralight">Wins</dt>
                  <dd className="text-5xl font-semibold">
                    {matchHistory.wins}
                  </dd>
                </div>
                <div className="flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">Losses</dt>
                  <dd className="text-5xl font-semibold">
                    {matchHistory.losses}
                  </dd>
                </div>
                <div className="flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">Win Ratio</dt>
                  <dd className="text-5xl font-semibold">
                    {matchHistory.winRate}
                  </dd>
                </div>
                <div className="flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 rounded-xl backdrop-blur">
                  <dt className="tracking-wide font-extralight">LP Gain</dt>
                  <dd className="text-5xl font-semibold">
                    {matchHistory.lpGain}
                  </dd>
                </div>
              </dl>
            )}
          </>
        )}
        {!(isCurrentlyTracking || isLoading) && (
          <form
            className="mt-10 mx-auto"
            onSubmit={(e) => {
              e.preventDefault();
              const cfn = (e.target as any).cfn.value;
              if (cfn == "") return;
              setCFN(cfn)
              setLoading(true);

              const x = async () => {
                const isTracking = await Track(cfn);
                setTracking(isTracking);
                if (isTracking == false) {
                  alert("Failed to track CFN");
                } else {
                  console.log("is Tracking");
                  setTracking(true);
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
              className="py-3 px-4 block w-full border-gray-200 rounded-md text-lg focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
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
