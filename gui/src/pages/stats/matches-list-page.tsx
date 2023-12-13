import React from "react";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";

import { GetMatches } from "@@/go/core/CommandHandler";
import type { model } from "@@/go/models";

import { PageHeader } from "@/ui/page-header";
import { useParams } from "react-router-dom";

export const MatchesListPage: React.FC = () => {
  const { t } = useTranslation();
  const params = useParams();
  const sessionId = Number(params["sessionId"]);

  const [matches, setMatches] = React.useState<model.Match[]>([]);
  const [totalWinRate, setTotalWinRate] = React.useState<number | null>(null);

  React.useEffect(() => {
    GetMatches(sessionId, "", 0, 0)
      .then((matches) => {
        setMatches(matches);
      })
      .catch((err) => console.error(err));
  }, [sessionId])

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
  };

  return (
    <>
      <PageHeader text={t("history")} />

      <div className="relative w-full">
        {matches && (
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
