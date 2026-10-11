import React from 'react'
import { useLoaderData, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'

import type { model } from '@model'

import * as Page from '@/ui/page'
import { Button } from '@/ui/button'
import { BigStat } from '@/ui/stat'
import { CharacterPortrait } from '@/ui/character-portrait'

type CharacterStats = {
  character: string
  wins: number
  losses: number
  lpGain: number
  mrGain: number
}

// Matches come newest first, and a character's latest match carries its gains for the session.
function statsByCharacter(matches: model.Match[]): CharacterStats[] {
  const stats = new Map<string, CharacterStats>()
  for (const match of matches) {
    let character = stats.get(match.character)
    if (!character) {
      character = {
        character: match.character,
        wins: 0,
        losses: 0,
        lpGain: match.lpGain,
        mrGain: match.mrGain
      }
      stats.set(match.character, character)
    }
    if (match.victory) {
      character.wins++
    } else {
      character.losses++
    }
  }
  return [...stats.values()].sort((a, b) => b.wins + b.losses - (a.wins + a.losses))
}

const signed = (value: number) => (value > 0 ? `+${value}` : value)

const winRate = (wins: number, losses: number) =>
  wins + losses === 0 ? 0 : Math.floor((wins / (wins + losses)) * 100)

export type SessionData = { session: model.Session; matches: model.Match[] }

export function SessionBreadcrumbs({
  session,
  matchLog
}: {
  session: model.Session
  matchLog?: boolean
}) {
  const { t, i18n } = useTranslation()
  const started = Intl.DateTimeFormat(i18n.resolvedLanguage, {
    dateStyle: 'medium',
    timeStyle: 'short'
  }).format(new Date(session.createdAt))
  return (
    <Page.Breadcrumbs
      items={[
        { label: t('sessions'), to: '/sessions' },
        { label: started, to: matchLog ? `/sessions/${session.id}` : undefined },
        ...(matchLog ? [{ label: t('history') }] : [])
      ]}
    />
  )
}

export function SessionOverviewPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { session, matches } = useLoaderData() as SessionData

  const characters = statsByCharacter(matches)
  const wins = characters.reduce((sum, c) => sum + c.wins, 0)
  const losses = characters.reduce((sum, c) => sum + c.losses, 0)
  // Tekken has no LP or MR.
  const hasPoints = matches.some(m => m.lp !== 0 || m.mr !== 0)

  return (
    <Page.Root>
      <Page.Header>
        <SessionBreadcrumbs session={session} />
      </Page.Header>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='relative w-full'
      >
        <div className='mb-2 flex h-10 items-center justify-between border-b border-slate-50/10 px-8 pt-1'>
          <Button
            className='flex items-center gap-1 px-1 py-[2px]'
            style={{ filter: 'hue-rotate(-65deg)' }}
            onClick={() => navigate('/sessions')}
          >
            <Icon
              width={20}
              height={20}
              icon='material-symbols:chevron-left'
              className='rtl:rotate-180'
            />
            <span className='relative top-[2px]'>{t('goBack')}</span>
          </Button>
          <Button
            className='flex items-center gap-2 px-2 py-[2px]'
            disabled={matches.length === 0}
            onClick={() => navigate(`/sessions/${session.id}/matches`)}
          >
            <Icon width={18} height={18} icon='material-symbols:list-alt-outline' />
            <span className='relative top-[2px]'>{t('history')}</span>
          </Button>
        </div>

        <div className='flex max-h-[340px] flex-col gap-4 overflow-y-auto px-8 pt-2 pb-4'>
          {characters.length > 1 && (
            <dl className='flex gap-2'>
              <BigStat text={t('wins')} value={wins} />
              <BigStat text={t('losses')} value={losses} />
              <BigStat text={t('winRate')} value={`${winRate(wins, losses)}%`} />
            </dl>
          )}
          {characters.map(c => (
            <CharacterSection key={c.character} stats={c} hasPoints={hasPoints} />
          ))}
        </div>
      </motion.div>
    </Page.Root>
  )
}

function CharacterSection({ stats, hasPoints }: { stats: CharacterStats; hasPoints: boolean }) {
  const { t } = useTranslation()
  // The portrait is a square as tall as the stats, which CSS can't size from a sibling.
  const statsRef = React.useRef<HTMLDivElement>(null)
  const [size, setSize] = React.useState(0)
  React.useLayoutEffect(() => {
    const observer = new ResizeObserver(([entry]) => setSize(entry.borderBoxSize[0].blockSize))
    observer.observe(statsRef.current!)
    return () => observer.disconnect()
  }, [])

  return (
    <section className='flex items-start gap-4'>
      <CharacterPortrait
        character={stats.character}
        className='relative shrink-0 overflow-hidden rounded-xl bg-slate-50/5'
        style={{ width: size, height: size }}
      />
      <div ref={statsRef} className='flex-1'>
        <h3 className='mb-2 text-lg font-bold'>{stats.character}</h3>
        <dl className='grid grid-cols-3 gap-2 [&>div]:mb-0'>
          <BigStat text={t('wins')} value={stats.wins} />
          <BigStat text={t('losses')} value={stats.losses} />
          <BigStat text={t('winRate')} value={`${winRate(stats.wins, stats.losses)}%`} />
          {hasPoints && (
            <>
              <BigStat text={t('lpGain')} value={signed(stats.lpGain)} />
              <BigStat text={t('mrGain')} value={signed(stats.mrGain)} />
            </>
          )}
        </dl>
      </div>
    </section>
  )
}
