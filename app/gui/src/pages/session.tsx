import React from 'react'
import { useLoaderData, useLocation, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { motion } from 'framer-motion'

import type { model } from '@model'

import * as Page from '@/ui/page'
import { SettingButton } from '@/ui/setting-button'
import { BigStat } from '@/ui/stat'
import { CharacterStatsCard } from '@/ui/character-stats'
import { cn } from '@/helpers/cn'
import { signed, winRate } from '@/helpers/stats'

import { HistoryChart } from './statistics/history-chart'

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

export type SessionData = { session: model.Session; matches: model.Match[] }

// Where a session was opened from, carried between its overview and match log so
// going back returns there.
export type SessionLocationState = { backTo?: string } | null

export function useSessionBackTo() {
  const state = useLocation().state as SessionLocationState
  return { backTo: state?.backTo ?? '/sessions', state }
}

export function SessionBreadcrumbs({
  session,
  matchLog,
  character
}: {
  session: model.Session
  matchLog?: boolean
  character?: string
}) {
  const { t, i18n } = useTranslation()
  const { backTo, state } = useSessionBackTo()
  const started = Intl.DateTimeFormat(i18n.resolvedLanguage, {
    dateStyle: 'medium',
    timeStyle: 'short'
  }).format(new Date(session.createdAt))
  return (
    <Page.Breadcrumbs
      items={[
        backTo.startsWith('/statistics')
          ? { label: t('statistics'), to: backTo }
          : { label: t('sessions'), to: '/sessions' },
        { label: started, to: matchLog ? `/sessions/${session.id}` : undefined, state },
        ...(matchLog
          ? [{ label: character ? `${t('history')} · ${character}` : t('history') }]
          : [])
      ]}
    />
  )
}

export function SessionOverviewPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { state } = useSessionBackTo()
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
        <SettingButton
          icon='material-symbols:list-alt-outline'
          // Clickable although the header is the window's drag region.
          style={{ '--draggable': 'no-drag' } as React.CSSProperties}
          onClick={() => navigate(`/sessions/${session.id}/matches`, { state })}
        >
          {t('history')}
        </SettingButton>
      </Page.Header>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='relative flex min-h-0 w-full min-w-0 flex-col'
      >
        <div className='flex flex-col gap-2 overflow-y-auto px-8 pt-2 pb-2'>
          {characters.length > 1 && (
            <section className='-mx-8 border-b border-slate-50/10 px-8 pb-4'>
              <SectionHeading title={t('statistics')} />
              <dl className='grid grid-cols-3 gap-2 [&>div]:mb-0'>
                <BigStat text={t('wins')} value={wins} />
                <BigStat text={t('losses')} value={losses} />
                <BigStat text={t('winRate')} value={`${winRate(wins, losses)}%`} />
              </dl>
            </section>
          )}
          {characters.map(c => (
            <CharacterSection
              key={c.character}
              sessionId={session.id}
              stats={c}
              hasPoints={hasPoints}
              // With one character its match log is the whole session's, which the header opens.
              matchLog={characters.length > 1}
            />
          ))}
          {characters.length > 0 && (
            <HistoryChart
              history={matches.toReversed()}
              mainCharacter={characters[0].character}
              timeSpans={false}
            />
          )}
        </div>
      </motion.div>
    </Page.Root>
  )
}

function SectionHeading({ title, matchLog }: { title: string; matchLog?: string }) {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { state } = useSessionBackTo()
  return (
    <div className='mb-2 flex items-center justify-between gap-4'>
      <h3 className='text-lg font-bold'>{title}</h3>
      {matchLog && (
        <SettingButton
          icon='material-symbols:list-alt-outline'
          className='shrink-0'
          onClick={() => navigate(matchLog, { state })}
        >
          {t('history')}
        </SettingButton>
      )}
    </div>
  )
}

function CharacterSection({
  sessionId,
  stats,
  hasPoints,
  matchLog
}: {
  sessionId: number
  stats: CharacterStats
  hasPoints: boolean
  matchLog: boolean
}) {
  const { t } = useTranslation()
  return (
    <CharacterStatsCard
      character={stats.character}
      heading={
        <SectionHeading
          title={stats.character}
          matchLog={
            matchLog
              ? `/sessions/${sessionId}/matches?${new URLSearchParams({ character: stats.character })}`
              : undefined
          }
        />
      }
    >
      <BigStat text={t('wins')} value={stats.wins} />
      <BigStat text={t('losses')} value={stats.losses} />
      <BigStat text={t('winRate')} value={`${winRate(stats.wins, stats.losses)}%`} />
      {/* Tekken has no points, but keeps the row so every character is the same size. */}
      <BigStat
        text={t('lpGain')}
        value={signed(stats.lpGain)}
        className={cn(!hasPoints && 'invisible')}
      />
      <BigStat
        text={t('mrGain')}
        value={signed(stats.mrGain)}
        className={cn(!hasPoints && 'invisible')}
      />
    </CharacterStatsCard>
  )
}
