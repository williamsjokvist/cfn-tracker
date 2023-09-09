import { motion } from "framer-motion";
import { useTranslation } from "react-i18next";
import { Icon } from "@iconify/react";
import { PieChart as ReactMinimalPieChart } from "react-minimal-pie-chart";

import { CFNMachineContext } from "@/main/machine";
import { ActionButton } from "@/ui/action-button";
import { StopTracking } from "@@/go/core/CommandHandler";

type StatProps = {
  text: string;
  value: string | number;
};
const BigStat: React.FC<StatProps> = ({ text, value }) => (
  <div className="mb-2 flex flex-1 gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl">
    <dt className="tracking-wider font-extralight">{text}</dt>
    <dd className="text-4xl font-semibold">{value}</dd>
  </div>
);

const SmallStat: React.FC<StatProps> = ({ text, value }) => (
  <div className="flex gap-3 text-2xl">
    <dt className='text-xl leading-8'>{text}</dt>
    <dd className="font-bold">{value}</dd>
  </div>
)

export const TrackingLiveUpdater: React.FC = () => {
  const { t } = useTranslation();
  const [state, send] = CFNMachineContext.useActor();
  const {
    cfn,
    lp,
    mr,
    wins,
    losses,
    winStreak,
    lpGain,
    mrGain,
    winRate,
    opponent,
    opponentCharacter,
    opponentLeague,
    result,
  } = state.context.matchHistory;

  const PieChart = (
    <ReactMinimalPieChart
      className="pie-chart animate-enter w-[150px] h-[150px] mx-auto"
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
    <motion.section           
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.125 }} 
      className="px-6 pt-4 h-full"
    >
      <dl className="flex whitespace-nowrap items-start justify-between w-full">
        <SmallStat text="CFN" value={cfn} />
        <SmallStat text="LP" value={`${lp == -1 ? t("placement") : lp}`} />
        <SmallStat text="MR" value={`${mr == -1 ? t("placement") : mr}`} />
        <SmallStat text={t("winRate")} value={`${winRate}%`} />
      </dl>
      <div className="flex gap-12 pt-3 pb-5 h-[calc(100%-32px)]">
        <dl className="w-full text-lg whitespace-nowrap">
          <div className="flex justify-between gap-2">
            <BigStat text={t("wins")} value={wins} />
            <BigStat text={t("losses")} value={losses} />
          </div>
          <BigStat text={t("winStreak")} value={winStreak} />
          <div className="flex justify-between gap-2">
            <BigStat text={t("lpGain")} value={`${lpGain > 0 ? `+` : ``}${lpGain}`} />
            <BigStat text={t("mrGain")} value={`${mrGain > 0 ? `+` : ``}${mrGain}`} />
          </div>
        </dl>
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.35 }}
          className="relative h-full max-w-[280px] w-full text-center gap-5"
        >
          {PieChart}
          <div className="w-full max-w-[320px] mx-auto whitespace-nowrap text-lg mt-5">
            {opponent != "" && (
              <ul>
                <li className="group leading-none flex justify-between p-2">
                  <span>Last opponent</span>
                  <div className="relative">
                    <b>{opponent}</b>
                    <div className="group-hover:opacity-100 opacity-0 pointer-events-none absolute top-5 transition-all bg-[rgba(0,0,0,.525)] backdrop-blur-xl px-3 py-2 right-0 rounded-md">
                      <b className="mr-2">{opponentCharacter}</b>
                      <span>{`(${opponentLeague})`}</span>
                    </div>
                  </div>
                </li>
                <li className="leading-none flex justify-between p-2">
                  <span>Last match result</span>
                  <div>
                    <b>{result ? "WIN" : "LOSE"}</b>
                  </div>
                </li>
              </ul>
            )}
          </div>
          <ActionButton
            className="absolute bottom-0 right-0"
            onClick={() => {
              StopTracking(); // TODO: this should be part of the state machine
              send("stoppedTracking");
            }}
          >
            <Icon icon="fa6-solid:stop" className="mr-3 w-5 h-5" />
            {t("stop")}
          </ActionButton>
        </motion.div>
      </div>
    </motion.section>
  );
};
