import React from 'react'
import { useTranslation } from 'react-i18next'
import { useLoaderData, useNavigate } from 'react-router-dom'

import * as Page from '@/ui/page'
import * as Table from '@/ui/table'
import { motion } from 'framer-motion'

import type { model } from '@model'
import { Button } from '@/ui/button'
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/ui/hover-card'

type DayGroup = Record<string, model.Session[]>
type MonthGroup = Record<string, DayGroup>
type YearGroup = Record<string, MonthGroup>

export function SessionsListPage() {
  const sessions = (useLoaderData() ?? []) as model.Session[]
  const { i18n, t } = useTranslation()
  const navigate = useNavigate()

  const groupedSessions: YearGroup = sessions.reduce((group, sesh) => {
    const date = new Date(sesh.createdAt)
    const year = date.getFullYear()
    const month = date.getMonth() + 1
    const day = date.getDate()

    group[year] = group[year] ?? {}
    group[year][month] = group[year][month] ?? []
    group[year][month][day] = group[year][month][day] ?? []
    group[year][month][day].push(sesh)

    return group
  }, {})

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
        {Object.keys(groupedSessions)
          .reverse()
          .map(year => (
            <section key={year} className=''>
              <h2 className='px-8 py-6 text-4xl font-bold'>{year}</h2>
              {Object.keys(groupedSessions[year])
                .reverse()
                .map(month => (
                  <div key={month}>
                    <h3 className='px-8 py-4 text-xl font-bold'>
                      {Intl.DateTimeFormat(i18n.resolvedLanguage, {
                        month: 'long'
                      }).format(new Date(`2024-${Number(month) < 10 ? '0' + month : month}-01`))}
                    </h3>
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
                      {Object.keys(groupedSessions[year][month]).map(day => (
                        <div className='flex w-[193.5px] flex-col border-b-[0.5px] border-solid border-divider px-2'>
                          <span className='text-center text-xl font-bold'>{day}</span>
                          {groupedSessions[year][month][day].reverse().map(s => (
                            <HoverCard>
                              <HoverCardTrigger>
                                <Button
                                  key={day}
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
                  </div>
                ))}
            </section>
          ))}
      </motion.div>
    </Page.Root>
  )
}
