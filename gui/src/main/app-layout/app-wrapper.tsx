import React from 'react'
import { useTranslation } from 'react-i18next'
import { useSelector } from '@xstate/react'
import { Outlet } from 'react-router-dom'
import { Icon } from '@iconify/react'

import { cn } from '@/helpers/cn'
import { AuthMachineContext } from '@/state/auth-machine'
import { CheckForUpdate } from '@@/go/core/CommandHandler'
import { BrowserOpenURL } from '@@/runtime/runtime'

import { useErrorPopup } from '../error-popup'
import { AppSidebar } from './app-sidebar'

export function AppWrapper() {
  return (
    <>
      <AppSidebar />
      <div className='flex-[1]'>
        <LoadingBar />
        <UpdatePrompt />
        <React.StrictMode>
          <Outlet />
        </React.StrictMode>
      </div>
    </>
  )
}

function LoadingBar() {
  const authActor = AuthMachineContext.useActorRef()
  const progress = useSelector(authActor, ({ context }) => context.progress)
  return (
    <div className='fixed top-[53px] h-1 w-full'>
      <div
        className='h-1 bg-yellow-500'
        style={{
          width: `${progress}%`,
          transition: progress > 10 ? 'width 3s ease-out' : 'width .25 ease-in'
        }}
      />
    </div>
  )
}

function UpdatePrompt() {
  const { t } = useTranslation()
  const [hasUpdate, setHasUpdate] = React.useState(false)
  const setError = useErrorPopup()

  React.useEffect(() => {
    CheckForUpdate().then(setHasUpdate).catch(setError)
  }, [])

  if (hasUpdate === false) {
    return null
  }

  return (
    <a
      className={cn(
        'group absolute bottom-2 left-0 z-50',
        'cursor-pointer text-base leading-5',
        'bg-[rgba(0,0,0,.625)] text-[#bfbcff] backdrop-blur transition-colors hover:bg-[rgba(0,0,0,.525)] hover:text-white',
        'ml-2 rounded-lg px-3 py-2'
      )}
      onClick={() => {
        BrowserOpenURL('https://cfn.williamsjokvist.se/')
        setHasUpdate(false)
      }}
    >
      <Icon
        icon='radix-icons:update'
        className='mr-2 inline h-4 w-4 text-[#49b3f5] transition-colors group-hover:text-white'
      />
      {t('newVersionAvailable')}
    </a>
  )
}
