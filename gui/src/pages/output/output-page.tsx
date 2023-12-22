import React from "react";
import { Trans, useTranslation } from "react-i18next";
import { motion, useAnimate } from "framer-motion";
import { Icon } from "@iconify/react";

import { PageHeader } from "@/ui/page-header";
import { ActionButton } from "@/ui/action-button";
import { OpenResultsDirectory } from "@@/go/core/CommandHandler";

import { ThemeDialog } from "./theme-dialog";
import { StatsDialog } from "./stats-dialog";

export type StatOptions = typeof defaultOptions;

const defaultOptions = {
  theme: "default",
  cfn: false,
  wins: true,
  losses: true,
  winRate: true,
  lp: false,
  mr: true,
  lpGain: false,
  mrGain: true,
  opponent: false,
  opponentCharacter: false,
  opponentLeague: false,
  opponentLP: false,
  totalMatches: false,
  result: false,
  winStreak: false,
  timestamp: false,
  date: false,
};

export const OutputPage: React.FC = () => {
  const { t } = useTranslation();
  const [options, setOptions] = React.useState(defaultOptions);
  const [scope, animate] = useAnimate();

  const copyUrlToClipBoard = () => {
    let url = `http://localhost:4242?theme=${options.theme}`;
    for (const [key, value] of Object.entries(options)) {
      if (key == "theme") continue;
      url += `&${key}=${value}`;
    }

    navigator.clipboard.writeText(url);

    animate("#ok",
      { opacity: [0, 1], y: [0, -10] },
      { delay: 0.125 }
    ).then(() => {
      animate("#ok",
        { opacity: [1, 0] },
        { delay: .5 }
      );
    })

  };
  return (
    <>
      <PageHeader text={t("output")} />

      <div className="z-40 h-full w-full overflow-y-scroll">
        <motion.section
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className="grid gap-4 py-6 px-8 w-full border-b-[rgba(255,255,255,.125)] border-b-[1px] border-solid"
        >
          <div className="flex items-center justify-between gap-8">
            <div>
              <h2 className="text-2xl font-bold">{t("usingBrowserSource")}</h2>
              <p className="my-4 max-w-[340px]">
                <Trans t={t} i18nKey="browserSourceDescription">
                  The easiest way to display the stats in OBS is to use
                  <i className="whitespace-nowrap">Browser Source</i>. After
                  editing the options, copy the browser source link.
                </Trans>
              </p>
              <div className='relative' ref={scope}>
                <Icon icon="twemoji:ok-hand" width={45} id="ok" className='absolute top-0 right-7 opacity-0' /> 
              </div>
              <ActionButton
                onClick={copyUrlToClipBoard}
                style={{ filter: "hue-rotate(-120deg)" }}
              >
                <Icon icon="mdi:internet" className="mr-3 w-6 h-6" />
                {t("copyBrowserSourceLink")}
              </ActionButton>
            </div>
            <div className='flex flex-col gap-4'>
              <ThemeDialog 
                onSelect={theme => setOptions({ ...options, theme }) } 
                selectedTheme={options.theme}
              />
              <StatsDialog 
                options={options} 
                onSelect={(key, value) => setOptions({ ...options, [key]: value })}
              />
            </div>
          </div>
          <div className="flex justify-end items-center gap-8">

          </div>
        </motion.section>
        <motion.section
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className="flex justify-between items-center gap-8 py-6 px-8 w-full"
        >
          <header>
            <h2 className="text-xl font-bold">{t("importFiles")}</h2>
            <p className="mt-2 max-w-[420px]">
              <Trans t={t} i18nKey="obsCustomize">
                If you want to customize through OBS, you can use
                <i>Text Labels</i> and add the text files in the results folder
                as sources
              </Trans>
            </p>
          </header>
          <ActionButton
            onClick={OpenResultsDirectory}
            style={{ filter: "hue-rotate(-120deg)" }}
            className='mx-auto'
          >
            <Icon icon="fa6-solid:folder-open" className="mr-3 w-6 h-6" />
            {t("files")}
          </ActionButton>
        </motion.section>
      </div>
    </>
  );
};
