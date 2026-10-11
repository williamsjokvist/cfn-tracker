import React from 'react'
import { useTranslation } from 'react-i18next'

import type { model } from '@model'

import * as Select from '@/ui/select'
import { cn } from '@/helpers/cn'
import { signed, winRate } from '@/helpers/stats'

type Metric = 'mr' | 'lp' | 'winRate' | 'netWins'

type Point = { date: Date; value: number }

function toPoints(history: model.HistoryMatch[], metric: Metric): Point[] {
  let wins = 0
  let losses = 0
  const points: Point[] = []
  for (const m of history) {
    if (m.victory) wins++
    else losses++
    const value = {
      mr: m.mr,
      lp: m.lp,
      winRate: winRate(wins, losses),
      netWins: wins - losses
    }[metric]
    // Characters have no MR before they reach Master.
    if (metric === 'mr' && value === 0) continue
    points.push({ date: new Date(`${m.date}T${m.time}`), value })
  }
  return points
}

type Span = { key: string; unit?: 'week' | 'month' | 'year'; count?: number }

const spans: Span[] = [
  { key: '1w', unit: 'week', count: 1 },
  { key: '1m', unit: 'month', count: 1 },
  { key: '3m', unit: 'month', count: 3 },
  { key: '6m', unit: 'month', count: 6 },
  { key: '1y', unit: 'year', count: 1 },
  { key: 'all' }
]

// Spans count back from the latest match rather than today, so players who haven't
// played in a while still get every span.
function pointsInSpan(points: Point[], span: Span) {
  if (!span.unit || !span.count || points.length === 0) return points
  const start = new Date(points[points.length - 1].date)
  if (span.unit === 'week') start.setDate(start.getDate() - 7 * span.count)
  if (span.unit === 'month') start.setMonth(start.getMonth() - span.count)
  if (span.unit === 'year') start.setFullYear(start.getFullYear() - span.count)
  return points.filter(p => p.date >= start)
}

const change = (points: Point[]) =>
  points.length > 0 ? points[points.length - 1].value - points[0].value : 0

