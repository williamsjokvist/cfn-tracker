import React from 'react'
import { createRoot } from 'react-dom/client'

import { TrackingMachineContext } from '@/state/tracking-machine'
import { AuthMachineContext } from '@/state/auth-machine'

import { cn } from '@/helpers/cn'

import { ErrorPopupProvider } from './error-popup'
import { ConfigProvider } from './config'
import { RouterProvider } from './router'
import { I18nProvider } from './i18n'

import './style.sass'

function AppLoader() {
  return (
    <main className='grid h-screen w-full items-center justify-center text-white'>
      <i
        aria-label='loading'
        className={cn(
          'inline-block h-12 w-12',
          'animate-spin rounded-full',
          'border-[4px] border-current border-t-transparent',
          'text-white'
        )}
        role='status'
      />
    </main>
  )
}

createRoot(document.getElementById('root')!).render(
  <React.Suspense fallback={<AppLoader />}>
    <I18nProvider>
      <AuthMachineContext.Provider>
        <TrackingMachineContext.Provider>
          <React.Suspense fallback={<AppLoader />}>
            <ErrorPopupProvider>
              <ConfigProvider>
                <RouterProvider />
              </ConfigProvider>
            </ErrorPopupProvider>
          </React.Suspense>
        </TrackingMachineContext.Provider>
      </AuthMachineContext.Provider>
    </I18nProvider>
  </React.Suspense>
)
