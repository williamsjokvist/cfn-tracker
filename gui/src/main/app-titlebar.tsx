import { Icon } from '@iconify/react'
import { Quit, WindowMinimise } from '@@/runtime/runtime'

export function AppTitleBar() {
  return (
    <div style={{ '--draggable': 'drag' } as React.CSSProperties}>
      <div className='flex justify-start'>
        <div className='group group mb-3 ml-2 mt-2 flex'>
          <button
            aria-label='close'
            className='mr-[8px] flex h-[14px] w-[14px] items-center justify-center rounded-full bg-slate-600 p-[2px] transition-all group-hover:bg-red-500'
            onClick={Quit}
          >
            <Icon
              icon='ep:close-bold'
              className='text-red-800 opacity-0 transition-all group-hover:opacity-100'
            />
          </button>
          <button
            aria-label='close'
            className='flex h-[14px] w-[14px] items-center justify-center rounded-full bg-slate-600 p-[2px] transition-all group-hover:bg-yellow-500'
            onClick={WindowMinimise}
          >
            <Icon
              icon='mingcute:minimize-fill'
              className='text-yellow-800 opacity-0 transition-all group-hover:opacity-100'
            />
          </button>
        </div>
      </div>
    </div>
  )
}
