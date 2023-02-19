import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { FaChevronLeft } from "react-icons/fa";
import { MdOutlineDelete } from "react-icons/md";
import {
  GetAvailableLogs,
  GetMatchLog,
  DeleteMatchLog,
} from "../../wailsjs/go/backend/App";
import { backend } from "../../wailsjs/go/models";

const History = () => {
  const { t } = useTranslation();

  const [availableLogs, setAvailableLogs] = useState<string[]>([]);
  const [chosenLog, setLog] = useState<string>();
  const [matchLog, setMatchLog] = useState<backend.MatchHistory[]>();
  const [isSpecified, setSpecified] = useState(false);
  const [totalWinRate, setTotalWinRate] = useState<number | null>(null);

  const fetchLog = async (cfn: string) => {
    const log = await GetMatchLog(cfn);
    setMatchLog([]);
    setTimeout(() => {
      setMatchLog(log);
      setSpecified(false);
    }, 50);
  };

  useEffect(() => {
    if (matchLog) {
      const wonMatches = matchLog.filter((log) => log.result == true).length;
      const lostMatches = matchLog.filter((log) => log.result == false).length;
      const winRate = Math.floor(
        (wonMatches / (wonMatches + lostMatches)) * 100
      );
      setTotalWinRate(winRate);
    }

    const fetchAvailableLogs = async () => {
      const logs = await GetAvailableLogs();
      setAvailableLogs(logs);
    };

    fetchAvailableLogs();

    chosenLog && !matchLog && fetchLog(chosenLog);
  }, [chosenLog, matchLog]);

  const filterLog = (property: string, value: string) => {
    if (!matchLog) return;

    setMatchLog([]);
    setTimeout(() => {
      setMatchLog(
        matchLog.filter((ml) => (ml as any)[property] == value)
      );
      setSpecified(true);
    }, 50);
  };
  return (
    <main className="grid grid-rows-[0fr_1fr] min-h-screen max-h-screen z-40 flex-1 text-white mx-auto">
      <header
        className="px-8 flex items-center border-b border-slate-50 border-opacity-10 backdrop-blur select-none"
        style={
          {
            "--wails-draggable": "drag",
          } as React.CSSProperties
        }
      >
        <h2 className="pt-4 whitespace-nowrap flex items-center justify-between gap-5 uppercase text-sm tracking-widest mb-4">
          {t("history")}
        </h2>
        {chosenLog && (
          <div className='flex items-center justify-between w-full ml-4'>
            <button
              onClick={() => {
                isSpecified && fetchLog(chosenLog);
                !isSpecified && setLog(undefined);
              }}
              className="bg-[rgb(0,0,0,0.28)] hover:bg-[rgb(255,255,255,0.125)] backdrop-blur h-8 inline-block mr-3 rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
            >
              <FaChevronLeft className="w-4 h-4 inline mr-2" />
              {t("goBack")}
            </button>
            <button
            onClick={() => {
              chosenLog && DeleteMatchLog(chosenLog);
              setTimeout(() => {
                setMatchLog([]);
              }, 50);
            }}
            className="bg-[rgb(0,0,0,0.28)] hover:bg-[rgb(255,255,255,0.125)] backdrop-blur h-8 inline-block float-right rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
          >
            <MdOutlineDelete className="w-4 h-4 inline mr-2" />
            {t("delete")}
          </button>
          </div>
        )}
      </header>
      <div className="relative w-full">
        {chosenLog && (
          <div className="flex items-center px-8 mb-2 h-10 border-b border-slate-50 border-opacity-10 ">
            {totalWinRate != null && (
              <span>
                {t("winRate")}: <b>{totalWinRate}</b>%
              </span>
            )}
          </div>
        )}
        {!chosenLog && availableLogs.length > 0 && (
          <div className="grid h-full justify-center content-center overflow-y-scroll">
            {availableLogs.map((cfn, index) => {
              return (
                <button
                  className="bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-xl backdrop-blur mb-5 rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
                  key={index}
                  onClick={() => {
                    setLog(cfn);
                    fetchLog(cfn);
                  }}
                >
                  {cfn}
                </button>
              );
            })}
          </div>
        )}
        {chosenLog && (
          <div className="overflow-y-scroll max-h-[340px] h-full px-8">
            <table className="w-full border-spacing-y-1 border-separate min-w-[525px]">
              <thead>
                <tr>
                  <th className="text-left px-3 whitespace-nowrap">
                    {t("date")}
                  </th>
                  <th className="text-left px-3 whitespace-nowrap">
                    {t("time")}
                  </th>
                  <th className="text-left px-3 whitespace-nowrap">
                    {t("opponent")}
                  </th>
                  <th className="text-left px-3 whitespace-nowrap">
                    {t("character")}
                  </th>
                  <th className="text-left px-3 whitespace-nowrap">
                    {t("result")}
                  </th>
                  <th className="text-left px-3 whitespace-nowrap">
                    {t("lpGain")}
                  </th>
                </tr>
              </thead>
              <tbody>
                {matchLog &&
                  matchLog.map((log, index) => {
                    return (
                      <tr key={index} className="w-full backdrop-blur">
                        <td 
                          onClick={() => filterLog('date', log.date)}
                          className="whitespace-nowrap text-center rounded-l-xl rounded-r-none bg-slate-50 bg-opacity-5 px-3 py-2 hover:underline cursor-pointer">
                          {log.date}
                        </td>
                        <td className="whitespace-nowrap text-center bg-slate-50 bg-opacity-5 px-3 py-2">
                          {log.timestamp}
                        </td>
                        <td
                          onClick={() => filterLog('opponent', log.opponent)}
                          className="whitespace-nowrap w-full rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 hover:underline cursor-pointer"
                        >
                          {log.opponent}
                        </td>
                        <td
                          onClick={() => filterLog('opponentCharacter', log.opponentCharacter)}
                          className="rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 text-center hover:underline cursor-pointer"
                        >
                          {log.opponentCharacter}
                        </td>
                        <td
                          className="rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 text-center"
                          style={{
                            color: log.result == true ? "lime" : "red",
                          }}
                        >
                          {log.result == true ? "W" : "L"}
                        </td>
                        <td className="rounded-r-xl rounded-l-none bg-slate-50 bg-opacity-5 px-3 py-2 text-center">
                          {log.lpGain > 0 && "+"}
                          {log.lpGain}
                        </td>
                      </tr>
                    );
                  })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </main>
  );
};

export default History;
