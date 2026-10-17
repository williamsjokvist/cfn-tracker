import React from 'react'
import { useTranslation } from 'react-i18next'
import { motion } from 'framer-motion'
import { ArrowRightIcon } from '@radix-ui/react-icons'

import { model } from '@model'
import { useSelector } from '@xstate/react'

import { AuthMachineContext } from '@/state/auth-machine'
import { ConfigContext } from '@/main/config'

import * as Page from '@/ui/page'
import { cn } from '@/helpers/cn'

import sf6 from './games/sf6.webp'
import t8 from './games/t8.svg'
import fightcade from './games/fightcade.svg'

type Game = {
  code: model.GameType
  name: string
  logo: string
  // The logos have different amounts of empty space around them. Cropping each to
  // a similar wide shape makes them show at about the same size.
  logoClassName: string
  // An RGB triplet that tints the game's card
  accent: string
}

const GAMES: Game[] = [
  {
    code: model.GameType.STREET_FIGHTER_6,
    name: 'Street Fighter 6',
    logo: sf6,
    logoClassName: 'aspect-[1280/200] object-cover invert',
    accent: '236 72 153'
  },
  {
    code: model.GameType.TEKKEN_8,
    name: 'Tekken 8',
    logo: t8,
    logoClassName: 'aspect-[2498/345]',
    accent: '239 68 68'
  },
  {
    code: model.GameType.FIGHTCADE,
    name: 'Fightcade',
    logo: fightcade,
    logoClassName: 'aspect-[290/40]',
    accent: '51 241 229'
  }
]

export const gameName = (game?: model.GameType) => GAMES.find(g => g.code === game)?.name ?? ''

const accentStyle = (game: Game) => ({ '--accent': game.accent }) as React.CSSProperties

export function TrackingGamePicker(props: { onSubmit: (game: model.GameType) => void }) {
  const { t } = useTranslation()
  const [cfg] = React.useContext(ConfigContext)
  const authActor = AuthMachineContext.useActorRef()
  // The game in the auth machine is newer than the one the config had at startup.
  const currentGame = useSelector(authActor, ({ context }) => context.game)
  const lastGame = GAMES.find(g => g.code === (currentGame ?? cfg.lastGame))
  const otherGames = GAMES.filter(g => g !== lastGame)

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('pickGame')}</Page.Title>
      </Page.Header>
      <motion.div
        initial={{ opacity: 0, y: 8 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className='mx-auto flex w-full max-w-3xl flex-col gap-8 self-center px-8'
      >
        {lastGame && (
          <button
            type='button'
            style={accentStyle(lastGame)}
            onClick={() => props.onSubmit(lastGame.code)}
            className={cn(
              'group flex flex-col items-center gap-8 rounded-2xl px-10 pt-12 pb-6',
              'border border-white/10 transition-colors hover:border-white/25',
              'bg-[radial-gradient(ellipse_at_top,rgb(var(--accent)/0.18),rgb(255_255_255/0.05)_70%)]'
            )}
          >
            <img
              src={lastGame.logo}
              alt={lastGame.name}
              className={cn(
                'pointer-events-none w-3/5 transition-transform select-none group-hover:scale-[1.03]',
                lastGame.logoClassName
              )}
            />
            <span className='flex items-center gap-2 text-lg text-gray-300 transition-colors group-hover:text-white'>
              {t('continueWithGame', { game: lastGame.name })}
              <ArrowRightIcon className='h-5 w-5 transition-transform group-hover:translate-x-1' />
            </span>
          </button>
        )}
        <section className='flex flex-col gap-3'>
          {lastGame && (
            <h2 className='text-sm tracking-widest text-gray-400 uppercase'>{t('otherGames')}</h2>
          )}
          {/* Fills rows as games are added, rather than squeezing them onto one */}
          <ul className='grid grid-cols-[repeat(auto-fill,minmax(12rem,1fr))] gap-4'>
            {otherGames.map(game => (
              <li key={game.code}>
                <button
                  type='button'
                  style={accentStyle(game)}
                  onClick={() => props.onSubmit(game.code)}
                  aria-label={game.name}
                  className={cn(
                    'group flex h-24 w-full items-center justify-center rounded-xl px-6',
                    'border border-white/10 bg-white/5 transition-colors',
                    'hover:border-[rgb(var(--accent)/0.5)] hover:bg-[rgb(var(--accent)/0.08)]'
                  )}
                >
                  <img
                    src={game.logo}
                    alt=''
                    className={cn(
                      'pointer-events-none w-full transition-transform select-none group-hover:scale-105',
                      game.logoClassName
                    )}
                  />
                </button>
              </li>
            ))}
          </ul>
        </section>
      </motion.div>
    </Page.Root>
  )
}
