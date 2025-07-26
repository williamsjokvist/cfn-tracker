import type { JSX } from 'react'

import { cn } from '@/helpers/cn'

export function Checkbox(props: JSX.IntrinsicElements['input']) {
  const { className, ...restProps } = props
  return (
    <input
      ref={props.ref}
      type='checkbox'
      className={cn(
        'h-7 w-7 cursor-pointer rounded-md bg-transparent text-transparent',
        'border-2 border-[rgba(255,255,255,.25)] focus:border-2',
        'checked:border-2 checked:border-[rgba(255,255,255,.25)] checked:hover:border-[rgba(255,255,255,.25)] checked:focus:border-[rgba(255,255,255,.25)]',
        'focus:ring-transparent focus:ring-offset-transparent',
        className
      )}
      {...restProps}
    />
  )
}
