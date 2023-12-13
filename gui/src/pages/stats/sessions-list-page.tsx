import React from "react";
import { useTranslation } from "react-i18next";

import { GetSessions } from "@@/go/core/CommandHandler";
import type { model } from "@@/go/models";

import { PageHeader } from "@/ui/page-header";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";

type MonthGroup = Record<number, model.Session[]>;
type YearGroup = Record<number, MonthGroup>;

export const SessionsListPage: React.FC = () => {
  const { i18n, t } = useTranslation();
  const navigate = useNavigate();
  const [sessions, setSessions] = React.useState<model.Session[]>([]);

  React.useEffect(() => {
    GetSessions("").then((seshs) => {
      seshs && setSessions(seshs);
    });
  }, []);

  const groupedSessions: YearGroup = sessions.reduce((group, sesh) => {
    const date = new Date(sesh.createdAt);
    const year = date.getFullYear();
    const month = date.getMonth();

    group[year] = group[year] ?? {};
    group[year][month] = group[year][month] ?? [];
    group[year][month].push(sesh);

    return group;
  }, {});

  return (
    <>
      <PageHeader text={t("sessions")} />
      <div className="overflow-y-scroll max-h-[340px] h-full mx-4 px-4 pb-4 mt-3">
        {Object.keys(groupedSessions).map((year) => (
          <motion.section
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.125 }}
          >
            <h2 className="text-4xl mt-2">{year}</h2>
            <div>
              {Object.keys(groupedSessions[year])
                .reverse()
                .map((month) => (
                  <section>
                    <h3 className="text-2xl mt-2">
                      {Intl.DateTimeFormat(i18n.resolvedLanguage, {
                        month: "long",
                      }).format(new Date(`1999-${month}-01`))}
                    </h3>
                    <table
                      className="w-full border-spacing-y-1 border-separate"
                    >
                      <thead>
                        <tr>
                          <th className="text-left px-3 whitespace-nowrap w-[120px]">
                            {t("started")}
                          </th>
                          <th className="text-left px-3 whitespace-nowrap w-full">
                            {t("user")}
                          </th>
                          <th className="text-left px-3 whitespace-nowrap">
                            {t("mrGain")}
                          </th>
                          <th className="text-left px-3 whitespace-nowrap">
                            {t("lpGain")}
                          </th>
                          <th className="text-center px-3 whitespace-nowrap">
                            {t("matchesWon")}
                          </th>
                          <th className="text-center px-3 whitespace-nowrap">
                            {t("matchesLost")}
                          </th>
                        </tr>
                      </thead>
                      <tbody>
                        {groupedSessions[year][month].map((sesh) => (
                          <tr
                            key={sesh.id}
                            className="backdrop-blur group cursor-pointer"
                            onClick={() => navigate(`/sessions/${sesh.id}`)}
                          >
                            <td className="whitespace-nowrap text-left rounded-l-xl rounded-r-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                              <time dateTime={sesh.createdAt}>
                                {new Date(sesh.createdAt).toLocaleDateString(
                                  i18n.resolvedLanguage ?? "en-GB",
                                  {
                                    day: "numeric",
                                    weekday: "short",
                                    hour: "2-digit",
                                    minute: "2-digit",
                                  },
                                )}
                              </time>
                            </td>
                            <td className="whitespace-nowrap text-left bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                              {sesh.userName}
                            </td>
                            <td className="whitespace-nowrap text-center bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                              {sesh.mrGain > 0 && "+"}
                              {sesh.mrGain}
                            </td>
                            <td className="whitespace-nowrap text-center bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                              {sesh.lpGain > 0 && "+"}
                              {sesh.lpGain}
                            </td>
                            <td className="whitespace-nowrap text-center bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                              {sesh.matchesWon}
                            </td>
                            <td className="rounded-r-xl rounded-l-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2 text-center">
                              {sesh.matchesLost}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </section>
                ))}
            </div>
          </motion.section>
        ))}
      </div>
    </>
  );
};
