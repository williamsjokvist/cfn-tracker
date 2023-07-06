import React from "react";
import { Trans, useTranslation } from "react-i18next";
import { motion } from "framer-motion";
import { Icon } from "@iconify/react";

import { PageHeader } from "@/ui/header";
import { ActionButton } from "@/ui/action-button";
import { GetThemeList, OpenResultsDirectory } from "@@/go/core/CommandHandler";

var defaultOptions = {
  theme: 'default',
  cfn: false,
  wins: true,
  losses: true,
  winRate: true,
  lpGain: true,
  lp: false,
  opponent: false,
  opponentCharacter: false,
  opponentLeague: false,
  opponentLP: false,
  totalLosses: false,
  totalMatches: false,
  totalWins: false,
  result: false,
  winStreak: false,
  timestamp: false,
  date: false,
}

export const OutputPage: React.FC = () => {
  const { t } = useTranslation();
  const [options, setOptions] = React.useState(defaultOptions)
  const [themes, setThemes] = React.useState<string[]>([])
  React.useEffect(() => {
    if (themes.length == 0) {
      GetThemeList().then(themes => setThemes(themes))
    }
  }, [])
  return (
    <>
      <PageHeader text={t('output')}/>
      <div className="z-40 h-full w-full overflow-y-scroll">
        <motion.section            
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className='grid gap-8 py-6 px-8 w-full border-b-[rgba(255,255,255,.125)] border-b-[1px] border-solid'
        >
          <div className='flex items-center justify-between gap-8'>
            <div>
              <h2 className="text-xl font-bold">{t('usingBrowserSource')}</h2>
                <p className='my-2 '>
                  <Trans t={t} i18nKey='browserSourceDescription'>
                    The easiest way to display the stats in OBS is to use <i className="whitespace-nowrap">Browser Source</i>. After editing the options, copy the browser source link.
                  </Trans>
                </p>
            </div>
            <div>
              <ActionButton 
                onClick={() => {
                  let url = `http://localhost:4242?theme=${options.theme}`
                  for (const [key, value] of Object.entries(options)){
                    if (key == 'theme') continue
                    url += `&${key}=${value}`
                  }
                  
                  navigator.clipboard.writeText(url)
                }} 
                style={{ filter: "hue-rotate(-120deg)" }}
              >
                <Icon icon='mdi:internet' className="mr-3 w-6 h-6" /> 
                {t('copyBrowserSourceLink')}
              </ActionButton>
            </div>
          </div>
          <div className='flex items-center justify-between gap-8'>
            <div>
              <b className='block'>{t('theme')}</b>
              <select 
                defaultValue='default' 
                onChange={(e) => { setOptions({...options, theme: e.target.value}) }}
                className='w-48 appearance-none bg-[rgb(0,0,0,.525)] rounded-lg border-none !drop-shadow-none'
              >
                <option value='default'>{t('defaultTheme')}</option>
                {themes && themes.map(theme => <option key={theme} value={theme}>{theme.charAt(0).toUpperCase() + theme.slice(1)}</option>)}
              </select>
            </div>
            <div>
              <h3 className='text-lg font-bold'>{t('pickTheme')}</h3>
              <p>
                <Trans t={t} i18nKey='createTheme'>
                  To create your own theme in CSS, create a <i>yourthemename.css</i> file in the <i>themes folder</i> and it will appear in the menu here. You can look at the <i>nord.css</i> file for reference.
                </Trans>
              </p>
            </div>
          </div>

          <div className="grid gap-2 items-center my-4 select-none">
            <h3 className='text-lg font-bold'>{t('displayStats')}</h3>
            {Object.entries(options).map(([key, value]) => {
              if (key == 'theme') return null
              const sfvOnly = key == 'totalLosses' || key == 'totalWins'
              return (
                <div key={key} className='flex gap-2'>
                  <input 
                    type='checkbox' 
                    id={`show-${key}`} 
                    defaultChecked={(value == true)} 
                    className='cursor-pointer'
                    onChange={e => {
                      const opt = options
                      opt[key] = e.target.value == 'on'
                      setOptions(opt)
                    }}
                  />
                  <label htmlFor={`show-${key}`} className='cursor-pointer'>
                    {key}
                    {sfvOnly && <i className='ml-2'>(SFV only)</i>}
                  </label>
                </div>
              )
            })}
          </div>
        </motion.section>
        <motion.section
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className='flex justify-between items-center gap-8 py-6 px-8 w-full'
        >
          <header>
            <h2 className="text-xl font-bold">{t('importFiles')}</h2>
            <p className='my-2'>
              <Trans t={t} i18nKey='obsCustomize'>
                If you want to customize through OBS, you can use <i>Text Labels</i> and add the text files in the results folder as sources              
              </Trans>
            </p>
          </header>
          <ActionButton onClick={OpenResultsDirectory} style={{ filter: "hue-rotate(-120deg)" }}>
            <Icon icon='fa6-solid:folder-open' className="mr-3 w-6 h-6" /> 
            {t("files")}
          </ActionButton>
        </motion.section>
      </div>
    </>
  );
};

