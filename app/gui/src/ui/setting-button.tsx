import type React from 'react'
import { Icon } from '@iconify/react'

import { cn } from '@/helpers/cn'

export function SettingButton(
  props: React.PropsWithChildren<{ icon: string; disabled?: boolean; onClick: () => void }>
) {
  return (
    <button
      disabled={props.disabled}
      onClick={props.onClick}
      className={cn(
        'flex items-center gap-2 rounded-lg px-3 py-1.5 text-base',
        'border-[1px] border-white/20 bg-white/5 transition-colors',
        'hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-50'
      )}
    >
      <Icon icon={props.icon} className='text-icon h-5 w-5' />
      {props.children}
    </button>
  )
}
