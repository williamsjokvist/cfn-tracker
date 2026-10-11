import { useLoaderData, useNavigate, useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'

import type { model } from '@model'

import * as Page from '@/ui/page'
import * as Table from '@/ui/table'
import { Button } from '@/ui/button'
import { cn } from '@/helpers/cn'

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

const winRate = (wins: number, losses: number) =>
  wins + losses === 0 ? 0 : Math.floor((wins / (wins + losses)) * 100)

export function SessionOverviewPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { sessionId } = useParams()
  const matches = (useLoaderData() ?? []) as model.Match[]

  const characters = statsByCharacter(matches)
  const wins = characters.reduce((sum, c) => sum + c.wins, 0)
  const losses = characters.reduce((sum, c) => sum + c.losses, 0)
  const lastMatch = matches[0]
  // Tekken has no LP or MR.
  const hasPoints = matches.some(m => m.lp !== 0 || m.mr !== 0)

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('statistics')}</Page.Title>
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
            onClick={() => navigate(`/sessions/${sessionId}/matches`)}
          >
            <Icon width={18} height={18} icon='material-symbols:list-alt-outline' />
            <span className='relative top-[2px]'>{t('history')}</span>
          </Button>
        </div>

        <dl className='mx-8 mt-3 flex flex-wrap gap-2'>
          {lastMatch && <Stat text={t('user')} value={lastMatch.userName} />}
          <Stat text={t('wins')} value={wins} />
          <Stat text={t('losses')} value={losses} />
          <Stat text={t('winRate')} value={`${winRate(wins, losses)}%`} />
        </dl>

        {characters.length > 0 && (
          <Table.Page className='max-h-[250px]'>
            <Table.Content>
              <thead>
                <Table.Tr>
                  <Table.Th>{t('character')}</Table.Th>
                  <Table.Th className='text-center'>{t('wins')}</Table.Th>
                  <Table.Th className='text-center'>{t('losses')}</Table.Th>
                  <Table.Th className='text-center'>{t('winRate')}</Table.Th>
                  {hasPoints && (
                    <>
                      <Table.Th className='text-center'>{t('lpGain')}</Table.Th>
                      <Table.Th className='text-center'>{t('mrGain')}</Table.Th>
                    </>
                  )}
                </Table.Tr>
              </thead>
              <tbody>
                {characters.map(c => (
                  <Table.Tr key={c.character}>
                    <Table.Td className='font-bold'>{c.character}</Table.Td>
                    <Table.Td className='text-center'>{c.wins}</Table.Td>
                    <Table.Td className='text-center'>{c.losses}</Table.Td>
                    <Table.Td className='text-center'>{winRate(c.wins, c.losses)}%</Table.Td>
                    {hasPoints && (
                      <>
                        <Table.Td className='text-center'>
                          <Gain value={c.lpGain} />
                        </Table.Td>
                        <Table.Td className='text-center'>
                          <Gain value={c.mrGain} />
                        </Table.Td>
                      </>
                    )}
                  </Table.Tr>
                ))}
              </tbody>
            </Table.Content>
          </Table.Page>
        )}
      </motion.div>
    </Page.Root>
  )
}

function Stat({ text, value }: { text: string; value: string | number }) {
  return (
    <div className='flex items-baseline gap-3 rounded-xl bg-slate-50/5 px-3 pt-2 pb-1'>
      <dt className='font-extralight tracking-wider'>{text}</dt>
      <dd className='text-xl font-semibold'>{value}</dd>
    </div>
  )
}

function Gain({ value }: { value: number }) {
  return (
    <span className={cn(value > 0 && 'text-lime-400', value < 0 && 'text-red-400')}>
      {value > 0 ? `+${value}` : value}
    </span>
  )
}
