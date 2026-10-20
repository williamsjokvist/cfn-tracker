import React from 'react'
import { createRoot, type Root } from 'react-dom/client'
import { Icon } from '@iconify/react'
import { Trans, useTranslation } from 'react-i18next'
import { motion, useAnimate } from 'framer-motion'

import * as Page from '@/ui/page'
import * as Dialog from '@/ui/dialog'
import { SettingButton, settingTextClass } from '@/ui/setting-button'

import { Checkbox } from '@/ui/checkbox'

import { model } from '@model'
import {
  GetStreamOverlayConfig,
  GetThemes,
  OpenResultsDirectory,
  SaveStreamOverlayConfig
} from '@cmd/CommandHandler'
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
  banner: 'browserSourceThemeBanner',
  blackout: 'browserSourceThemeBlackout',
  blades: 'browserSourceThemeBlades',
  broadcast: 'browserSourceThemeBroadcast',
  chips: 'browserSourceThemeChips',
  comic: 'browserSourceThemeComic',
  crimson: 'browserSourceThemeCrimson',
  deco: 'browserSourceThemeDeco',
  'flash-portal': 'browserSourceThemeFlashPortal',
  glass: 'browserSourceThemeGlass',
  glitch: 'browserSourceThemeGlitch',
  hearts: 'browserSourceThemeHearts',
  jaeger: 'browserSourceThemeJaeger',
  link: 'browserSourceThemeLink',
  minimal: 'browserSourceThemeMinimal',
  'minimal-dark': 'browserSourceThemeMinimalDark',
  'most-wanted': 'browserSourceThemeMostWanted',
  'most-wanted-strip': 'browserSourceThemeMostWantedStrip',
  neon: 'browserSourceThemeNeon',
  nord: 'browserSourceThemeNord',
  notebook: 'browserSourceThemeNotebook',
  pills: 'browserSourceThemePills',
  pivot: 'browserSourceThemePivot',
  retro: 'browserSourceThemeRetro',
  silver: 'browserSourceThemeSilver',
  sticker: 'browserSourceThemeSticker',
  terminal: 'browserSourceThemeTerminal',
  'text-only': 'browserSourceThemeTextOnly',
  ticker: 'browserSourceThemeTicker',
  'top-8': 'browserSourceThemeTop8',
  ultra: 'browserSourceThemeUltra',
  vapor: 'browserSourceThemeVapor',
  velvet: 'browserSourceThemeVelvet',
  versus: 'browserSourceThemeVersus',
  webcore: 'browserSourceThemeWebcore',
  zebra: 'browserSourceThemeZebra'
}

// Custom themes aren't listed here, so they show under rows
const sideBySideThemes = new Set([
  'banner',
  'chips',
  'flash-portal',
  'link',
  'most-wanted-strip',
  'nord',
  'sticker',
  'ticker',
  'top-8',
  'velvet'
])

type ThemeLayout = 'rows' | 'sideBySide'
const themeLayouts: { key: ThemeLayout; label: LocalizationKey }[] = [
  { key: 'rows', label: 'themeLayoutRows' },
  { key: 'sideBySide', label: 'themeLayoutSideBySide' }
]
const layoutOf = (theme: string): ThemeLayout =>
  sideBySideThemes.has(theme) ? 'sideBySide' : 'rows'

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

// Saved stats are the ones that are on, so a stat added in a later version starts off
function optionsFromConfig(cfg: model.StreamOverlayConfig): StatOptions {
  if (!cfg.stats) {
    return { ...defaultOptions, theme: cfg.theme }
  }
  const options: StatOptions = { ...defaultOptions, theme: cfg.theme }
  for (const key of Object.keys(options) as (keyof StatOptions)[]) {
    if (key !== 'theme') {
      options[key] = cfg.stats.includes(key)
    }
  }
  return options
}

export function OutputPage() {
  const { t } = useTranslation()
  const [options, setOptions] = React.useState(defaultOptions)
  // English rather than the app language, since the stream's audience may not share it
  const [statsLocale, setStatsLocale] = React.useState('en-GB')
  const [scope, animate] = useAnimate()
  const [cfg] = React.useContext(ConfigContext)
  const setError = useErrorPopup()

  React.useEffect(() => {
    GetStreamOverlayConfig()
      .then(overlay => {
        setOptions(optionsFromConfig(overlay))
        setStatsLocale(overlay.lang)
      })
      .catch(setError)
  }, [])

  function save(options: StatOptions, lang: string) {
    setOptions(options)
    setStatsLocale(lang)
    const { theme, ...stats } = options
    SaveStreamOverlayConfig({
      theme,
      lang,
      stats: Object.entries(stats)
        .filter(([, on]) => on)
        .map(([key]) => key)
    }).catch(setError)
  }

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
        <Page.Title>{t('output')}</Page.Title>
        <SettingButton
          icon='mdi:brush'
          // Clickable although the header is the window's drag region.
          style={{ '--draggable': 'no-drag' } as React.CSSProperties}
          onClick={() => BrowserOpenURL('https://cfn.williamsjokvist.se/theme-guide/')}
        >
          {t('createYourOwnTheme')}
        </SettingButton>
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
            <div className='grid auto-rows-fr gap-2'>
              <div className='flex items-center justify-between'>
                <span className={settingTextClass}>{t('language')}</span>
                <LanguageSelect value={statsLocale} onValueChange={lang => save(options, lang)} />
              </div>
              <div className='flex items-center justify-between'>
                <span className={settingTextClass}>{t('displayStats')}</span>
                <StatSelect
                  options={statOptions}
                  onSelect={(opt, checked) => save({ ...options, [opt]: checked }, statsLocale)}
                />
              </div>
              <div className='flex items-center justify-between'>
                <span className={settingTextClass}>{t('theme')}</span>
                <ThemeSelect
                  value={options.theme}
                  lang={statsLocale}
                  onSelect={theme => save({ ...options, theme }, statsLocale)}
                />
              </div>
            </div>
            <div className='flex flex-wrap justify-end gap-2'>
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

// The same width as the language select next to them
const rowControlClass = 'w-48 text-start'

function StatSelect(props: {
  options: Omit<StatOptions, 'theme'>
  onSelect: (option: string, checked: boolean) => void
}) {
  const { t } = useTranslation()
  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <SettingButton icon='bx:stats' className={rowControlClass}>
          <span className='truncate'>
            {Object.entries(props.options)
              .filter(([, on]) => on)
              .map(([opt]) => t(opt as keyof typeof props.options))
              .join(', ') || '–'}
          </span>
        </SettingButton>
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
                <span className='ms-4 cursor-pointer text-center'>
                  {t(opt as keyof typeof props.options)}
                </span>
              </button>
            </li>
          ))}
        </ul>
      </Dialog.Content>
    </Dialog.Root>
  )
}

