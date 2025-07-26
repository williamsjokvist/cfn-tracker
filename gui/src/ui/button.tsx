import type { JSX } from 'react'

import { cn } from '@/helpers/cn'

export function Button(props: JSX.IntrinsicElements['button']) {
  const { disabled, style, className, children, onClick, ...restProps } = props
  return (
    <button
      onClick={disabled ? undefined : onClick}
      {...(disabled
        ? {
            style: {
              backgroundColor: 'rgba(0,0,0,.25)',
              border: '1px solid rgba(255,255,255,.25)',
              cursor: 'not-allowed',
              ...style
            }
          }
        : { style })}
      className={cn(
        'flex items-center justify-between',
        'text-md font-semibold whitespace-nowrap',
        'bg-[rgba(255,10,10,.1)] transition-colors',
        'hover:bg-[#FF3D51] active:bg-[#ff6474]',
        'rounded-[18px] border-[1px] border-[#FF3D51]',
        'px-5 py-3',
        className
      )}
      {...restProps}
    >
      {children}
    </button>
  )
}
