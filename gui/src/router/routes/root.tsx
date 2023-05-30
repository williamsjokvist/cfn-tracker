import React from "react";
import { useForm, useWatch } from 'react-hook-form'
import { useTranslation } from "react-i18next";
import { PieChart } from "react-minimal-pie-chart";
import { FaStop } from "react-icons/fa";
import { AiFillFolderOpen } from "react-icons/ai";

import {
  OpenResultsDirectory,
  GetAvailableLogs,
} from "@@/go/core/App";

import { CFNMachineContext } from '@/state-machine/machine'
import { PageHeader } from "@/ui/header";

type CfnFormValues = {
  cfn: string;
  restore?: boolean;
}

const EntryForm: React.FC = () => {
  const { t } = useTranslation();
  const [state, send] = CFNMachineContext.useActor();
  const [oldCfns, setOldCfns] = React.useState<string[] | null>(null);
  const [lastJSONExist, setLastJSONExist] = React.useState<boolean>(false);
  const cfnInputRef = React.useRef<HTMLInputElement>(null)

  React.useCallback(() => {
    GetAvailableLogs().then(logs => setOldCfns(logs))
    console.log(oldCfns)
  }, [oldCfns])

  const form = useForm<CfnFormValues>({
    defaultValues: {
      cfn: '',
      restore: false
    },
  })

  const { cfn, restore } = useWatch({
    control: form.control
  })

  const onSubmit = ({ cfn, restore }: CfnFormValues) => {
    state.context.cfn = cfn
    state.context.restore = restore

    send('submit')
  }

  return (
    <form
      className="max-w-[450px] mx-auto"
      onSubmit={form.handleSubmit(onSubmit)}
    >
      <h3 className="mb-2 text-lg">{t("enterCfnName")}:</h3>
      <input
        {...form.register('cfn')}
        ref={cfnInputRef}
        className="bg-transparent border-b-2 border-0 focus:ring-offset-transparent focus:ring-transparent border-b-[rgba(255,255,255,0.275)] focus:border-white hover:border-white outline-none focus:outline-none hover:text-white transition-colors py-3 px-4 block w-full text-lg text-gray-300"
        disabled={restore}
        type="text"
        placeholder={t("cfnName")!}
        autoCapitalize="off"
        autoComplete="off"
        autoCorrect="off"
        autoSave="off"
      />
      {oldCfns && (
        <div className="mt-3 flex flex-wrap gap-2 content-center items-center text-center pr-3">
          {oldCfns.map((cfn, index) => {
            return (
              <button
                disabled={restore}
                onClick={_ => {
                  cfnInputRef.current.value = cfn
                }}
                className="whitespace-nowrap bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-base backdrop-blur rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
                type="button"
                key={index}
              >
                {cfn}
              </button>
            );
          })}
        </div>
      )}
      {lastJSONExist && (
        <div className={`text-lg flex items-center mt-4`}>
          <input
            {...form.register('restore')}
            type="checkbox"
            className="w-7 h-7 rounded-md checked:border-2 checked:focus:border-[rgba(255,255,255,.25)] checked:hover:border-[rgba(255,255,255,.25)] checked:border-[rgba(255,255,255,.25)] border-2 border-[rgba(255,255,255,.25)] focus:border-2 cursor-pointer bg-transparent text-transparent focus:ring-offset-transparent focus:ring-transparent mr-4"
            onChange={e => {
              if (e.target.checked)
                cfnInputRef.current.value = ''
            }}
          />
          {t("restoreSession")}
        </div>
      )}
      <div className="flex justify-end">
        <button
          type="submit"
          className="mt-4 select-none text-center bg-[rgba(255,10,10,.1)] rounded-md px-7 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
          style={{ filter: "hue-rotate(156deg)" }}
        >
          {t("start")}
        </button>
      </div>
    </form>
  )
}

const Tracking: React.FC = () => {
  const { t } = useTranslation();
  const [state, send] = CFNMachineContext.useActor();
  const { cfn, lp, wins, losses, winStreak, lpGain, winRate } = state.context.matchHistory

  const openResultBtn = (
    <button
      onClick={() => OpenResultsDirectory()}
      style={{ filter: "hue-rotate(-120deg)" }}
      type="button"
      className="whitespace-nowrap flex items-center justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
    >
      <AiFillFolderOpen className="w-4 h-4 mr-2" />
      {t("openResultFolder")}
    </button>
  )

  const stopBtn = (
    <button
      onClick={() => send('stop')}
      type="button"
      className="flex items-center mb-2 justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
    >
      <FaStop className="mr-3" /> {t("stop")}
    </button>
  )

  const pieChart = (
    <PieChart
      className="pie-chart mt-6 animate-enter max-w-[200px] max-h-[200px] backdrop-blur"
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
        <h4 className="text-2xl">
          <span className="text-sm block">LP</span>
          {lp}
        </h4>
        <dl className="stat-grid-item w-full mt-2 relative text-center text-lg whitespace-nowrap">
          <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
            <dt className="tracking-wider font-extralight">{t("wins")}</dt>
            <dd className="text-4xl font-semibold">{wins}</dd>
          </div>
          <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
            <dt className="tracking-wide font-extralight">{t("losses")}</dt>
            <dd className="text-4xl font-semibold">{losses}</dd>
          </div>
          <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
            <dt className="tracking-wide font-extralight">{t("winStreak")}</dt>
            <dd className="text-4xl font-semibold">{winStreak}</dd>
          </div>
          <div className="mb-2 flex gap-4 justify-between bg-slate-50 bg-opacity-5 p-3 pb-1 rounded-xl backdrop-blur">
            <dt className="tracking-wide font-extralight">{t("lpGain")}</dt>
            <dd className="text-4xl font-semibold">{lpGain > 0 && "+"} {lpGain}</dd>
          </div>
        </dl>
      </section>
      <section className="relative mr-4 text-center h-full grid content-between justify-items-center">
        <b className='absolute top-[30%] z-50 text-4xl'>{(winRate > 0) && (winRate + '%')}</b>
        {pieChart}
        <div className="relative bottom-[10px] grid justify-items-end right-[-10px]">
          {stopBtn}
          {openResultBtn}
        </div>
      </section>
    </>
  ) 
}

export const RootPage: React.FC = () => {
  const { t } = useTranslation();
  const [state] = CFNMachineContext.useActor();
  const [headerText, setHeaderText] = React.useState<string>('')
  React.useEffect(() => {
    if (state.matches('initialized')) 
      setHeaderText(t("startTracking"))
    else if (state.matches('loading'))
      setHeaderText(t("loading"))
    else if (state.matches('tracking'))
      setHeaderText(t("tracking"))
    else if (state.matches('idle'))
      setHeaderText(t('startTracking'))
  }, [state])
  
  return (
    <>
      <PageHeader 
        text={headerText} 
        {...(state.matches('loading') && { showSpinner: true })} 
      />
      <div className="z-40 h-full w-full justify-self-center flex justify-between items-center px-8 py-4">
        {state.matches('tracking') && <Tracking/>}
        {state.matches('idle') && <EntryForm/>}
      </div>
    </>
  );
};

