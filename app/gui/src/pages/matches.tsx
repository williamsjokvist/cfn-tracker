import React from 'react'
import { useLoaderData, useSearchParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

import * as Page from '@/ui/page'
import * as Table from '@/ui/table'

import type { model } from '@model'
import { type LocalizationKey } from '@/main/i18n'
import { Tooltip } from '@/ui/tooltip'
import { SessionBreadcrumbs, type SessionData } from '@/pages/session'

export function MatchesListPage() {
  const { t } = useTranslation()
  const { session, matches: sessionMatches } = useLoaderData() as SessionData
  const [searchParams] = useSearchParams()
  const character = searchParams.get('character') ?? undefined
  const totalMatches = React.useMemo(
    () => (character ? sessionMatches.filter(m => m.character === character) : sessionMatches),
    [sessionMatches, character]
  )
  const [filters, setFilters] = React.useState<Partial<Record<keyof model.Match, string>>>({})
  const matches = totalMatches.filter(m =>
    Object.entries(filters).every(
      ([property, value]) => String(m[property as keyof model.Match]).toLowerCase() === value
    )
  )

  const filterLog = (property: keyof model.Match, value: string) => {
    const { [property]: active, ...rest } = filters
    setFilters(
      active === value.toLowerCase() ? rest : { ...filters, [property]: value.toLowerCase() }
    )
  }

  if (matches.length === 0) {
    return (
      <Page.Root>
        <Page.Header>
          <SessionBreadcrumbs session={session} matchLog character={character} />
        </Page.Header>
      </Page.Root>
    )
  }

  return (
    <Page.Root>
      <Page.Header>
        <SessionBreadcrumbs session={session} matchLog character={character} />
      </Page.Header>

      <div className='relative w-full'>
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
                      <span className='block w-12 overflow-hidden text-ellipsis'>
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
