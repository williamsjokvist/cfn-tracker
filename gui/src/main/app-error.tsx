import React from 'react'
import { useTranslation } from 'react-i18next'
import { useRouteError } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'

import { model } from '@model'
import * as Page from '@/ui/page'

import { AppTitleBar } from './app-titlebar'

export function AppErrorBoundary() {
  const err = useFGCTrackerError()
  return (
    <ErrorWrapper err={err}>
      <AppTitleBar />
    </ErrorWrapper>
  )
}

export function PageErrorBoundary() {
  const { t } = useTranslation()
  const err = useFGCTrackerError()
  if (!err) {
    return null
  }
  return (
    <ErrorWrapper err={err}>
      <Page.Header>
        <Page.Title>{t(err.localizationKey)}</Page.Title>
      </Page.Header>
    </ErrorWrapper>
  )
}

const isFGCTrackerError = (error: unknown) => error instanceof Object && 'localizationKey' in error

function useFGCTrackerError() {
  const thrownError = useRouteError()
  const [err, setErr] = React.useState<model.FGCTrackerError>()

  React.useEffect(() => {
    console.error(thrownError)
    if (isFGCTrackerError(thrownError)) {
      setErr(thrownError as model.FGCTrackerError)
    } else if (thrownError instanceof Error) {
      setErr({
        localizationKey: model.ErrorLocalizationKey.errUnknown,
        message: thrownError.message,
        InnerError: thrownError
      })
    }
  }, [thrownError])

  return err
}

function ErrorWrapper(props: React.PropsWithChildren & { err?: model.FGCTrackerError }) {
  const { t } = useTranslation()
  if (!props.err) {
    return null
  }
  return (
    <motion.section
      className='h-screen w-full text-white'
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.125 }}
    >
      {props.children}
      <div className='mt-8 flex w-full flex-col items-center justify-center rounded-md pb-16 text-center'>
        <Icon icon='material-symbols:warning-outline' className='h-40 w-40 text-[#ff6388]' />
        <h1 className='text-center text-2xl font-bold'>{t(props.err.localizationKey)}</h1>
        <p className='text-xl'>{props.err?.message}</p>
      </div>
    </motion.section>
  )
}
