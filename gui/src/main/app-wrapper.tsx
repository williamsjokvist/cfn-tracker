import React from 'react'
import { useSelector } from '@xstate/react'
import { Outlet } from 'react-router-dom'

import { AuthMachineContext } from '@/state/auth-machine'

import { AppSidebar } from './app-sidebar'

export function AppWrapper() {
  return (
    <>
      <AppSidebar />
      <div className='flex-1'>
        <LoadingBar />
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
