import React from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useSelector } from '@xstate/react'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'

import { ExportUserMatches, GetSessions, GetUsers, GetUserStatistics } from '@cmd/CommandHandler'
import { model } from '@model'

import { AuthMachineContext } from '@/state/auth-machine'
import { TrackingMachineContext } from '@/state/tracking-machine'
import { useErrorPopup } from '@/main/error-popup'
import * as Page from '@/ui/page'
import * as Select from '@/ui/select'
import { Button } from '@/ui/button'
import { BigStat } from '@/ui/stat'
import { CharacterStatsCard } from '@/ui/character-stats'
import { cn } from '@/helpers/cn'
import { signed, winRate } from '@/helpers/stats'

const recentSessionCount = 5

export function StatisticsPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const setError = useErrorPopup()
  const [searchParams, setSearchParams] = useSearchParams()

  const authActor = AuthMachineContext.useActorRef()
  const game = useSelector(authActor, ({ context }) => context.game)
  const trackingActor = TrackingMachineContext.useActorRef()
  const trackedCode = useSelector(trackingActor, ({ context }) => context.user?.code)

  const [users, setUsers] = React.useState<model.User[]>()
  React.useEffect(() => {
    // Users are saved per game, so before a game is picked there's no reason to hide any.
    const games = game ? [game] : Object.values(model.GameType)
    Promise.all(games.map(g => GetUsers(g)))
      .then(results => setUsers(results.flatMap(users => users ?? [])))
      .catch(setError)
  }, [game])

  const user =
    users?.find(u => u.code === searchParams.get('user')) ??
    users?.find(u => u.code === trackedCode) ??
    users?.[0]

  const [stats, setStats] = React.useState<model.UserStatistics>()
  const [sessions, setSessions] = React.useState<model.Session[]>([])
  React.useEffect(() => {
    if (!user) return
    setStats(undefined)
    setSessions([])
    GetUserStatistics(user.code).then(setStats).catch(setError)
    GetSessions(user.code, '', recentSessionCount, 0)
      .then(sessions => setSessions(sessions ?? []))
      .catch(setError)
  }, [user?.code])

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('statistics')}</Page.Title>
      </Page.Header>

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='relative flex min-h-0 w-full flex-col'
      >
        <div className='flex min-h-14 items-center justify-between gap-4 border-b border-slate-50/10 px-8 py-2'>
          {users && users.length > 1 ? (
            <Select.Root
              size='sm'
              value={user?.code}
              onValueChange={code => setSearchParams({ user: code }, { replace: true })}
            >
              {users.map(u => (
                <Select.Item key={u.code} value={u.code}>
                  {u.displayName}
                </Select.Item>
              ))}
            </Select.Root>
          ) : (
            <h3 className='text-lg font-bold'>{user?.displayName}</h3>
          )}
          <div className='flex items-center gap-2'>
            {user?.game === model.GameType.STREET_FIGHTER_6 && stats?.mainCharacter && (
              <Button
                className='flex items-center gap-2 px-2 py-[2px]'
                onClick={() =>
                  navigate(`/statistics/comparison?${new URLSearchParams({ user: user.code })}`)
                }
              >
                <Icon width={18} height={18} icon='mdi:compare-horizontal' />
                <span className='relative top-[2px]'>
                  {t('compareWithTopPlayers', { character: stats.mainCharacter })}
                </span>
              </Button>
            )}
            <Button
              className='flex items-center gap-2 px-2 py-[2px]'
              onClick={() => navigate('/sessions')}
            >
              <Icon width={18} height={18} icon='ion:document-text-outline' />
              <span className='relative top-[2px]'>{t('sessions')}</span>
            </Button>
            {user && stats && stats.sessions > 0 && (
              <Button
                className='flex items-center gap-2 px-2 py-[2px]'
                style={{ filter: 'hue-rotate(145deg)' }}
                onClick={() => ExportUserMatches(user.code).catch(setError)}
              >
                <Icon width={18} height={18} icon='material-symbols:download' />
                <span className='relative top-[2px]'>{t('export')}</span>
              </Button>
            )}
          </div>
        </div>

        {users?.length === 0 && <p className='px-8 py-6 text-lg'>{t('noStatistics')}</p>}

        {stats && (
          <div className='flex flex-col gap-4 overflow-y-auto px-8 pt-4 pb-4'>
            <dl className='grid grid-cols-4 gap-2 [&>div]:mb-0'>
              <BigStat text={t('sessions')} value={stats.sessions} />
              <BigStat text={t('wins')} value={stats.wins} />
              <BigStat text={t('losses')} value={stats.losses} />
              <BigStat text={t('winRate')} value={`${winRate(stats.wins, stats.losses)}%`} />
            </dl>
            <CharacterList characters={stats.characters} />
            {sessions.length > 0 && (
              <RecentSessions
                sessions={sessions}
                hasPoints={user?.game === model.GameType.STREET_FIGHTER_6}
              />
            )}
          </div>
        )}
      </motion.div>
    </Page.Root>
  )
}

