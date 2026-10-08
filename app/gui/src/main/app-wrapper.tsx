import React from 'react'
import { Outlet } from 'react-router-dom'

import { AppSidebar } from './app-sidebar'

export function AppWrapper() {
  return (
    <>
      <AppSidebar />
      <div className='flex-1'>
        <React.StrictMode>
          <Outlet />
        </React.StrictMode>
      </div>
    </>
  )
}
