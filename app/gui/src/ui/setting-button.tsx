import type { JSX } from 'react'
import { Icon } from '@iconify/react'

import { cn } from '@/helpers/cn'

export const settingTextClass = 'text-base font-extralight text-[#d6d4ff]'

// Spreads the rest of the props so it also works as a Radix trigger with asChild.
export function SettingButton(props: JSX.IntrinsicElements['button'] & { icon: string }) {
  const { icon, className, children, ...restProps } = props
  return (
    <button
      {...restProps}
      className={cn(
        'flex items-center gap-2 rounded-lg px-3 py-1.5 text-base',
        'border-[1px] border-white/20 bg-white/5 transition-colors',
        'hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-50',
        className
      )}
    >
      <Icon icon={icon} className='text-icon h-5 w-5' />
      {children}
    </button>
  )
}