// Strips are laid out on one line at their natural size, then scaled down to fit the preview box
function StripPreview(props: { css: string; stats: [string, string][] }) {
  const ref = React.useRef<HTMLDivElement>(null)
  const [scale, setScale] = React.useState(1)

  React.useLayoutEffect(() => {
    const el = ref.current
    if (!el) {
      return
    }
    // Re-measures once the theme's web fonts arrive and change the strip's size
    const observer = new ResizeObserver(() => {
      setScale(Math.min(1, 350 / el.offsetWidth, 60 / el.offsetHeight))
    })
    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  return (
    <div
      ref={ref}
      style={{
        position: 'absolute',
        left: '50%',
        top: 0,
        display: 'flow-root',
        width: 'max-content',
        transform: `translateX(-50%) scale(${scale})`,
        transformOrigin: 'top center'
      }}
    >
      <style>{props.css}</style>
      <div className='stat-list'>
        {props.stats.map(([title, value]) => (
          <div key={title} className='stat-item'>
            <span className='stat-title'>{title}</span>
            <span className='stat-value'>{value}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

const previewRoots = new WeakMap<Element, Root>()

function ThemeSelect(props: { value: string; lang: string; onSelect: (theme: string) => void }) {
  const { t, i18n } = useTranslation()
  const setError = useErrorPopup()
  // The previews are in the overlay's language, which only gets loaded if it's also the app's
  const [statsT, setStatsT] = React.useState(() => i18n.getFixedT(props.lang))
  React.useEffect(() => {
    i18n
      .loadLanguages(props.lang)
      .then(() => setStatsT(() => i18n.getFixedT(props.lang)))
      .catch(setError)
  }, [props.lang])
  const [themes, setThemes] = React.useState<model.Theme[]>([])
  const [isOpen, setOpen] = React.useState(false)
  const [layout, setLayout] = React.useState(layoutOf(props.value))

  const selectedTheme = themes.find(t => t.name === props.value) ?? themes[0]
  const stripPreviewStats: [string, string][] = [
    [statsT('mr'), '444'],
    [statsT('wins'), '12'],
    [statsT('losses'), '5']
  ]

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
      // Rendered again once the overlay language loads, and a shadow root can only be attached once
      let root = previewRoots.get(container)
      if (!root) {
        root = createRoot(container.attachShadow({ mode: 'open' }))
        previewRoots.set(container, root)
      }
      root.render(
        layoutOf(theme.name) === 'sideBySide' ? (
          <StripPreview css={theme.css} stats={stripPreviewStats} />
        ) : (
          <div className='stat-list'>
            <style>{theme.css}</style>
            <div className='stat-item'>
              <span className='stat-title'>{statsT('mr')}</span>
              <span className='stat-value'>444</span>
            </div>
          </div>
        )
      )
    }
  }, [isOpen, themes, statsT])

  return (
    <Dialog.Root
      onOpenChange={open => {
        setOpen(open)
        if (open) setLayout(layoutOf(props.value))
      }}
    >
      <Dialog.Trigger asChild>
        <SettingButton icon='ph:paint-bucket-fill' className={cn(rowControlClass, 'capitalize')}>
          <span className='truncate'>
            {themeLabels[props.value] ? t(themeLabels[props.value]) : props.value}
          </span>
        </SettingButton>
      </Dialog.Trigger>
      <Dialog.Content title='selectTheme' className='flex flex-col'>
        <div role='tablist' className='mt-3 grid grid-cols-2 border-b border-slate-50/10'>
          {themeLayouts.map(l => (
            <button
              key={l.key}
              type='button'
              role='tab'
              aria-selected={l.key === layout}
              onClick={() => setLayout(l.key)}
              className={cn(
                '-mb-px border-b-2 border-transparent py-1 transition-colors hover:bg-slate-50/5',
                l.key === layout && 'border-green-300 font-semibold'
              )}
            >
              {t(l.label)}
            </button>
          ))}
        </div>
        <ul role='tabpanel' className='my-4 min-h-0 w-full flex-1 overflow-y-auto'>
          {themes
            .filter(theme => layoutOf(theme.name) === layout)
            .map(theme => (
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
        <div className='relative mx-auto h-[60px] w-[350px] shrink-0'>
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
              <style>{theme.css.match(/@import.*;|@font-face\s*{[^}]*}/g)}</style>
            </div>
          ))}
        </div>
      </Dialog.Content>
    </Dialog.Root>
  )
}
