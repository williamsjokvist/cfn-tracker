import React from 'react'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'


import { GetSessions, GetSessionsStatistics } from '@cmd/CommandHandler'
import type { model } from '@model'


import { useErrorPopup } from '@/main/error-popup'
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/ui/hover-card'
import * as Page from '@/ui/page'
import { Button } from '@/ui/button'

export function SessionsListPage() {
  const { i18n, t } = useTranslation()
  const navigate = useNavigate()
  const setError = useErrorPopup()

  const [sessions, setSessions] = React.useState<model.Session[]>([])
  const [sessionStatistics, setSessionStatistics] = React.useState<model.SessionsStatistics>()
  const [year, setYear] = React.useState("")
  const [month, setMonth] = React.useState("01")
  const [monthIndex, setMonthIndex] = React.useState(0)

  const months = sessionStatistics?.Months ?? []

  React.useEffect(() => {
    GetSessionsStatistics('').then(setSessionStatistics).catch(setError)
  }, [])

  React.useEffect(() => {
    if (months.length > 0 && months[monthIndex]) {
      const [month, year] = months[monthIndex].Date.split('-')
      setMonth(month)
      setYear(year)
    }
  }, [sessionStatistics, monthIndex])

  React.useEffect(() => {
    GetSessions("", month, 0, 0).then(setSessions).catch(setError)
  }, [month])

  const sessionsByDay = (sessions ?? []).reduce((group, session) => {
    const date = new Date(session.createdAt)
    const day = date.getDate()
    group[day] = group[day] ?? []
    group[day].push(session)
    return group
  }, {} as Record<string, model.Session[]>)

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
        <header className='px-8 py-4 text-xl flex gap-2 items-center'>
          <Button
            className="!py-0 !px-0 !text-md !font-normal"
            disabled={months[monthIndex + 1] === undefined}
            onClick={() => setMonthIndex(monthIndex + 1)}>
            <Icon width={26} height={26} icon='material-symbols:chevron-left' />
          </Button>
          <Button
            className="!py-0 !px-0 !text-md !font-normal"
            disabled={monthIndex === 0}
            onClick={() => setMonthIndex(monthIndex - 1)}>
            <Icon width={26} height={26} icon='material-symbols:chevron-left' className='rotate-180' />
          </Button>
          <h2 className='ml-2 font-bold'>
            {year}{" "}/{" "}
            {Intl.DateTimeFormat(i18n.resolvedLanguage, {
              month: 'long'
            }).format(new Date(`2024-${month}-01`))}
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
          className='relative flex flex-wrap items-stretch border-y-[0.5px] border-solid border-divider px-8'
        >
          {Object.keys(sessionsByDay).map(day => (
            <div
              key={day}
              className='flex w-[193.5px] flex-col border-b-[0.5px] border-solid border-divider px-2'
            >
              <span className='text-center text-xl font-bold'>{day}</span>
              {sessionsByDay[day].reverse().map(s => (
                <HoverCard key={s.id} openDelay={250}>
                  <HoverCardTrigger>
                    <Button
                      className='mb-1 w-full !justify-between gap-2 rounded-xl !px-[6px] !py-0 !pt-[2px] text-xl'
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
                  </HoverCardTrigger>
                  <HoverCardContent side='bottom'>
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
                  </HoverCardContent>
                </HoverCard>
              ))}
            </div>
          ))}
        </div>
      </motion.div>
    </Page.Root>
  )
}
