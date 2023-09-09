import React from "react";
import { Trans, useTranslation } from "react-i18next";
import { motion, useAnimate } from "framer-motion";
import { Icon } from "@iconify/react";

import { PageHeader } from "@/ui/page-header";
import { ActionButton } from "@/ui/action-button";
import { GetThemeList, OpenResultsDirectory } from "@@/go/core/CommandHandler";
import { Checkbox } from "@/ui/checkbox";

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
      { opacity: [0, 1], y: [-30, -40] },
      { delay: 0.125 }
    ).then(() => {
      animate("#ok",
        { opacity: [1, 0], y: [-40, -40] },
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
          className="grid gap-8 py-6 px-8 w-full border-b-[rgba(255,255,255,.125)] border-b-[1px] border-solid"
        >
          <div className="flex items-center justify-between gap-8">
            <div>
              <h2 className="text-2xl font-bold">{t("usingBrowserSource")}</h2>
              <p className="my-2 ">
                <Trans t={t} i18nKey="browserSourceDescription">
                  The easiest way to display the stats in OBS is to use
                  <i className="whitespace-nowrap">Browser Source</i>. After
                  editing the options, copy the browser source link.
                </Trans>
              </p>
            </div>
            <div ref={scope}>
              <motion.span initial={{ opacity: 0 }} id="ok" className='h-0 top-2 -left-6 block relative text-2xl font-bold'>👌</motion.span>
              <ActionButton
                onClick={copyUrlToClipBoard}
                style={{ filter: "hue-rotate(-120deg)" }}
              >
                <Icon icon="mdi:internet" className="mr-3 w-6 h-6" />
                {t("copyBrowserSourceLink")}
              </ActionButton>
            </div>
          </div>
          <div className="flex items-center justify-between gap-8">
            <div>
              <b className="block">{t("theme")}</b>
              <select
                defaultValue="default"
                onChange={(e) => {
                  setOptions({ ...options, theme: e.target.value });
                }}
                className="w-48 appearance-none bg-[rgb(0,0,0,.525)] rounded-lg border-none !drop-shadow-none cursor-pointer"
              >
                <option value="default">{t("defaultTheme")}</option>
                {themes &&
                  themes.map((theme) => (
                    <option key={theme} value={theme}>
                      {theme.charAt(0).toUpperCase() + theme.slice(1)}
                    </option>
                  ))}
              </select>
            </div>
            <div>
              
              <ActionButton
                onClick={() => {

                }}
                style={{ filter: "hue-rotate(-320deg)" }}
              >
                <Icon icon="ph:paint-bucket-fill" className="mr-3 w-6 h-6" />
                Pick Theme
              </ActionButton>
            </div>
          </div>

          <div className="grid items-center my-4 select-none">
            <h3 className="text-xl font-bold mb-2">{t("displayStats")}</h3>
            <ul>
              {Object.entries(options).map(([key, value]) => {
                if (key == "theme") return null;

                return (
                  <li key={key}>
                    <button
                      className="w-full cursor-pointer flex py-1 px-6 justify-between items-center text-lg hover:bg-[rgba(255,255,255,0.075)]"
                      onClick={() => {
                        setOptions({ ...options, [key]: !value });
                      }}
                    >
                      <div>
                        <Checkbox checked={options[key] == true} readOnly />
                        <span className="ml-2 text-center cursor-pointer capitalize">
                          {key}
                        </span>
                      </div>
                    </button>
                  </li>
                );
              })}
            </ul>
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
            <p className="my-2">
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
          >
            <Icon icon="fa6-solid:folder-open" className="mr-3 w-6 h-6" />
            {t("files")}
          </ActionButton>
        </motion.section>
      </div>
    </>
  );
};
