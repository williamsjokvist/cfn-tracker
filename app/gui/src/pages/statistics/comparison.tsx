import React from 'react'
import { useTranslation } from 'react-i18next'
import { useSearchParams } from 'react-router-dom'
import { Icon } from '@iconify/react'

import * as Page from '@/ui/page'
import { LoginPrompt } from '@/ui/login-prompt'
import { Banner } from '@/ui/banner'
import type { LocalizationKey } from '@/main/i18n'
import { useErrorPopup } from '@/main/error-popup'
import * as Table from '@/ui/table'
import { Button } from '@/ui/button'
import { cn } from '@/helpers/cn'

import {
  CancelSF6BattleStatsComparison,
  ExportSF6BattleStatsComparison,
  GetSF6BattleStatsComparison
} from '@cmd/ComparisonHandler'
import type { model } from '@model'
import { EventsOff, EventsOn } from '@runtime'

type ComparisonRun = {
  userCode: string
  runId: string
  report: Promise<model.SF6BattleStatsComparisonReport>
  settled: boolean
  cancelTimer?: number
}

function startComparison(userCode: string): ComparisonRun {
  const runId = crypto.randomUUID()
  const run: ComparisonRun = {
    userCode,
    runId,
    settled: false,
    report: GetSF6BattleStatsComparison(userCode, runId).finally(() => {
      run.settled = true
    })
  }
  return run
}

export function ComparisonPage() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [searchParams] = useSearchParams()
  const userCode = searchParams.get('user') ?? ''

  const [loading, setLoading] = React.useState(false)
  const [report, setReport] = React.useState<model.SF6BattleStatsComparisonReport | null>(null)
  const [authAction, setAuthAction] = React.useState<{
    localizationKey: string
    secondsLeft: number
  } | null>(null)

  const [progress, setProgress] = React.useState<{
    localizationKey: string
    params?: Record<string, unknown>
  } | null>(null)

  React.useEffect(() => {
    EventsOn('comparison-progress', progress => {
      // A step after the login prompt means the user is done in the Chrome window.
      setAuthAction(null)
      setProgress(progress)
    })
    EventsOn('comparison-auth-action-required', setAuthAction)
    EventsOn('comparison-auth-success', () => setAuthAction(null))
    return () => {
      EventsOff('comparison-progress')
      EventsOff('comparison-auth-action-required')
      EventsOff('comparison-auth-success')
    }
  }, [])

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

  const run = React.useRef<ComparisonRun | null>(null)
  React.useEffect(() => {
    if (!userCode) return
    let current = run.current
    if (current?.userCode === userCode && current.cancelTimer !== undefined) {
      // StrictMode and Fast Refresh remount the page right away, so carry on with
      // the run instead of interrupting it and starting over.
      window.clearTimeout(current.cancelTimer)
      current.cancelTimer = undefined
    } else {
      setReport(null)
      setProgress(null)
      setLoading(true)
      current = startComparison(userCode)
      run.current = current
    }
    const thisRun = current
    let left = false
    thisRun.report
      .then(report => !left && setReport(report))
      .catch(err => !left && setError(err))
      .finally(() => {
        if (left) return
        setLoading(false)
        setAuthAction(null)
        setProgress(null)
      })
    return () => {
      left = true
      if (thisRun.settled) return
      thisRun.cancelTimer = window.setTimeout(() => CancelSF6BattleStatsComparison(thisRun.runId))
    }
  }, [userCode])

  return (
    <Page.Root>
      <Page.Header>
        <Page.Breadcrumbs
          items={[
            {
              label: t('statistics'),
              to: `/statistics?${new URLSearchParams({ user: userCode })}`
            },
            {
              label: report?.userName
                ? `${t('sf6Comparison')} · ${report.userName}`
                : t('sf6Comparison')
            }
          ]}
        />
        <div className='flex items-center gap-3'>
          {loading && <Page.LoadingIcon />}
          {report && (
            <Button
              className='flex items-center gap-2 px-2 py-[2px]'
              // Clickable although the header is the window's drag region.
              style={
                { filter: 'hue-rotate(145deg)', '--draggable': 'no-drag' } as React.CSSProperties
              }
              onClick={() => ExportSF6BattleStatsComparison(report).catch(setError)}
            >
              <Icon width={18} height={18} icon='material-symbols:download' />
              <span className='relative top-[2px]'>{t('export')}</span>
            </Button>
          )}
        </div>
      </Page.Header>
      {authAction ? (
        <LoginPrompt />
      ) : (
        <div className='flex h-full min-h-0 flex-col gap-4 px-8 pt-3 text-base text-gray-200'>
          {report && (
            <div className='min-h-0 flex-1 overflow-y-auto'>
              <Table.Content className='border-separate border-spacing-y-1'>
                <thead>
                  <Table.Tr>
                    <Table.Th className='w-[52%]'>{t('comparisonMetric')}</Table.Th>
                    <Table.Th className='w-[16%] text-end!'>{report.userName}</Table.Th>
                    <Table.Th className='w-[16%] text-end!'>{t('comparisonTopAvg')}</Table.Th>
                    <Table.Th className='w-[16%] text-end!'>
                      <Icon icon='mdi:delta' aria-label='Δ' className='inline h-5 w-5' />
                    </Table.Th>
                  </Table.Tr>
                </thead>
                <tbody>
                  {metrics.map(m => {
                    const effective = (m.polarity ?? 1) * m.delta
                    return (
                      <Table.Tr key={m.key}>
                        <Table.Td className='leading-tight whitespace-normal'>
                          <span className='text-white'>{m.key}</span>
                          {m.kind === 'int' && m.unit && (
                            <span className='ml-2 text-xs text-white/60'>({m.unit})</span>
                          )}
                        </Table.Td>
                        <Table.Td className='text-end font-semibold tabular-nums'>
                          {formatValue(m, m.current)}
                        </Table.Td>
                        <Table.Td className='text-end tabular-nums'>
                          {formatValue(m, m.topAvg)}
                        </Table.Td>
                        <Table.Td
                          className={cn(
                            'text-end font-semibold tabular-nums',
                            effective > 0 ? 'text-green-300' : '',
                            effective < 0 ? 'text-red-300' : ''
                          )}
                        >
                          {formatDelta(m, m.delta)}
                        </Table.Td>
                      </Table.Tr>
                    )
                  })}
                </tbody>
              </Table.Content>
            </div>
          )}
        </div>
      )}
      <Banner open={loading} variant='loading'>
        <span>
          {authAction
            ? t(authAction.localizationKey as 'authNeedRelogin', {
                seconds: authAction.secondsLeft
              })
            : progress
              ? t(progress.localizationKey as LocalizationKey, progress.params)
              : t('loading')}
        </span>
      </Banner>
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
