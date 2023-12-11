import React from "react";
import { useTranslation } from "react-i18next";

import { GetSessions } from "@@/go/core/CommandHandler";
import type { model } from "@@/go/models";

import { PageHeader } from "@/ui/page-header";
import { useNavigate } from "react-router-dom";
import { motion } from "framer-motion";

export const SessionsListPage: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [sessions, setSessions] = React.useState<model.Session[]>([]);

  React.useEffect(() => {
    GetSessions("").then((seshs) => setSessions(seshs));
  }, []);

  return (
    <>
      <PageHeader text={"Sessions"} />
      <div className="overflow-y-scroll max-h-[340px] h-full mx-4 px-4 pb-4 mt-3">
        <motion.table
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.25 }}
          className="w-full border-spacing-y-1 border-separate min-w-[525px]"
        >
          <thead>
            <tr>
              <th className="text-left px-3 whitespace-nowrap w-[120px]">
                Started
              </th>
              <th className="text-left px-3 whitespace-nowrap w-[50%]">User</th>
              <th className="text-left px-3 whitespace-nowrap">LP</th>
              <th className="text-center px-3 whitespace-nowrap">MR</th>
            </tr>
          </thead>
          <tbody>
            {sessions.map((sesh) => (
              <tr
                key={sesh.id}
                className="backdrop-blur group cursor-pointer"
                onClick={() => navigate(`/sessions/${sesh.id}`)}
              >
                <td className="whitespace-nowrap text-left rounded-l-xl rounded-r-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                  <time dateTime={sesh.createdAt}>
                    {new Date(sesh.createdAt).toLocaleDateString("en-GB", {
                      weekday: "short",
                      month: "numeric",
                      day: "numeric",
                      year: "2-digit",
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </time>
                </td>
                <td className="whitespace-nowrap text-left bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                  {sesh.userName}
                </td>
                <td className="whitespace-nowrap text-left bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2">
                  {sesh.lp}
                </td>
                <td className="rounded-r-xl rounded-l-none bg-slate-50 bg-opacity-5 group-hover:bg-opacity-10 transition-colors px-3 py-2 text-center">
                  {sesh.mr}
                </td>
              </tr>
            ))}
          </tbody>
        </motion.table>
      </div>
    </>
  );
};
