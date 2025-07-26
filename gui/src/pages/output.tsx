import React from 'react'
import { createRoot } from 'react-dom/client'
import { Icon } from '@iconify/react'
import { Trans, useTranslation } from 'react-i18next'
import { motion, useAnimate } from 'framer-motion'

import * as Page from '@/ui/page'
import * as Dialog from '@/ui/dialog'
import { Button } from '@/ui/button'

import { Checkbox } from '@/ui/checkbox'

import { model } from '@model'
import { GetThemes, OpenResultsDirectory } from '@cmd/CommandHandler'

import { useErrorPopup } from '@/main/error-popup'
import { cn } from '@/helpers/cn'

type StatOptions = Omit<Record<keyof model.Match, boolean>, 'replayId' | 'sessionId'> & {
  theme: string
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
  const [scope, animate] = useAnimate()

  const { theme: _, ...statOptions } = options

  function copyUrlToClipBoard() {
    let url = `http://localhost:4242?theme=${options.theme}`
    for (const [key, value] of Object.entries(options)) {
      if (key === 'theme') {
        continue
      }
      url += `&${key}=${value}`
    }

    navigator.clipboard.writeText(url)

    animate('#ok', { opacity: [0, 1], y: [0, -10] }, { delay: 0.125 }).then(() => {
      animate('#ok', { opacity: [1, 0] }, { delay: 0.5 })
    })
  }

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('output')}</Page.Title>
      </Page.Header>

      <div className='z-40 h-full w-full overflow-y-scroll'>
        <motion.section
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className='border-b-divider grid w-full gap-4 border-b-[1px] border-solid px-8 py-6'
        >
          <div className='flex items-center justify-between gap-8'>
            <div>
              <h2 className='text-2xl font-bold'>{t('usingBrowserSource')}</h2>
              <p className='my-4 max-w-[340px]'>
                <Trans t={t} i18nKey='browserSourceDescription'>
                  The easiest way to display the stats in OBS is to use
                  <i className='whitespace-nowrap'>Browser Source</i>. After editing the options,
                  copy the browser source link.
                </Trans>
              </p>
              <div className='relative' ref={scope}>
                <Icon
                  icon='twemoji:ok-hand'
                  width={45}
                  id='ok'
                  className='absolute top-0 right-7 opacity-0'
                />
              </div>
              <Button onClick={copyUrlToClipBoard} style={{ filter: 'hue-rotate(-120deg)' }}>
                <Icon icon='mdi:internet' className='mr-3 h-6 w-6' />
                {t('copyBrowserSourceLink')}
              </Button>
            </div>
            <div className='flex flex-col gap-4'>
              <ThemeSelect
                value={options.theme}
                onSelect={theme => setOptions({ ...options, theme })}
              />
              <StatSelect
                options={statOptions}
                onSelect={(opt, checked) => setOptions({ ...options, [opt]: checked })}
              />
            </div>
          </div>
        </motion.section>
        <motion.section
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className='flex w-full items-center justify-between gap-8 px-8 py-6'
        >
          <header>
            <h2 className='text-xl font-bold'>{t('importFiles')}</h2>
            <p className='mt-2 max-w-[420px]'>
              <Trans t={t} i18nKey='obsCustomize'>
                If you want to customize through OBS, you can use
                <i>Text Labels</i> and add the text files in the results folder as sources
              </Trans>
            </p>
          </header>
          <Button
            onClick={OpenResultsDirectory}
            style={{ filter: 'hue-rotate(-120deg)' }}
            className='mx-auto'
          >
            <Icon icon='fa6-solid:folder-open' className='mr-3 h-6 w-6' />
            {t('files')}
          </Button>
        </motion.section>
      </div>
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
        <Button style={{ filter: 'hue-rotate(-45deg)' }}>
          <Icon icon='bx:stats' className='mr-3 h-6 w-6' />
          {t('displayStats')}
        </Button>
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
                <span className='ml-4 cursor-pointer text-center capitalize'>{opt}</span>
              </button>
            </li>
          ))}
        </ul>
      </Dialog.Content>
    </Dialog.Root>
  )
}

function ThemeSelect(props: { value: string; onSelect: (theme: string) => void }) {
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
        <Button
          className='capitalize'
          style={{ filter: 'hue-rotate(-180deg)', justifyContent: 'center' }}
        >
          <Icon icon='ph:paint-bucket-fill' className='mr-3 h-6 w-6' />
          {props.value}
        </Button>
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
                className='absolute top-1 left-2'
              />
              <label
                htmlFor={`${theme.name}-checkbox`}
                className='font-spartan w-full cursor-pointer py-1 pl-14 capitalize'
              >
                {theme.name}
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
                'absolute top-0 left-0 w-full transition-opacity',
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
