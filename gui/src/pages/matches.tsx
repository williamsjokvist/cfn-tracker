import React from 'react'
import { useLoaderData, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Icon } from '@iconify/react'

import * as Page from '@/ui/page'
import * as Table from '@/ui/table'

import type { model } from '@model'
import { Button } from '@/ui/button'
import { type LocalizationKey } from '@/main/i18n'
import { Tooltip } from '@/ui/tooltip'

export function MatchesListPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
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
        <div className='mb-2 flex h-10 items-center justify-between border-b border-slate-50 border-opacity-10 px-8 pt-1'>
          <Button
            className='flex items-center gap-1 px-1 py-[2px]'
            style={{ filter: 'hue-rotate(-65deg)' }}
            onClick={() => {
              if (matches.length != totalMatches.length) {
                setMatches(totalMatches)
              } else {
                navigate('/sessions')
              }
            }}
          >
            <Icon width={20} height={20} icon='material-symbols:chevron-left' />
            <span className='relative top-[2px]'>{t('goBack')}</span>
          </Button>
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
                <Table.Th className='text-center'>{t('replayId')}</Table.Th>
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
                  <Table.Td
                    interactive
                    onClick={() => filterLog('opponentLeague', log.opponentLeague)}
                  >
                    {t(log.opponentLeague as LocalizationKey)}
                  </Table.Td>
                  <Table.Td
                    interactive
                    className='text-center'
                    onClick={() => filterLog('opponentCharacter', log.opponentCharacter)}
                  >
                    {log.opponentCharacter}
                  </Table.Td>
                  <Table.Td
                    className='text-center'
                    style={{ color: log.victory === true ? 'lime' : 'red' }}
                  >
                    {log.victory === true ? 'W' : 'L'}
                  </Table.Td>
                  <Table.Td
                    onClick={() => navigator.clipboard.writeText(log.replayId)}
                    className='group text-center'
                    interactive
                  >
                    <Tooltip text={t('copy')}>
                      <span className='block w-12 overflow-hidden overflow-ellipsis'>
                        {log.replayId}
                      </span>
                    </Tooltip>
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
