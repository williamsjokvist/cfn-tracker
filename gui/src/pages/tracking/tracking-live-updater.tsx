import React from 'react'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'
import { useTranslation } from 'react-i18next'
import { useSelector } from '@xstate/react'
import { PieChart } from 'react-minimal-pie-chart'

import { TrackingMachineContext } from '@/state/tracking-machine'
import { Button } from '@/ui/button'
import { Tooltip } from '@/ui/tooltip'
import * as Page from '@/ui/page'
import { type LocalizationKey } from '@/main/i18n'

export function TrackingLiveUpdater() {
  const { t } = useTranslation()
  const trackingActor = TrackingMachineContext.useActorRef()

  const {
    lp,
    mr,
    wins,
    losses,
    winStreak,
    lpGain,
    mrGain,
    winRate,
    opponent,
    opponentCharacter,
    character,
    opponentLeague,
    userName,
    victory
  } = useSelector(trackingActor, ({ context: { match } }) => match)

  const [refreshDisabled, setRefreshDisabled] = React.useState(false)

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('tracking')}</Page.Title>
        <Page.LoadingIcon />
      </Page.Header>
      <motion.section
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='h-full px-6 pt-4'
      >
        <dl className='flex w-full items-center justify-between whitespace-nowrap'>
          <SmallStat text='CFN' value={userName} />
          <div className='flex justify-between gap-8'>
            {lp !== 0 && <SmallStat text='LP' value={`${lp == -1 ? t('placement') : lp}`} />}
            {mr !== 0 && <SmallStat text='MR' value={`${mr == -1 ? t('placement') : mr}`} />}
          </div>
        </dl>
        <div className='flex h-[calc(100%-32px)] flex-1 pt-3 pb-5'>
          <div className='w-full'>
            <dl className='text-lg whitespace-nowrap'>
              <div className='flex justify-between gap-2'>
                <BigStat text={t('wins')} value={wins} />
                <BigStat text={t('losses')} value={losses} />
              </div>
              <div className='flex justify-between gap-2'>
                <BigStat text={t('winRate')} value={`${winRate}%`} />
                <BigStat text={t('winStreak')} value={winStreak} />
              </div>
              <div className='flex justify-between gap-2'>
                {lpGain > 0 && (
                  <BigStat text={t('lpGain')} value={`${lpGain > 0 ? `+` : ``}${lpGain}`} />
                )}
                {mrGain > 0 && (
                  <BigStat text={t('mrGain')} value={`${mrGain > 0 ? `+` : ``}${mrGain}`} />
                )}
              </div>
            </dl>
            {opponent != '' && (
              <div className='group flex items-center justify-between rounded-xl bg-slate-50/5 p-3 pb-2 text-lg leading-none'>
                <span>{t('lastMatch')}</span>
                <div className='relative flex items-center gap-2'>
                  <Icon
                    icon={victory ? 'twemoji:victory-hand' : 'twemoji:pensive-face'}
                    width={25}
                  />{' '}
                  vs
                  <b>{opponent}</b> - {opponentCharacter} ({t(opponentLeague as LocalizationKey)})
                </div>
              </div>
            )}
          </div>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.35 }}
            className='relative grid h-full flex-0'
          >
            <PieChart
              className='mx-auto h-52 w-full'
              animate
              animationDuration={2000}
              lineWidth={85}
              data={[
                {
                  title: t('wins'),
                  value: wins,
                  color: 'rgba(0, 255, 77, .65)'
                },
                {
                  title: t('losses'),
                  value: wins == 0 && losses == 0 ? 1 : losses,
                  color: 'rgba(251, 73, 73, 0.25)'
                }
              ]}
            >
              <defs>
                <linearGradient id='blue-gradient' direction={-65}>
                  <stop offset='0%' stopColor='#20BF55' />
                  <stop offset='100%' stopColor='#347fd0' />
                </linearGradient>
                <linearGradient id='red-gradient' direction={120}>
                  <stop offset='0%' stopColor='#EC9F05' />
                  <stop offset='100%' stopColor='#EE9617' />
                </linearGradient>
              </defs>
            </PieChart>
            <div className='flex justify-between gap-2 self-end'>
              <Tooltip text={t('cooldown')} disabled={!refreshDisabled}>
                <Button
                  style={{ filter: 'hue-rotate(-160deg)' }}
                  disabled={refreshDisabled}
                  onClick={() => {
                    if (!refreshDisabled) {
                      trackingActor.send({ type: 'forcePoll' })
                      setRefreshDisabled(true)
                      setTimeout(() => setRefreshDisabled(false), 15000)
                    }
                  }}
                >
                  <Icon icon='fa6-solid:recycle' className='mr-3 h-5 w-5' />
                  {t('refresh')}
                </Button>
              </Tooltip>
              <Button onClick={() => trackingActor.send({ type: 'cease' })}>
                <Icon icon='fa6-solid:stop' className='mr-3 h-5 w-5' />
                {t('stop')}
              </Button>
            </div>
          </motion.div>
        </div>
        {/* TODO: fix character image for tekken 8 */}
        <img
          className='pointer-events-none absolute top-0 -right-20 z-[-1] h-full opacity-10 grayscale'
          src={`https://www.streetfighter.com/6/buckler/assets/images/material/character/character_${character
            .toLowerCase()
            .replace(/\s/g, '')
            .replace('.', '')}_r.png`}
          alt={''}
        />
      </motion.section>
    </Page.Root>
  )
}

type StatProps = { text: string; value: string | number }
const BigStat = ({ text, value }: StatProps) => (
  <div className='mb-2 flex flex-1 justify-between gap-4 rounded-xl bg-slate-50/5 p-3 pb-1'>
    <dt className='font-extralight tracking-wider'>{text}</dt>
    <dd className='text-4xl font-semibold'>{value}</dd>
  </div>
)

const SmallStat = ({ text, value }: StatProps) => (
  <div className='flex gap-3 text-2xl'>
    <dt className='text-xl leading-8'>{text}</dt>
    <dd className='font-bold'>{value}</dd>
  </div>
)
