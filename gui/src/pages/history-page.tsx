import React from "react";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";
import { Icon } from "@iconify/react";

import {
  GetUsers,
  GetAllMatchesForUser,
  DeleteMatchLog,
  ExportLogToCSV,
  OpenResultsDirectory,
} from "@@/go/core/CommandHandler";
import type { model } from "@@/go/models";

import { PageHeader } from "@/ui/page-header";

export const HistoryPage: React.FC = () => {
  const { t } = useTranslation();

  const [users, setUsers] = React.useState<model.User[]>([]);
  
  const [user, setUser] = React.useState<model.User>();
  const [matches, setMatches] = React.useState<model.Match[]>([]);

  const [isSpecified, setSpecified] = React.useState(false);
  const [totalWinRate, setTotalWinRate] = React.useState<number | null>(null);

  const getMatches = (u: string) => {
    GetAllMatchesForUser(u).then((matches) => {
      setMatches(matches);
      setSpecified(false);
    }).catch(err => console.error(err));
  };

  React.useEffect(() => {
    GetUsers().then(
      (users) => setUsers((_) => users)
    ).catch(err => console.error(err));
  }, []);

  React.useEffect(() => {
    if (user && !matches) getMatches(user.code);
  }, [user, matches]);

  React.useEffect(() => {
    if (!matches) return;
    const wonMatches = matches.filter((log) => log.victory === true).length;
    const winRate = Math.floor((wonMatches / matches.length) * 100);
    !isNaN(winRate) && setTotalWinRate(winRate);
  }, [matches]);

  const filterLog = (property: string, value: string) => {
    if (!matches) return;

    setMatches(
      matches.filter(
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
        text={user ? `${t("history")}/${user.displayName}` : t("history")}
      >
        {user && (
          <div className="flex items-center justify-end w-full ml-4">
            <motion.button
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              onClick={() =>
                isSpecified ? getMatches(user.code) : setUser(undefined)
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
                ExportLogToCSV(user.code);
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
                user && DeleteMatchLog(user.code);
                setTimeout(() => setMatches([]), 50);
                setUser(undefined);
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
        {!user && users && users.length > 0 && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.225 }}
            className="grid max-h-[340px] h-full m-4 justify-center content-start overflow-y-scroll gap-5"
          >
            {users.map((u) => (
              <button
                className="bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-xl backdrop-blur rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
                key={u.displayName}
                onClick={() => {
                  setUser(u);
                  getMatches(u.code);
                }}
              >
                {u.displayName}
              </button>
            ))}
          </motion.div>
        )}
        {matches && user && (
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
                  {matches.map((log) => (
                    <tr
                      key={`${log.date}-${log.time}`}
                      className="backdrop-blur group"
                    >
                      <td
                        onClick={() => filterLog("date", log.date)}
                        className="whitespace-nowrap text-left rounded-l-xl rounded-r-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2 hover:underline cursor-pointer"
                      >
                        {log.date}
                      </td>
                      <td className="whitespace-nowrap text-left bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                        {log.time}
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
                        style={{ color: log.victory == true ? "lime" : "red" }}
                      >
                        {log.victory == true ? "W" : "L"}
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
