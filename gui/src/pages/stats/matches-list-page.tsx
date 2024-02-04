import React from 'react'
import { useLoaderData } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

import * as Page from '@/ui/page'
import * as Table from '@/ui/table'

import type { model } from '@@/go/models'

export function MatchesListPage() {
  const { t } = useTranslation()
  const totalMatches = (useLoaderData() ?? []) as model.Match[]
  const [matches, setMatches] = React.useState(totalMatches)
  const [totalWinRate, setTotalWinRate] = React.useState(0)

  React.useEffect(() => {
    const wonMatches = matches.filter(log => log.victory === true).length
    const winRate = Math.floor((wonMatches / matches.length) * 100)
    !isNaN(winRate) && setTotalWinRate(winRate)
  }, [matches])

  const filterLog = (property: string, value: string) => {
    setMatches(matches.filter(ml => ml[property].toLowerCase() === value.toLowerCase()))
  }

  if (matches.length === 0) {
    return (
      <Page.Root>
        <Page.Header>
          <Page.Title>{t('history')}</Page.Title>
        </Page.Header>
      </Page.Root>
    )
  }

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('history')}</Page.Title>
      </Page.Header>

      <div className='relative w-full'>
        <div className='mb-2 flex h-10 items-center border-b border-slate-50 border-opacity-10 px-8 pt-1 '>
          <span>
            {t('winRate')}: <b>{totalWinRate}</b>%
          </span>
        </div>
        <Table.Page>
          <Table.Content>
            <thead>
              <Table.Tr>
                <Table.Th className='w-[120px]'>{t('date')}</Table.Th>
                <Table.Th className='w-[70px]'>{t('time')}</Table.Th>
                <Table.Th className='w-[180px]'>{t('opponent')}</Table.Th>
                <Table.Th>{t('league')}</Table.Th>
                <Table.Th className='text-center'>{t('character')}</Table.Th>
                <Table.Th className='text-center'>{t('result')}</Table.Th>
              </Table.Tr>
            </thead>
            <tbody>
              {matches.map(log => (
                <Table.Tr key={`${log.date}${log.time}`}>
                  <Table.Td interactive onClick={() => filterLog('date', log.date)}>
                    {log.date}
                  </Table.Td>
                  <Table.Td>{log.time}</Table.Td>
                  <Table.Td interactive onClick={() => filterLog('opponent', log.opponent)}>
                    {log.opponent}
                  </Table.Td>
                  <Table.Td interactive onClick={() => filterLog('opponentLeague', log.opponentLeague)}>
                    {log.opponentLeague}
                  </Table.Td>
                  <Table.Td interactive className='text-center' onClick={() => filterLog('opponentCharacter', log.opponentCharacter)}>
                    {log.opponentCharacter}
                  </Table.Td>
                  <Table.Td
                    className='text-center'
                    style={{ color: log.victory === true ? 'lime' : 'red' }}
                  >
                    {log.victory === true ? 'W' : 'L'}
                  </Table.Td>
                </Table.Tr>
              ))}
            </tbody>
          </Table.Content>
        </Table.Page>
      </div>
    </Page.Root>
  )
}
