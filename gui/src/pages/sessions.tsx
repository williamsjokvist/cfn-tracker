import React from 'react'
import { useTranslation } from 'react-i18next'
import { useLoaderData, useNavigate } from 'react-router-dom'

import * as Page from '@/ui/page'
import * as Table from '@/ui/table'

import type { model } from '@model'

type MonthGroup = Record<string, model.Session[]>
type YearGroup = Record<string, MonthGroup>

export function SessionsListPage() {
  const sessions = (useLoaderData() ?? []) as model.Session[]
  const { i18n, t } = useTranslation()
  const navigate = useNavigate()

  const groupedSessions: YearGroup = sessions.reduce((group, sesh) => {
    const date = new Date(sesh.createdAt)
    const year = date.getFullYear()
    const month = date.getMonth() + 1

    group[year] = group[year] ?? {}
    group[year][month] = group[year][month] ?? []
    group[year][month].push(sesh)

    return group
  }, {})

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('sessions')}</Page.Title>
      </Page.Header>
      <Table.Page>
        {Object.keys(groupedSessions)
          .reverse()
          .map(year => (
            <section key={year}>
              <h2 className='mt-2 text-4xl'>{year}</h2>
              {Object.keys(groupedSessions[year])
                .reverse()
                .map(month => (
                  <section key={month}>
                    <h3 className='mt-2 text-2xl'>
                      {Intl.DateTimeFormat(i18n.resolvedLanguage, {
                        month: 'long'
                      }).format(new Date(`2024-${Number(month) < 10 ? '0' + month : month}-01`))}
                    </h3>
                    <Table.Content>
                      <thead>
                        <Table.Tr>
                          <Table.Th className='w-[120px]'>{t('started')}</Table.Th>
                          <Table.Th className='w-full'>{t('user')}</Table.Th>
                          <Table.Th>{t('mrGain')}</Table.Th>
                          <Table.Th>{t('lpGain')}</Table.Th>
                          <Table.Th className='text-center'>{t('matchesWon')}</Table.Th>
                          <Table.Th className='text-center'>{t('matchesLost')}</Table.Th>
                        </Table.Tr>
                      </thead>
                      <tbody>
                        {groupedSessions[year][month].map(sesh => (
                          <Table.Tr
                            key={sesh.id}
                            className='group cursor-pointer'
                            onClick={() => navigate(`/sessions/${sesh.id}/matches`)}
                          >
                            <Table.Td>
                              <time dateTime={sesh.createdAt}>
                                {new Date(sesh.createdAt).toLocaleDateString(
                                  i18n.resolvedLanguage,
                                  {
                                    day: 'numeric',
                                    weekday: 'short',
                                    hour: '2-digit',
                                    minute: '2-digit'
                                  }
                                )}
                              </time>
                            </Table.Td>
                            <Table.Td>{sesh.userName}</Table.Td>
                            <Table.Td className='text-center'>
                              {`${sesh.mrGain > 0 ? `+` : ``}${sesh.mrGain}`}
                            </Table.Td>
                            <Table.Td className='text-center'>
                              {`${sesh.lpGain > 0 ? `+` : ``}${sesh.lpGain}`}
                            </Table.Td>
                            <Table.Td className='text-center'>{sesh.matchesWon}</Table.Td>
                            <Table.Td className='text-center'>{sesh.matchesLost}</Table.Td>
                          </Table.Tr>
                        ))}
                      </tbody>
                    </Table.Content>
                  </section>
                ))}
            </section>
          ))}
      </Table.Page>
    </Page.Root>
  )
}
