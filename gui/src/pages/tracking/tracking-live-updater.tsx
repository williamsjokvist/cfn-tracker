import { motion } from "framer-motion";
import { useTranslation } from "react-i18next";
import { Icon } from "@iconify/react";
import { PieChart as ReactMinimalPieChart } from "react-minimal-pie-chart";

import { CFNMachineContext } from "@/main/machine";
import { ActionButton } from "@/ui/action-button";
import { StopTracking } from "@@/go/core/CommandHandler";

type BigStatProps = {
  text: string;
  value: string | number;
};
const BigStat: React.FC<BigStatProps> = ({ text, value }) => {
  return (
    <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl">
      <dt className="tracking-wider font-extralight">{text}</dt>
      <dd className="text-4xl font-semibold">{value}</dd>
    </div>
  );
};

export const TrackingLiveUpdater: React.FC = () => {
  const { t } = useTranslation();
  const [state, send] = CFNMachineContext.useActor();
  const {
    cfn,
    lp,
    wins,
    losses,
    winStreak,
    lpGain,
    winRate,
    opponent,
    opponentCharacter,
    opponentLeague,
    result,
  } = state.context.matchHistory;

  const PieChart = (
    <ReactMinimalPieChart
      className="pie-chart animate-enter w-[45px] backdrop-blur"
      animate
      animationDuration={750}
      lineWidth={85}
      paddingAngle={0}
      viewBoxSize={[60, 60]}
      center={[30, 30]}
      animationEasing="ease-in-out"
      data={[
        {
          title: t("wins"),
          value: wins,
          color: "rgba(0, 255, 77, .65)",
        },
        {
          title: t("losses"),
          value: wins == 0 && losses == 0 ? 1 : losses,
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
    </ReactMinimalPieChart>
  );

  return (
    <>
      <motion.section
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className="relative w-full h-full max-w-[360px]"
      >
        <header className="flex whitespace-nowrap items-start justify-between w-full">
          <h3 className="whitespace-nowrap max-w-[145px] text-2xl">
            <span className="text-sm block">CFN</span>
            <b className="text-ellipsis block overflow-hidden">{cfn}</b>
          </h3>
          <h3 className="text-2xl">
            <span className="text-sm block">LP</span>
            <b>{`${lp == -1 ? t("placement") : lp}`}</b>
          </h3>
          <div className="relative flex gap-5 text-right">
            <h3 className="text-2xl relative">
              <span className="text-sm block">{t("winRate")}</span>
              <b>{winRate + "%"}</b>
            </h3>
            {PieChart}
          </div>
        </header>
        <dl className="stat-grid-item w-full mt-2 relative text-center text-lg whitespace-nowrap">
          <BigStat text={t("wins")} value={wins} />
          <BigStat text={t("losses")} value={losses} />
          <BigStat text={t("winStreak")} value={winStreak} />
          <BigStat
            text={t("lpGain")}
            value={`${lpGain > 0 ? `+` : ``}${lpGain}`}
          />
        </dl>
      </motion.section>
      <motion.section
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.35 }}
        className="relative text-center w-[300px] h-full grid content-between justify-items-end"
      >
        <div className="w-full max-w-[280px] mx-auto whitespace-nowrap text-md">
          {opponent != "" && (
            <ul>
              <li className="group leading-none flex justify-between p-2 border-b-[1px] border-white border-opacity-25 border-solid">
                <span>Last opponent</span>
                <div className="relative">
                  <b>{opponent}</b>
                  <div className="group-hover:opacity-100 opacity-0 pointer-events-none absolute top-5 transition-all bg-[rgba(0,0,0,.525)] backdrop-blur-xl px-3 py-2 right-0 rounded-md">
                    <b className="mr-2">{opponentCharacter}</b>
                    <span>{`(${opponentLeague})`}</span>
                  </div>
                </div>
              </li>
              <li className="leading-none flex justify-between p-2 border-b-[1px] border-white border-opacity-25 border-solid">
                <span>Last match result</span>
                <div>
                  <b>{result ? "WIN" : "LOSE"}</b>
                </div>
              </li>
            </ul>
          )}
        </div>

        <div className="flex items-start justify-end w-full m-3 gap-5">
          <ActionButton
            onClick={() => {
              StopTracking(); // TODO: this should be part of the state machine
              send("stoppedTracking");
            }}
          >
            <Icon icon="fa6-solid:stop" className="mr-3 w-5 h-5" />
            {t("stop")}
          </ActionButton>
        </div>
      </motion.section>
    </>
  );
};
