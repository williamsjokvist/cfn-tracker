import React from "react";
import { useTranslation } from "react-i18next";
import { FaChevronLeft } from "react-icons/fa";
import { RiFileExcel2Fill } from "react-icons/ri";

import { MdDelete } from "react-icons/md";
import {
  GetAvailableLogs,
  GetMatchLog,
  DeleteMatchLog,
  ExportLogToCSV,
  OpenResultsDirectory,
} from "@@/go/core/CommandHandler";
import { common } from "@@/go/models";
import { PageHeader } from "@/ui/header";

export const HistoryPage: React.FC = () => {
  const { t } = useTranslation();

  const [availableLogs, setAvailableLogs] = React.useState<string[]>([]);
  const [chosenLog, setLog] = React.useState<string>();
  const [matchLog, setMatchLog] = React.useState<common.MatchHistory[]>();

  const [isSpecified, setSpecified] = React.useState(false);
  const [totalWinRate, setTotalWinRate] = React.useState<number | null>(null);

  const fetchLog = (log: string) => {
    GetMatchLog(log).then((log) => {
      setMatchLog([]);
      setTimeout(() => {
        setMatchLog(log);
        setSpecified(false);
      }, 50);
    })
  }

  React.useEffect(() => {
    GetAvailableLogs().then((logs) => setAvailableLogs(logs))

    if (matchLog) {
      const wonMatches = matchLog.filter((log) => log.result == true).length;
      const lostMatches = matchLog.filter((log) => log.result == false).length;
      const winRate = Math.floor((wonMatches / (wonMatches + lostMatches)) * 100);
      setTotalWinRate(winRate);
    }

    if (chosenLog && !matchLog) {
      fetchLog(chosenLog)
    }
  }, [chosenLog, matchLog]);

  const filterLog = (property: string, value: string) => {
    if (!matchLog) return;

    setMatchLog(matchLog.filter((ml) => (ml as any)[property] == value));
    setSpecified(true);
  };

  return (
    <>
      <PageHeader text={t("history")}>
        {chosenLog && (
          <div className="flex items-center justify-end w-full ml-4">
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
                ExportLogToCSV(chosenLog);
                OpenResultsDirectory();
              }}
              className="bg-[rgb(0,0,0,0.28)] hover:bg-[rgb(255,255,255,0.125)] backdrop-blur h-8 inline-block mr-3 rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
            >
              <RiFileExcel2Fill className="w-4 h-4 inline mr-2 text-white" />
              {t("exportLog")}
            </button>
            <button
              onClick={() => {
                chosenLog && DeleteMatchLog(chosenLog);
                setTimeout(() => setMatchLog([]), 50);
                setLog(undefined);
              }}
              className="bg-[rgb(0,0,0,0.28)] hover:bg-[rgb(255,255,255,0.125)] backdrop-blur h-8 inline-block float-right rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
            >
              <MdDelete className="w-4 h-4 inline mr-2" />
              {t("deleteLog")}
            </button>
          </div>
        )}
      </PageHeader>

      <div className="relative w-full">
        {chosenLog && totalWinRate != null && (
          <div className="flex items-center pt-1 px-8 mb-2 h-10 border-b border-slate-50 border-opacity-10 ">
            <span>{t("winRate")}: <b>{totalWinRate}</b>%</span>
          </div>
        )}
        {!chosenLog && availableLogs.length > 0 && (
          <div className="grid h-full justify-center content-center overflow-y-scroll">
            {availableLogs.map(cfn => {
              return (
                <button
                  className="bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-xl backdrop-blur mb-5 rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
                  key={cfn}
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
          <div className="overflow-y-scroll max-h-[340px] h-full px-4 mx-4">
            <table className="w-full border-spacing-y-1 border-separate min-w-[525px]">
              <thead>
                <tr>
                  <th className="text-left px-3 whitespace-nowrap">{t("date")}</th>
                  <th className="text-left px-3 whitespace-nowrap">{t("time")}</th>
                  <th className="text-left px-3 whitespace-nowrap">{t("opponent")}</th>
                  <th className="text-left px-3 whitespace-nowrap">{t("league")}</th>
                  <th className="text-center px-3 whitespace-nowrap">{t("character")}</th>
                  <th className="text-center px-3 whitespace-nowrap">{t("result")}</th>
                </tr>
              </thead>
              <tbody>
                {matchLog && matchLog.map((log, index) => {
                  return (
                    <tr key={index} className="backdrop-blur">
                      <td
                        onClick={() => filterLog("date", log.date)}
                        className="whitespace-nowrap text-left rounded-l-xl rounded-r-none bg-slate-50 bg-opacity-5 px-3 py-2 hover:underline cursor-pointer"
                      >
                        {log.date}
                      </td>
                      <td className="whitespace-nowrap text-left bg-slate-50 bg-opacity-5 px-3 py-2">{log.timestamp}</td>
                      <td
                        onClick={() => filterLog("opponent", log.opponent)}
                        className="whitespace-nowrap rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 hover:underline cursor-pointer"
                      >
                        {log.opponent}
                      </td>
                      <td
                        onClick={() => filterLog("opponentLeague", log.opponentLeague)}
                        className="whitespace-nowrap rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 hover:underline cursor-pointer"
                      >
                        {log.opponentLeague}
                      </td>
                      <td
                        onClick={() => filterLog("opponentCharacter", log.opponentCharacter)}
                        className="rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 text-center hover:underline cursor-pointer"
                      >
                        {log.opponentCharacter}
                      </td>
                      <td
                        className="rounded-r-xl rounded-l-none bg-slate-50 bg-opacity-5 px-3 py-2 text-center"
                        style={{ color: log.result == true ? "lime" : "red" }}
                      >
                        {log.result == true ? "W" : "L"}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </>
  );
};