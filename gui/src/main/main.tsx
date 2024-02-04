import React from 'react'
import { createRoot } from 'react-dom/client'

import { TrackingMachineContext } from '@/state/tracking-machine'
import { AuthMachineContext } from '@/state/auth-machine'

import { AppLoader } from './app-layout/app-loader'
import { ErrorPopupProvider } from './error-popup'
import { ConfigProvider } from './config'
import { RouterProvider } from './router'
import { I18nProvider } from './i18n'

import './style.sass'

createRoot(document.getElementById('root')!).render(
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
)
