import React from "react";
import { Trans, useTranslation } from "react-i18next";
import { motion, useAnimate } from "framer-motion";
import { Icon } from "@iconify/react";

import { PageHeader } from "@/ui/page-header";
import { ActionButton } from "@/ui/action-button";
import { GetThemeList, OpenResultsDirectory } from "@@/go/core/CommandHandler";
import { Checkbox } from "@/ui/checkbox";
import * as Dialog from "@radix-ui/react-dialog";

import defaultThemePreview from '@/img/theme-previews/default.png'
import bladesThemePreview from '@/img/theme-previews/blades.png'
import jaegerThemePreview from '@/img/theme-previews/jaeger.png'
import nordThemePreview from '@/img/theme-previews/nord.png'
import pillsThemePreview from '@/img/theme-previews/pills.png'

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
  const [themes, setThemes] = React.useState<string[]>([]);
  const [scope, animate] = useAnimate();

  React.useEffect(() => {
    if (themes.length == 0) {
      GetThemeList().then((themes) => setThemes(themes));
    }
  }, []);

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
              <ThemePickModal 
                themes={themes}
                onSelect={theme => setOptions({ ...options, theme }) } 
                selectedTheme={options.theme}
              />
              <StatPickModal 
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

const CloseButton = () => (
  <Dialog.Close className="bg-[#202020] hover:bg-[#2b2a33] w-11 h-11 rounded-full">
    <Icon icon="ci:close-big" width={28} className='mx-auto' />
  </Dialog.Close>
)

const ThemePickModal = (props: { selectedTheme: string, themes: string[], onSelect: (theme: string) => void }) => {
  const { t } = useTranslation()

  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <ActionButton style={{ filter: "hue-rotate(-180deg)", justifyContent: "center" }}>
          <Icon icon="ph:paint-bucket-fill" className="mr-3 w-6 h-6" />
          {t("selectTheme")}
        </ActionButton>
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Content className="DialogContent">
          <header className="flex justify-between">
            <Dialog.Title className="text-2xl font-bold">
              {t("selectTheme")}
            </Dialog.Title>
            <CloseButton />
          </header>
          <Dialog.Description>
          </Dialog.Description>

          <ul className='mt-2 h-80 w-full pr-2 overflow-y-scroll'>
            <li>
              <button onClick={e => {
                props.onSelect('default')
              }}>
                <img src={defaultThemePreview} className='object-cover object-top overflow-hidden w-full rounded-md h-[76px]' 
                  style={{ border: props.selectedTheme == 'default' ? '2px solid lime' : '2px solid transparent'}}
                />
              </button>
            </li>
            <li>
              <button onClick={e => {
                props.onSelect('blades')
              }}>
                <img src={bladesThemePreview} className='object-cover object-top overflow-hidden w-full rounded-md h-[67px]'
                  style={{ border: props.selectedTheme == 'blades' ? '2px solid lime' : '2px solid transparent'}}
                />
              </button>
            </li>
            <li>
              <button onClick={e => {
                props.onSelect('jaeger')
              }}>
                <img src={jaegerThemePreview} className='object-cover object-top overflow-hidden w-full rounded-md h-[60px]'
                  style={{ border: props.selectedTheme == 'jaeger' ? '2px solid lime' : '2px solid transparent'}}
                />
              </button>
            </li>
            <li>
              <button onClick={e => {
                props.onSelect('pills')
              }}>
                <img src={pillsThemePreview} className='object-cover object-top overflow-hidden w-full rounded-md h-[60px]' 
                  style={{ border: props.selectedTheme == 'pills' ? '2px solid lime' : '2px solid transparent'}}
                />
              </button>
            </li>
            <li>
              <button onClick={e => {
                props.onSelect('nord')
              }}>
                <img src={nordThemePreview} className='object-cover object-top overflow-hidden w-full rounded-md h-[90px]' 
                  style={{ border: props.selectedTheme == 'nord' ? '2px solid lime' : '2px solid transparent'}}
                />
              </button>
            </li>
          </ul>
        </Dialog.Content>
        <Dialog.Overlay className="DialogOverlay" />
      </Dialog.Portal>
    </Dialog.Root>
  )
}

const StatPickModal = (props: { onSelect: (option: string, checked: boolean) => void, options: typeof defaultOptions }) => {
  const { t } = useTranslation()
  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <ActionButton style={{ filter: "hue-rotate(-45deg)", justifyContent: "center" }}>
          <Icon icon="bx:stats" className="mr-3 w-6 h-6" />
          {t("displayStats")}
        </ActionButton>
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Content className="DialogContent">
          <header className="flex justify-between">
            <Dialog.Title className="text-2xl font-bold">
              {t("displayStats")}
            </Dialog.Title>
            <CloseButton />
          </header>
          <Dialog.Description className="text-lg leading-none">
            {t('statsWillBeDisplayed')}
          </Dialog.Description>

          <ul className="overflow-y-scroll h-72 mt-4">
            {Object.entries(props.options).map(([key, value]) => {
              if (key == "theme") return null;
              return (
                <li key={key}>
                  <button
                    className="w-full cursor-pointer flex py-1 px-2 justify-between items-center text-lg hover:bg-[rgba(255,255,255,0.075)]"
                    onClick={() => props.onSelect(key, !value)}
                  >
                    <div>
                      <Checkbox checked={props.options[key] == true} readOnly />
                      <span className="ml-2 text-center cursor-pointer capitalize">
                        {key}
                      </span>
                    </div>
                  </button>
                </li>
              );
            })}
          </ul>
        </Dialog.Content>
        <Dialog.Overlay className="DialogOverlay" />
      </Dialog.Portal>
    </Dialog.Root>
  )
}