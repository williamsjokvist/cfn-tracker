import React from 'react'
import { useLoaderData } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'

import { TrackingMachineContext } from '@/state/tracking-machine'

import { Button } from '@/ui/button'
import { Checkbox } from '@/ui/checkbox'
import * as Page from '@/ui/page'

import { model } from '@model'

export function TrackingForm() {
  const { t } = useTranslation()
  const users = (useLoaderData() ?? []) as model.User[]
  const trackingActor = TrackingMachineContext.useActorRef()

  const playerIdInputRef = React.useRef<HTMLInputElement>(null)
  const restoreRef = React.useRef<HTMLInputElement>(null)
  const [playerIdInput, setPlayerIdInput] = React.useState('')

  const onSubmit: React.FormEventHandler<HTMLFormElement> = e => {
    e.preventDefault()
    if (playerIdInput == '') return
    trackingActor.send({
      type: 'submit',
      user: {
        displayName: users.find(old => old.code == playerIdInput) ?? playerIdInput,
        code: playerIdInput
      },
      restore: restoreRef.current && restoreRef.current.checked
    })
  }

  const playerChipClicked = (player: model.User) => {
    if (playerIdInputRef.current) {
      playerIdInputRef.current.value = player.code
      setPlayerIdInput(player.code)
    }
  }

  const clearInput = () => {
    if (playerIdInputRef.current) {
      playerIdInputRef.current.value = ''
      setPlayerIdInput('')
    }
  }

  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('startTracking')}</Page.Title>
      </Page.Header>
      <motion.form
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='relative flex h-full w-full flex-col gap-5 justify-self-center overflow-x-visible overflow-y-scroll px-56 pb-4 pt-12'
        onSubmit={onSubmit}
      >
        <h3 className='text-lg'>{t('enterCfnName')}</h3>
        <div className='relative'>
          <input
            ref={playerIdInputRef}
            onChange={e => setPlayerIdInput(e.target.value)}
            className='block w-full border-0 border-b-2 border-b-[rgba(255,255,255,0.275)] bg-transparent px-4 pb-3 pr-12 pt-4 text-lg text-gray-300 outline-hidden transition-colors hover:border-white hover:text-white focus:border-white focus:outline-hidden focus:ring-transparent focus:ring-offset-transparent'
            type='text'
            placeholder={t('cfnName')!}
            autoCapitalize='off'
            autoComplete='off'
            autoCorrect='off'
            autoSave='off'
          />
          {playerIdInput.length > 0 && (
            <button
              type='button'
              onClick={clearInput}
              aria-label='Clear'
              className='absolute right-0 top-0 mr-4 mt-4 rounded-md text-[#bfbcff] transition-colors hover:bg-(rgba(255,255,255,.11)) hover:text-white'
            >
              <Icon icon='ci:close-big' className='h-6 w-6' />
            </button>
          )}
        </div>
        {users.length > 0 && (
          <div className='flex flex-wrap content-center items-center gap-2 text-center'>
            {users.map(player => (
              <button
                key={player.displayName}
                type='button'
                onClick={() => playerChipClicked(player)}
                className='items-center whitespace-nowrap rounded-2xl bg-[rgb(255,255,255,0.075)] px-3 pt-1 text-base transition-all hover:bg-[rgb(255,255,255,0.125)]'
              >
                {player.displayName}
              </button>
            ))}
          </div>
        )}
        <footer className='flex w-full items-center'>
          {users.some(old => old.code === playerIdInput) && (
            <div className='group flex items-center gap-4'>
              <Checkbox ref={restoreRef} id='restore-session' />
              <label
                htmlFor='restore-session'
                className='cursor-pointer text-lg text-gray-300 transition-colors group-hover:text-white'
              >
                {t('restoreSession')}
              </label>
            </div>
          )}
          <Button type='submit' style={{ filter: 'hue-rotate(-65deg)' }} className='ml-auto'>
            {t('start')}
          </Button>
        </footer>
      </motion.form>
    </Page.Root>
  )
}
