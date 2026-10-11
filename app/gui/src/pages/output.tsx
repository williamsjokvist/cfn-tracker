import React from 'react'
import { createRoot } from 'react-dom/client'
import { Icon } from '@iconify/react'
import { Trans, useTranslation } from 'react-i18next'
import { motion, useAnimate } from 'framer-motion'

import * as Page from '@/ui/page'
import * as Dialog from '@/ui/dialog'
import { SettingButton, settingTextClass } from '@/ui/setting-button'

import { Checkbox } from '@/ui/checkbox'

import { model } from '@model'
import { GetThemes, OpenResultsDirectory } from '@cmd/CommandHandler'
import { BrowserOpenURL } from '@runtime'

import { useErrorPopup } from '@/main/error-popup'
import { ConfigContext } from '@/main/config'
import type { LocalizationKey } from '@/main/i18n'
import { LanguageSelect } from '@/main/language-select'
import { cn } from '@/helpers/cn'

type StatOptions = Omit<Record<keyof model.Match, boolean>, 'replayId' | 'sessionId'> & {
  theme: string
}

// Custom themes from the user's themes folder have no translation, so they show their file name
const themeLabels: Record<string, LocalizationKey> = {
  default: 'themeDefault',
  blades: 'browserSourceThemeBlades',
  jaeger: 'browserSourceThemeJaeger',
  nord: 'browserSourceThemeNord',
  pills: 'browserSourceThemePills'
}

const statLabels: Record<keyof Omit<StatOptions, 'theme'>, LocalizationKey> = {
  userName: 'user',
  userId: 'userId',
  character: 'character',
  lp: 'lp',
  lpGain: 'lpGain',
  mr: 'mr',
  mrGain: 'mrGain',
  opponent: 'opponent',
  opponentCharacter: 'opponentCharacter',
  opponentLp: 'opponentLp',
  opponentMr: 'opponentMr',
  opponentLeague: 'opponentLeague',
  victory: 'result',
  date: 'date',
  time: 'time',
  winStreak: 'winStreak',
  wins: 'wins',
  losses: 'losses',
  winRate: 'winRate'
}

const defaultOptions: StatOptions = {
  theme: 'default',
  userName: false,
  wins: true,
  losses: true,
  winRate: true,
  winStreak: false,
  lp: false,
  mr: true,
  lpGain: false,
  mrGain: true,
  opponent: false,
  opponentCharacter: false,
  opponentLeague: false,
  opponentLp: false,
  opponentMr: false,
  character: false,
  victory: false,
  userId: false,
  time: false,
  date: false
}

export function OutputPage() {
  const { t } = useTranslation()
  const [options, setOptions] = React.useState(defaultOptions)
  // English rather than the app language, since the stream's audience may not share it
  const [statsLocale, setStatsLocale] = React.useState('en-GB')
  const [scope, animate] = useAnimate()
  const [cfg] = React.useContext(ConfigContext)

  const { theme: _, ...statOptions } = options

  function copyUrlToClipBoard() {
    let url = `http://localhost:${cfg.browserSourcePort}?theme=${options.theme}`
    for (const [key, value] of Object.entries(options)) {
      if (key === 'theme') {
        continue
      }
      url += `&${key}=${value}`
    }
    url += `&lang=${statsLocale}`

    navigator.clipboard.writeText(url)

    animate('#ok', { opacity: [0, 1], y: [0, -10] }, { delay: 0.125 }).then(() => {
      animate('#ok', { opacity: [1, 0] }, { delay: 0.5 })
    })
  }

  return (
    <Page.Root>
      <Page.Header>
        <div className='mx-auto w-full max-w-xl text-start'>
          <Page.Title>{t('output')}</Page.Title>
        </div>
      </Page.Header>
      <motion.section
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='overflow-y-scroll text-xl'
      >
        <div className='border-b-divider border-b-[1px] border-solid px-8 py-6'>
          <div className='mx-auto grid max-w-xl gap-3'>
            <div>
              <h3 className='font-bold'>{t('usingBrowserSource')}</h3>
              <p className={settingTextClass}>
                <Trans t={t} i18nKey='browserSourceDescription'>
                  The easiest way to display the stats in OBS is to use
                  <i className='whitespace-nowrap'>Browser Source</i>. After editing the options,
                  copy the browser source link.
                </Trans>
              </p>
            </div>
            <div className='flex flex-wrap gap-2'>
              <ThemeSelect
                value={options.theme}
                onSelect={theme => setOptions({ ...options, theme })}
              />
              <StatSelect
                options={statOptions}
                onSelect={(opt, checked) => setOptions({ ...options, [opt]: checked })}
              />
              <LanguageSelect value={statsLocale} onValueChange={setStatsLocale} />
            </div>
            <div className='flex flex-wrap gap-2'>
              <div className='relative' ref={scope}>
                <Icon
                  icon='twemoji:ok-hand'
                  width={36}
                  id='ok'
                  className='pointer-events-none absolute end-2 -top-10 opacity-0'
                />
                <SettingButton icon='mdi:internet' onClick={copyUrlToClipBoard}>
                  {t('copyBrowserSourceLink')}
                </SettingButton>
              </div>
              <SettingButton
                icon='mdi:brush'
                onClick={() => BrowserOpenURL('https://cfn.williamsjokvist.se/theme-guide/')}
              >
                {t('createYourOwnTheme')}
              </SettingButton>
            </div>
          </div>
        </div>
        <div className='px-8 py-6'>
          <div className='mx-auto grid max-w-xl gap-3'>
            <div>
              <h3 className='font-bold'>{t('importFiles')}</h3>
              <p className={settingTextClass}>
                <Trans t={t} i18nKey='obsCustomize'>
                  If you want to customize through OBS, you can use
                  <i>Text Labels</i> and add the text files in the results folder as sources
                </Trans>
              </p>
            </div>
            <div className='flex flex-wrap gap-2'>
              <SettingButton icon='mdi:folder-open-outline' onClick={OpenResultsDirectory}>
                {t('files')}
              </SettingButton>
            </div>
          </div>
        </div>
      </motion.section>
    </Page.Root>
  )
}

