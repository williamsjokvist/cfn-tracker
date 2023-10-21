import React from "react";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";
import { Icon } from "@iconify/react";

import {
  GetAvailableLogs,
  GetMatchLog,
  DeleteMatchLog,
  ExportLogToCSV,
  OpenResultsDirectory,
} from "@@/go/core/CommandHandler";
import type { data as model } from "@@/go/models";

import { PageHeader } from "@/ui/page-header";

export const HistoryPage: React.FC = () => {
  const { t } = useTranslation();

  const [availablePlayers, setAvailablePlayers] = React.useState<model.PlayerInfo[]>([]);
  const [selectedPlayer, setLog] = React.useState<model.PlayerInfo>();
  const [matchLog, setMatchLog] = React.useState<model.MatchHistory[]>([]);

  const [isSpecified, setSpecified] = React.useState(false);
  const [totalWinRate, setTotalWinRate] = React.useState<number | null>(null);

  const fetchLog = (log: string) => {
    GetMatchLog(log).then((log) => {
      setMatchLog(log);
      setSpecified(false);
    });
  };

  React.useEffect(() => {
    GetAvailableLogs().then(
      (playersWithLogs) => setAvailablePlayers((_) => playersWithLogs)
    );
  }, []);

  React.useEffect(() => {
    if (selectedPlayer && !matchLog) fetchLog(selectedPlayer.id);
  }, [selectedPlayer, matchLog]);

  React.useEffect(() => {
    if (!matchLog) return;
    const wonMatches = matchLog.filter((log) => log.result == true).length;
    const winRate = Math.floor((wonMatches / matchLog.length) * 100);
    !isNaN(winRate) && setTotalWinRate(winRate);
  }, [matchLog]);

  const filterLog = (property: string, value: string) => {
    if (!matchLog) return;

    setMatchLog(
      matchLog.filter(
        (ml) =>
          ((ml as any)[property] as string).toLowerCase() ===
          value.toLowerCase()
      )
    );
    setSpecified(true);
  };

  return (
    <>
      <PageHeader
        text={selectedPlayer ? `${t("history")}/${selectedPlayer}` : t("history")}
      >
        {selectedPlayer && (
          <div className="flex items-center justify-end w-full ml-4">
            <motion.button
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              onClick={() =>
                isSpecified ? fetchLog(selectedPlayer.id) : setLog(undefined)
              }
              className="crumb-btn-dark mr-3"
            >
              <Icon
                icon="fa6-solid:chevron-left"
                className="w-4 h-4 inline mr-2"
              />
              {t("goBack")}
            </motion.button>
            <motion.button
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              onClick={() => {
                ExportLogToCSV(selectedPlayer.id);
                OpenResultsDirectory();
              }}
              className="crumb-btn-dark mr-3"
            >
              <Icon
                icon="ri:file-excel-2-fill"
                className="w-4 h-4 inline mr-2 text-white"
              />
              {t("exportLog")}
            </motion.button>
            <motion.button
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              onClick={() => {
                selectedPlayer && DeleteMatchLog(selectedPlayer.id);
                setTimeout(() => setMatchLog([]), 50);
                setLog(undefined);
              }}
              className="crumb-btn-dark"
            >
              <Icon icon="mdi:delete" className="w-4 h-4 inline mr-2" />
              {t("deleteLog")}
            </motion.button>
          </div>
        )}
      </PageHeader>

      <div className="relative w-full">
        {!selectedPlayer && availablePlayers && availablePlayers.length > 0 && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.225 }}
            className="grid max-h-[340px] h-full m-4 justify-center content-start overflow-y-scroll gap-5"
          >
            {availablePlayers.map((playerInfo) => (
              <button
                className="bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-xl backdrop-blur rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
                key={playerInfo.displayName}
                onClick={() => {
                  setLog(playerInfo);
                  fetchLog(playerInfo.id);
                }}
              >
                {playerInfo.displayName}
              </button>
            ))}
          </motion.div>
        )}
        {matchLog && selectedPlayer && (
          <>
            {totalWinRate != null && (
              <div className="flex items-center pt-1 px-8 mb-2 h-10 border-b border-slate-50 border-opacity-10 ">
                <span>
                  {t("winRate")}: <b>{totalWinRate}</b>%
                </span>
              </div>
            )}
            <div className="overflow-y-scroll max-h-[340px] h-full mx-4 px-4 pb-4">
              <motion.table
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.25 }}
                className="w-full border-spacing-y-1 border-separate min-w-[525px]"
              >
                <thead>
                  <tr>
                    <th className="text-left px-3 whitespace-nowrap w-[120px]">
                      {t("date")}
                    </th>
                    <th className="text-left px-3 whitespace-nowrap w-[70px]">
                      {t("time")}
                    </th>
                    <th className="text-left px-3 whitespace-nowrap w-[180px]">
                      {t("opponent")}
                    </th>
                    <th className="text-left px-3 whitespace-nowrap">
                      {t("league")}
                    </th>
                    <th className="text-center px-3 whitespace-nowrap">
                      {t("character")}
                    </th>
                    <th className="text-center px-3 whitespace-nowrap">
                      {t("result")}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {matchLog.map((log) => (
                    <tr
                      key={`${log.timestamp}-${log.totalMatches}`}
                      className="backdrop-blur group"
                    >
                      <td
                        onClick={() => filterLog("date", log.date)}
                        className="whitespace-nowrap text-left rounded-l-xl rounded-r-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2 hover:underline cursor-pointer"
                      >
                        {log.date}
                      </td>
                      <td className="whitespace-nowrap text-left bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                        {log.timestamp}
                      </td>
                      <td
                        onClick={() => filterLog("opponent", log.opponent)}
                        className="whitespace-nowrap rounded-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2 hover:underline cursor-pointer"
                      >
                        {log.opponent}
                      </td>
                      <td
                        onClick={() =>
                          filterLog("opponentLeague", log.opponentLeague)
                        }
                        className="whitespace-nowrap rounded-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2 hover:underline cursor-pointer"
                      >
                        {log.opponentLeague}
                      </td>
                      <td
                        onClick={() =>
                          filterLog("opponentCharacter", log.opponentCharacter)
                        }
                        className="rounded-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2 text-center hover:underline cursor-pointer"
                      >
                        {log.opponentCharacter}
                      </td>
                      <td
                        className="rounded-r-xl rounded-l-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2 text-center"
                        style={{ color: log.result == true ? "lime" : "red" }}
                      >
                        {log.result == true ? "W" : "L"}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </motion.table>
            </div>
          </>
        )}
      </div>
    </>
  );
};