export function HistoryChart({
  history,
  mainCharacter,
  timeSpans = true
}: {
  // Oldest match first.
  history: model.HistoryMatch[]
  mainCharacter: string
  // A single session is too short for spans of weeks and months.
  timeSpans?: boolean
}) {
  const { t, i18n } = useTranslation()
  const characters = [...new Set(history.map(h => h.character))]
  const [character, setCharacter] = React.useState(mainCharacter)
  const characterHistory = history.filter(h => h.character === character)

  // Tekken has no LP or MR.
  const metrics: Metric[] = [
    ...(characterHistory.some(h => h.mr !== 0) ? (['mr'] as const) : []),
    ...(characterHistory.some(h => h.lp !== 0) ? (['lp'] as const) : []),
    'winRate',
    'netWins'
  ]
  const [chosenMetric, setMetric] = React.useState<Metric>()
  const metric = chosenMetric && metrics.includes(chosenMetric) ? chosenMetric : metrics[0]
  const metricLabel = (m: Metric) =>
    ({ mr: 'MR', lp: 'LP', winRate: t('winRate'), netWins: t('netWins') })[m]

  const allPoints = toPoints(characterHistory, metric)
  const spanPoints = spans.map(span => pointsInSpan(allPoints, span))
  const [chosenSpan, setSpan] = React.useState(spans.length - 1)
  const span = spanPoints[chosenSpan].length >= 2 ? chosenSpan : spans.length - 1
  const points = spanPoints[span]
  const [hovered, setHovered] = React.useState<number>()
  const format = Intl.DateTimeFormat(i18n.resolvedLanguage, {
    dateStyle: 'medium',
    timeStyle: 'short'
  })

  const formatValue = (value: number) => (metric === 'winRate' ? `${value}%` : `${value}`)
  const formatChange = (value: number) =>
    metric === 'winRate' ? `${signed(value)}%` : signed(value)
  const spanLabel = ({ unit, count }: Span) =>
    unit && count
      ? new Intl.NumberFormat(i18n.resolvedLanguage, {
          style: 'unit',
          unit,
          unitDisplay: 'short'
        }).format(count)
      : t('allTime')

  const values = points.map(p => p.value)
  const min = Math.min(...values)
  const max = Math.max(...values)
  // Percentages of the plot, leaving room so the line never touches the edges.
  const x = (i: number) => (points.length === 1 ? 50 : (i / (points.length - 1)) * 100)
  const y = (v: number) => (max === min ? 50 : 90 - ((v - min) / (max - min)) * 80)
  const line = points.map((p, i) => `${x(i)},${y(p.value)}`).join(' ')
  const delta = change(points)
  const hoveredPoint = hovered !== undefined ? points[hovered] : undefined

  return (
    <section className='rounded-xl bg-slate-50/5 p-3'>
      <div className='mb-2 flex items-start justify-between gap-4'>
        <h3 className='font-extralight tracking-wider'>
          {t('statOverTime', { stat: metricLabel(metric) })}
        </h3>
        <div className='flex items-center gap-2'>
          {characters.length > 1 && (
            <Select.Root size='sm' className='w-40' value={character} onValueChange={setCharacter}>
              {characters.map(c => (
                <Select.Item key={c} value={c}>
                  {c}
                </Select.Item>
              ))}
            </Select.Root>
          )}
          <Select.Root
            size='sm'
            className='w-40'
            value={metric}
            onValueChange={m => setMetric(m as Metric)}
          >
            {metrics.map(m => (
              <Select.Item key={m} value={m}>
                {metricLabel(m)}
              </Select.Item>
            ))}
          </Select.Root>
          <span
            className={cn(
              'w-16 text-end text-lg font-semibold tabular-nums',
              delta > 0 && 'text-green-300',
              delta < 0 && 'text-red-300'
            )}
          >
            {formatChange(delta)}
          </span>
        </div>
      </div>
      {points.length < 2 ? (
        <p className='flex h-36 items-center justify-center text-white/60'>
          {t('notEnoughMatches')}
        </p>
      ) : (
        <div
          className='relative h-36 text-green-300'
          onPointerMove={e => {
            const rect = e.currentTarget.getBoundingClientRect()
            const ratio = (e.clientX - rect.left) / rect.width
            setHovered(Math.round(Math.min(Math.max(ratio, 0), 1) * (points.length - 1)))
          }}
          onPointerLeave={() => setHovered(undefined)}
        >
          <svg
            viewBox='0 0 100 100'
            preserveAspectRatio='none'
            className='absolute inset-0 h-full w-full overflow-visible'
            aria-hidden
          >
            <defs>
              <linearGradient id='history-chart-fill' x1='0' y1='0' x2='0' y2='1'>
                <stop offset='0%' stopColor='currentColor' stopOpacity={0.25} />
                <stop offset='100%' stopColor='currentColor' stopOpacity={0} />
              </linearGradient>
            </defs>
            <polygon points={`0,100 ${line} 100,100`} fill='url(#history-chart-fill)' />
            <polyline
              points={line}
              fill='none'
              stroke='currentColor'
              strokeWidth={2}
              strokeLinejoin='round'
              strokeLinecap='round'
              vectorEffect='non-scaling-stroke'
            />
          </svg>
          {hoveredPoint && hovered !== undefined && (
            <>
              <div
                className='pointer-events-none absolute inset-y-0 w-px bg-white/20'
                style={{ left: `${x(hovered)}%` }}
              />
              <div
                className='pointer-events-none absolute h-2.5 w-2.5 -translate-x-1/2 -translate-y-1/2 rounded-full bg-current ring-2 ring-slate-900'
                style={{ left: `${x(hovered)}%`, top: `${y(hoveredPoint.value)}%` }}
              />
              <div
                className={cn(
                  'pointer-events-none absolute top-0 rounded-lg bg-black/70 px-2 py-1 text-sm whitespace-nowrap text-white backdrop-blur-xl',
                  x(hovered) > 50 ? '-translate-x-[calc(100%+8px)]' : 'translate-x-2'
                )}
                style={{ left: `${x(hovered)}%` }}
              >
                <div className='text-white/60'>{format.format(hoveredPoint.date)}</div>
                <div className='font-semibold tabular-nums'>{formatValue(hoveredPoint.value)}</div>
              </div>
            </>
          )}
        </div>
      )}
      {timeSpans && (
        <div role='tablist' className='mt-3 grid grid-cols-6 border-t border-slate-50/10'>
          {spans.map((s, i) => {
            const spanChange = change(spanPoints[i])
            return (
              <button
                key={s.key}
                type='button'
                role='tab'
                aria-selected={i === span}
                disabled={spanPoints[i].length < 2}
                onClick={() => setSpan(i)}
                className={cn(
                  'flex flex-col items-center border-b-2 border-transparent pt-2 pb-1 tabular-nums transition-colors',
                  'hover:bg-slate-50/5 disabled:pointer-events-none disabled:opacity-40',
                  i === span && 'border-green-300 font-semibold'
                )}
              >
                <span>{spanLabel(s)}</span>
                <span
                  className={cn(
                    'text-sm',
                    spanChange > 0 && 'text-green-300',
                    spanChange < 0 && 'text-red-300'
                  )}
                >
                  {spanPoints[i].length < 2 ? '–' : formatChange(spanChange)}
                </span>
              </button>
            )
          })}
        </div>
      )}
    </section>
  )
}