function StatSelect(props: {
  options: Omit<StatOptions, 'theme'>
  onSelect: (option: string, checked: boolean) => void
}) {
  const { t } = useTranslation()
  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <SettingButton icon='bx:stats'>{t('displayStats')}</SettingButton>
      </Dialog.Trigger>
      <Dialog.Content title='displayStats' description='statsWillBeDisplayed'>
        <ul className='mt-4 h-72 overflow-y-scroll'>
          {Object.entries(props.options).map(([opt, checked]) => (
            <li key={opt}>
              <button
                className='flex w-full cursor-pointer items-center px-2 py-1 text-lg hover:bg-white/[0.075]'
                onClick={() => props.onSelect(opt, !checked)}
              >
                <Checkbox checked={props.options[opt] === true} readOnly />
                <span className='ms-4 cursor-pointer text-center'>{t(statLabels[opt])}</span>
              </button>
            </li>
          ))}
        </ul>
      </Dialog.Content>
    </Dialog.Root>
  )
}

function ThemeSelect(props: { value: string; onSelect: (theme: string) => void }) {
  const { t } = useTranslation()
  const [themes, setThemes] = React.useState<model.Theme[]>([])
  const [isOpen, setOpen] = React.useState(false)
  const setError = useErrorPopup()

  const selectedTheme = themes.find(t => t.name === props.value) ?? themes[0]

  React.useEffect(() => {
    GetThemes().then(setThemes).catch(setError)
  }, [])

  React.useEffect(() => {
    if (!isOpen) {
      return
    }
    for (const theme of themes) {
      const container = document.querySelector(`.${theme.name}-preview`)
      if (!container) {
        continue
      }
      const shadowRoot = container.attachShadow({ mode: 'open' })
      createRoot(shadowRoot).render(
        <div className='stat-list'>
          <style>{theme.css}</style>
          <div className='stat-item'>
            <span className='stat-title'>MR</span>
            <span className='stat-value'>444</span>
          </div>
        </div>
      )
    }
  }, [isOpen, themes])

  return (
    <Dialog.Root onOpenChange={setOpen}>
      <Dialog.Trigger asChild>
        <SettingButton icon='ph:paint-bucket-fill' className='capitalize'>
          {themeLabels[props.value] ? t(themeLabels[props.value]) : props.value}
        </SettingButton>
      </Dialog.Trigger>
      <Dialog.Content title='selectTheme'>
        <ul className='my-4 w-full overflow-y-scroll'>
          {themes.map(theme => (
            <li
              key={theme.name}
              className='relative flex w-full cursor-pointer items-center text-lg hover:bg-white/[0.075]'
            >
              <Checkbox
                id={`${theme.name}-checkbox`}
                checked={theme.name === props.value}
                onChange={e => props.onSelect(theme.name)}
                className='absolute start-2 top-1'
              />
              <label
                htmlFor={`${theme.name}-checkbox`}
                className='font-spartan w-full cursor-pointer py-1 ps-14 capitalize'
              >
                {themeLabels[theme.name] ? t(themeLabels[theme.name]) : theme.name}
              </label>
            </li>
          ))}
        </ul>
        <div className='relative mx-auto h-[60px] w-[350px]'>
          {themes.map(theme => (
            <div
              key={theme.name}
              style={{ opacity: theme.name === selectedTheme.name ? '100' : '0' }}
              className={cn(
                `${theme.name}-preview`,
                'absolute start-0 top-0 w-full transition-opacity',
                'pointer-events-none mx-auto h-[60px] w-[350px] select-none'
              )}
            >
              <style>{theme.css.match(/@import.*;/g)}</style>
            </div>
          ))}
        </div>
      </Dialog.Content>
    </Dialog.Root>
  )
}
