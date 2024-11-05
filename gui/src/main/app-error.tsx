import React from 'react'
import { useTranslation } from 'react-i18next'
import { useRouteError } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Icon } from '@iconify/react'

import { type model } from '@model'
import * as Page from '@/ui/page'

import { AppTitleBar } from './app-titlebar'
import { LocalizationKey } from './i18n'

export function AppErrorBoundary() {
  const err = useFormattedError()
  return (
    <ErrorWrapper err={err}>
      <AppTitleBar />
    </ErrorWrapper>
  )
}

export function PageErrorBoundary() {
  const { t } = useTranslation()
  const err = useFormattedError()
  if (!err?.localizationKey) {
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

const isFormattedError = (error: unknown) => error instanceof Object && 'translationKey' in error

function useFormattedError() {
  const thrownError = useRouteError()
  const [err, setErr] = React.useState<model.FormattedError | unknown>()

  React.useEffect(() => {
    console.error(thrownError)
    if (thrownError instanceof Error) {
      setErr({ translationKey: '', message: thrownError.message, error: thrownError })
    } else if (isFormattedError(thrownError)) {
      setErr(thrownError)
    }
  }, [thrownError])

  return err as model.FormattedError
}

function ErrorWrapper(props: React.PropsWithChildren & { err?: model.FormattedError }) {
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
        <h1 className='text-center text-2xl font-bold'>{t(props.err?.localizationKey)}</h1>
        <p className='text-xl'>{props.err?.message}</p>
      </div>
    </motion.section>
  )
}
