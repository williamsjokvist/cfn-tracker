import React from 'react'
import { useTranslation } from 'react-i18next'
import { motion } from 'framer-motion'

import { model } from '@model'

import { Button } from '@/ui/button'
import * as Page from '@/ui/page'
import { cn } from '@/helpers/cn'

import sf6 from './games/sf6.webp'
import t8 from './games/t8.png'

const GAMES = [
  {
    logo: t8,
    code: model.GameType.TEKKEN_8,
    alt: 'Tekken 8',
    style: {
      position: 'absolute',
      left: 0,
      top: -28,
      padding: '0 12px'
    }
  },
  {
    logo: sf6,
    code: model.GameType.STREET_FIGHTER_6,
    alt: 'Street Fighter 6',
    style: { filter: 'invert(1)' }
  }
] as const

export function TrackingGamePicker(props: { onSubmit: (game: model.GameType) => void }) {
  const { t } = useTranslation()
  const [selectedGame, setSelectedGame] = React.useState<model.GameType>()

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('pickGame')}</Page.Title>
      </Page.Header>
      <div className='flex flex-col items-center justify-center gap-10 justify-self-center'>
        <motion.ul
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className='flex w-full justify-center gap-24'
        >
          {GAMES.map(game => (
            <li key={game.code}>
              <button
                type='button'
                className={cn(
                  'relative h-[70px] w-60 rounded-2xl px-3',
                  'transition-colors hover:bg-slate-50/5'
                )}
                {...(game.code === selectedGame && {
                  style: {
                    outline: '1px solid lightblue',
                    background: 'rgb(248 250 252 / 0.05)'
                  }
                })}
                onClick={() => setSelectedGame(game.code)}
              >
                <img
                  src={game.logo}
                  alt={game.alt}
                  style={game.style}
                  className='pointer-events-none select-none'
                />
              </button>
            </li>
          ))}
        </motion.ul>
        <Button
          onClick={() => {
            selectedGame && props.onSubmit(selectedGame)
          }}
          disabled={!selectedGame}
        >
          {t('continueStep')}
        </Button>
      </div>
    </Page.Root>
  )
}
