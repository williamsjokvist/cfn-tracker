import React from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'

import { GetSessions, GetSessionsStatistics } from '@cmd/CommandHandler'
import type { model } from '@model'

import { useErrorPopup } from '@/main/error-popup'
import * as HoverCard from '@/ui/hover-card'
import * as Page from '@/ui/page'
import { Button } from '@/ui/button'
import { cn } from '@/helpers/cn'

export function SessionsListPage() {
  const { i18n, t } = useTranslation()
  const navigate = useNavigate()
  const setError = useErrorPopup()

  const [sessions, setSessions] = React.useState<model.Session[]>([])
  const [sessionStatistics, setSessionStatistics] = React.useState<model.SessionsStatistics>()
  const [monthIndex, setMonthIndex] = React.useState(0)

  const months = sessionStatistics?.Months ?? []

  const sessionsByDay = (sessions ?? []).reduce(
    (group, session) => {
      const date = new Date(session.createdAt)
      const day = date.getDate()
      group[day] = group[day] ?? []
      group[day].push(session)
      return group
    },
    {} as Record<string, model.Session[]>
  )

  React.useEffect(() => {
    GetSessionsStatistics('').then(setSessionStatistics).catch(setError)
  }, [])

  React.useEffect(() => {
    if (months.length > 0 && months[monthIndex]) {
      GetSessions('', months[monthIndex].Date, 0, 0).then(setSessions).catch(setError)
    }
  }, [monthIndex, months])

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('sessions')}</Page.Title>
      </Page.Header>
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='overflow-y-scroll'
      >
        <header className='flex items-center gap-2 px-8 py-4 text-xl'>
          <Button
            className='text-md! px-0! py-0! font-normal!'
            disabled={months[monthIndex + 1] === undefined}
            onClick={() => setMonthIndex(monthIndex + 1)}
          >
            <Icon width={26} height={26} icon='material-symbols:chevron-left' />
          </Button>
          <Button
            className='text-md! px-0! py-0! font-normal!'
            disabled={monthIndex === 0}
            onClick={() => setMonthIndex(monthIndex - 1)}
          >
            <Icon
              width={26}
              height={26}
              icon='material-symbols:chevron-left'
              className='rotate-180'
            />
          </Button>
          <h2 className='ml-2 font-bold'>
            {months.length > 0 && months[monthIndex] && (
              <>
                {months[monthIndex].Date.split('-')[0]} /{' '}
                {Intl.DateTimeFormat(i18n.resolvedLanguage, {
                  month: 'long'
                }).format(new Date(`2024-${months[monthIndex].Date.split('-')[1]}-01`))}
              </>
            )}
          </h2>
        </header>

        <div
          style={{
            background: `repeating-linear-gradient(
            90deg,
            transparent 31.5px,
            transparent 224px,
            rgba(255, 255, 255, 0.125) 225px
          )`
          }}
          className='border-divider relative flex flex-wrap items-stretch border-y-[0.5px] border-solid px-8'
        >
          {Object.keys(sessionsByDay).map(day => (
            <div
              key={day}
              className='border-divider flex w-[193.5px] flex-col border-b-[0.5px] border-solid px-2'
            >
              <span className='text-center text-xl font-bold'>{day}</span>
              {sessionsByDay[day].reverse().map(s => (
                <HoverCard.Root key={s.id} openDelay={250}>
                  <HoverCard.Trigger>
                    <Button
                      className='mb-1 w-full justify-between! gap-2 rounded-xl px-[6px]! py-0! pt-[2px]! text-xl'
                      onClick={() => navigate(`/sessions/${s.id}/matches`)}
                    >
                      <span className='text-base font-bold'>
                        {Intl.DateTimeFormat(i18n.resolvedLanguage, {
                          hour: '2-digit',
                          minute: '2-digit'
                        }).format(new Date(s.createdAt))}
                      </span>
                      <span className='text-base font-light'>{s.userName}</span>
                    </Button>
                  </HoverCard.Trigger>
                  <HoverCard.Content
                    side='bottom'
                    className={cn(
                      'overflow-hidden p-4 text-white',
                      'bg-black/90 backdrop-blur-xl',
                      'rounded-xl shadow-[0_3px_16px_rgba(0,0,0,.5)]'
                    )}
                  >
                    <dl>
                      <div className='flex justify-between gap-2'>
                        <dt>{t('wins')}</dt>
                        <dd>{s.matchesWon}</dd>
                      </div>
                      <div className='flex justify-between gap-2'>
                        <dt>{t('losses')}</dt>
                        <dd>{s.matchesLost}</dd>
                      </div>
                      {s.lpGain != 0 && s.mrGain != 0 && (
                        <>
                          <div className='flex justify-between gap-2'>
                            <dt>{t('mrGain')}</dt>
                            <dd>{s.mrGain}</dd>
                          </div>
                          <div className='flex justify-between gap-2'>
                            <dt>{t('lpGain')}</dt>
                            <dd>{s.lpGain}</dd>
                          </div>
                        </>
                      )}
                    </dl>
                  </HoverCard.Content>
                </HoverCard.Root>
              ))}
            </div>
          ))}
        </div>
      </motion.div>
    </Page.Root>
  )
}