function RecentSessions({
  sessions,
  hasPoints
}: {
  sessions: model.Session[]
  hasPoints: boolean
}) {
  const { t, i18n } = useTranslation()
  const navigate = useNavigate()
  const format = Intl.DateTimeFormat(i18n.resolvedLanguage, {
    dateStyle: 'medium',
    timeStyle: 'short'
  })
  return (
    <section>
      <h3 className='mb-2 text-lg font-bold'>{t('recentSessions')}</h3>
      <ul className='flex flex-col gap-1'>
        {sessions.map(s => (
          <li key={s.id}>
            <button
              type='button'
              onClick={() =>
                navigate(`/sessions/${s.id}`, {
                  state: { backTo: `/statistics?${new URLSearchParams({ user: s.userId })}` }
                })
              }
              className={cn(
                'grid w-full items-center gap-6 rounded-xl bg-slate-50/5 px-4 py-2 text-start tabular-nums transition-colors hover:bg-slate-50/10',
                hasPoints ? 'grid-cols-[1fr_auto_3rem_6rem_5rem]' : 'grid-cols-[1fr_auto_3rem]'
              )}
            >
              <span>{format.format(new Date(s.createdAt))}</span>
              <span>
                <span className='text-green-300'>{s.matchesWon}</span>
                {' – '}
                <span className='text-red-300'>{s.matchesLost}</span>
              </span>
              <span className='text-end'>{winRate(s.matchesWon, s.matchesLost)}%</span>
              {hasPoints && (
                <>
                  <Gain value={s.lpGain} unit='LP' />
                  <Gain value={s.mrGain} unit='MR' />
                </>
              )}
            </button>
          </li>
        ))}
      </ul>
    </section>
  )
}

function Gain({ value, unit }: { value: number; unit: string }) {
  return (
    <span className={cn('text-end', value > 0 && 'text-green-300', value < 0 && 'text-red-300')}>
      {signed(value)} {unit}
    </span>
  )
}

function CharacterList({ characters }: { characters: model.CharacterStatistics[] }) {
  const { t } = useTranslation()
  // Tekken has no LP or MR.
  const hasPoints = characters.some(c => c.lp !== 0 || c.mr !== 0)
  return characters.map(c => (
    <CharacterStatsCard
      key={c.character}
      character={c.character}
      heading={<h3 className='mb-2 text-lg font-bold'>{c.character}</h3>}
    >
      <BigStat text={t('wins')} value={c.wins} />
      <BigStat text={t('losses')} value={c.losses} />
      <BigStat text={t('winRate')} value={`${winRate(c.wins, c.losses)}%`} />
      {/* Tekken has no points, but keeps the row so every character is the same size. */}
      <BigStat text='LP' value={c.lp} className={cn(!hasPoints && 'invisible')} />
      <BigStat text='MR' value={c.mr} className={cn(!hasPoints && 'invisible')} />
    </CharacterStatsCard>
  ))
}
