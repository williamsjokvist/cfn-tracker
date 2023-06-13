import { CFNMachineContext } from "@/machine";
import { ActionButton } from "@/ui/action-button";
import { OpenResultsDirectory, StopTracking } from "@@/go/core/CommandHandler";
import { useTranslation } from "react-i18next";

import { FaFolderOpen } from "react-icons/fa";
import { FaStop } from "react-icons/fa";
import { PieChart } from "react-minimal-pie-chart";

type BigStatProps = {
  type: string
  value: string | number
}
const BigStat: React.FC<BigStatProps> = ({ type, value }) => {
  return (
    <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl">
      <dt className="tracking-wider font-extralight">{type}</dt>
      <dd className="text-4xl font-semibold">{value}</dd>
    </div>
  )
}

export const Tracking: React.FC = () => {
  const { t } = useTranslation();
  const [state, send] = CFNMachineContext.useActor();
  const { cfn, lp, wins, losses, winStreak, lpGain, winRate } = state.context.matchHistory

  const pieChart = (
    <PieChart
      className="pie-chart mt-16 animate-enter max-w-[180px] max-h-[180px] backdrop-blur"
      animate={true}
      lineWidth={75}
      paddingAngle={0}
      animationDuration={10}
      viewBoxSize={[60, 60]}
      center={[30, 30]}
      animationEasing={"ease-in-out"}
      data={[
        {
          title: "Wins",
          value: wins,
          color: "rgba(0, 255, 77, .65)",
        },
        {
          title: "Losses",
          value: losses,
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
    </PieChart>
  )
  
  return (
    <>
      <section className="relative w-full h-full grid grid-rows-[0fr_1fr] max-w-[320px]">
        <h3 className="whitespace-nowrap max-w-[145px] text-2xl">
          <span className="text-sm block">CFN</span>
          <span className="text-ellipsis block overflow-hidden">{cfn}</span>
        </h3>
        <h3 className="text-2xl">
          <span className="text-sm block">LP</span>
          {`${lp == -1 ? t('placement') : lp}`}
        </h3>
        <dl className="stat-grid-item w-full mt-2 relative text-center text-lg whitespace-nowrap">
          <BigStat type={t("wins")} value={wins}/>
          <BigStat type={t("losses")} value={losses}/>
          <BigStat type={t("winStreak")} value={winStreak}/>
          <BigStat type={t("lpGain")} value={`${(lpGain > 0) ? `+` : ``}${lpGain}`}/>
        </dl>
      </section>
      <section className="relative text-center w-[300px] h-full grid content-between justify-items-center">
        <b className='absolute top-[10px] z-50 text-4xl'>{(winRate > 0) && (winRate + '%')}</b>
        {pieChart}
        <div className="relative bottom-[10px] flex items-start gap-5">
          <ActionButton onClick={OpenResultsDirectory} style={{ filter: "hue-rotate(-120deg)" }}>
            <FaFolderOpen className="w-5 h-5 mr-2" />
            {t("files")}
          </ActionButton>
          <ActionButton onClick={() =>{
            StopTracking() // TODO: this should be part of the state machine
            send('stoppedTracking')}
          }>
            <FaStop className="mr-3" /> 
            {t("stop")}
          </ActionButton>
        </div>
      </section>
    </>
  ) 
}