import React from 'react'
import { Icon } from '@iconify/react'
import { Trans, useTranslation } from 'react-i18next'
import { motion, useAnimate } from 'framer-motion'

import * as Page from '@/ui/page'
import { Button } from '@/ui/button'

import { model } from '@@/go/models'
import { OpenResultsDirectory } from '@@/go/core/CommandHandler'

import { ThemeSelect } from './theme-select'
import { StatSelect } from './stat-select'

export type StatOptions = Omit<
  Record<keyof model.TrackingState, boolean>,
  'totalLosses' | 'totalWins'
> & {
  theme: string
}

const defaultOptions: StatOptions = {
  theme: 'default',
  cfn: false,
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
  opponentLP: false,
  totalMatches: false,
  character: false,
  result: false,
  userCode: false,
  timestamp: false,
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
          className='grid w-full gap-4 border-b-[1px] border-solid border-b-[rgba(255,255,255,.125)] px-8 py-6'
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
                  className='absolute right-7 top-0 opacity-0'
                />
              </div>
              <Button onClick={copyUrlToClipBoard} style={{ filter: 'hue-rotate(-120deg)' }}>
                <Icon icon='mdi:internet' className='mr-3 h-6 w-6' />
                {t('copyBrowserSourceLink')}
              </Button>
            </div>
            <div className='flex flex-col gap-4'>
              <ThemeSelect
                selectedTheme={options.theme}
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
