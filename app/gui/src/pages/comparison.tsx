import React from 'react'
import { useTranslation } from 'react-i18next'
import { useSelector } from '@xstate/react'
import { Link } from 'react-router-dom'

import * as Page from '@/ui/page'
import { AuthMachineContext } from '@/state/auth-machine'
import { TrackingMachineContext } from '@/state/tracking-machine'
import { useErrorPopup } from '@/main/error-popup'
import { Button } from '@/ui/button'
import * as Table from '@/ui/table'
import { cn } from '@/helpers/cn'

import { GetSF6BattleStatsComparison } from '@cmd/ComparisonHandler'
import { model } from '@model'

export function SF6ComparisonPage() {
  const { t } = useTranslation()
  const authActor = AuthMachineContext.useActorRef()
  const game = useSelector(authActor, ({ context }) => context.game)

  const trackingActor = TrackingMachineContext.useActorRef()
  const trackedUserCode = useSelector(trackingActor, ({ context }) => context.user?.code ?? '')

  const setError = useErrorPopup()

  const [userCode, setUserCode] = React.useState('')
  const [loading, setLoading] = React.useState(false)
  const [report, setReport] = React.useState<model.SF6BattleStatsComparisonReport | null>(null)

  React.useEffect(() => {
    if (!userCode && trackedUserCode) {
      setUserCode(trackedUserCode)
    }
  }, [trackedUserCode])

  const metrics = React.useMemo(() => {
    const items = report?.metrics ?? []
    const withDelta = items.map(m => ({
      ...m,
      delta: m.current - m.topAvg,
      absDelta: Math.abs(m.current - m.topAvg)
    }))
    withDelta.sort((a, b) => b.absDelta - a.absDelta)
    return withDelta
  }, [report])

  const fetchComparison = async () => {
    if (!userCode) {
      return
    }
    setReport(null)
    setError(null as unknown as model.FGCTrackerError)
    setLoading(true)
    try {
      const res = await GetSF6BattleStatsComparison(userCode)
      setReport(res)
    } catch (e) {
      setError(e as unknown as model.FGCTrackerError)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('sf6Comparison')}</Page.Title>
      </Page.Header>
      <div className='flex h-full min-h-0 flex-col gap-4 px-4 pt-3 text-base text-gray-200'>
        {game !== model.GameType.STREET_FIGHTER_6 ? (
          <div className='flex flex-col gap-3'>
            <div>{t('pickGame')}</div>
            <Link className='text-[#bfbcff] underline underline-offset-4' to='/tracking'>
              {t('tracking')}
            </Link>
          </div>
        ) : (
          <>
            <div className='flex flex-wrap items-end justify-between gap-3'>
              <div className='flex flex-col gap-2'>
                <label className='text-base text-white/80'>{t('cfnName')}</label>
                <input
                  value={userCode}
                  onChange={e => setUserCode(e.target.value)}
                  className='w-full max-w-[320px] rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-white outline-none focus:border-white/25'
                  placeholder={t('enterCfnName')}
                  autoCapitalize='off'
                  autoComplete='off'
                  autoCorrect='off'
                  autoSave='off'
                />
              </div>
              <Button
                disabled={!userCode || loading}
                onClick={fetchComparison}
                style={{ filter: 'hue-rotate(-160deg)' }}
                className='h-[46px]'
              >
                {loading ? t('loading') : t('refresh')}
              </Button>
            </div>

            {report && (
              <div className='flex flex-wrap items-center justify-between gap-3 rounded-xl bg-white/5 px-4 py-2'>
                <div className='flex flex-col leading-tight'>
                  <span className='text-base text-white/70'>{t('character')}</span>
                  <span className='text-lg font-semibold text-white'>{report.characterName}</span>
                </div>
                <div className='flex flex-col text-right leading-tight'>
                  <span className='text-base text-white/70'>Top Players</span>
                  <span className='text-lg font-semibold text-white'>Top {report.topN} Players</span>
                </div>
              </div>
            )}

            {report && (
              <div className='min-h-0 flex-1 overflow-y-auto'>
                <Table.Content className='border-separate border-spacing-y-1'>
                  <thead>
                    <Table.Tr>
                      <Table.Th className='w-[52%]'>Metric</Table.Th>
                      <Table.Th className='w-[16%] text-right'>You</Table.Th>
                      <Table.Th className='w-[16%] text-right'>Top Avg</Table.Th>
                      <Table.Th className='w-[16%] text-right'>Δ</Table.Th>
                    </Table.Tr>
                  </thead>
                  <tbody>
                    {metrics.map(m => (
                      <Table.Tr key={m.key}>
                        <Table.Td className='whitespace-normal leading-tight'>
                          <span className='text-white'>{m.key}</span>
                          {m.kind === 'int' && m.unit && (
                            <span className='ml-2 text-xs text-white/60'>({m.unit})</span>
                          )}
                        </Table.Td>
                        <Table.Td className='text-right font-semibold tabular-nums'>
                          {formatValue(m, m.current)}
                        </Table.Td>
                        <Table.Td className='text-right tabular-nums'>
                          {formatValue(m, m.topAvg)}
                        </Table.Td>
                        <Table.Td
                          className={cn(
                            'text-right font-semibold tabular-nums',
                            m.delta > 0 ? 'text-green-300' : '',
                            m.delta < 0 ? 'text-red-300' : ''
                          )}
                        >
                          {formatDelta(m, m.delta)}
                        </Table.Td>
                      </Table.Tr>
                    ))}
                  </tbody>
                </Table.Content>
              </div>
            )}
          </>
        )}
      </div>
    </Page.Root>
  )
}

function formatValue(metric: model.SF6BattleStatsMetric, value: number) {
  if (metric.kind === 'int') {
    return `${Math.round(value)}`
  }
  return `${value.toFixed(2)}${metric.unit ? metric.unit : ''}`
}

function formatDelta(metric: model.SF6BattleStatsMetric, delta: number) {
  if (metric.kind === 'int') {
    const v = Math.round(delta)
    return `${v > 0 ? '+' : ''}${v}`
  }
  return `${delta > 0 ? '+' : ''}${delta.toFixed(2)}${metric.unit ? metric.unit : ''}`
}

