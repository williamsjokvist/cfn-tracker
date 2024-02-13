import React from 'react'

import { cn } from '@/helpers/cn'

export const Checkbox = React.forwardRef<
  HTMLInputElement,
  React.InputHTMLAttributes<HTMLInputElement>
>((props, ref) => {
  const { className, ...restProps } = props
  return (
    <input
      ref={ref}
      type='checkbox'
      className={cn(
        'mr-4 h-7 w-7 cursor-pointer rounded-md bg-transparent text-transparent',
        'border-2 border-[rgba(255,255,255,.25)] focus:border-2',
        'checked:border-2 checked:border-[rgba(255,255,255,.25)] checked:hover:border-[rgba(255,255,255,.25)] checked:focus:border-[rgba(255,255,255,.25)]',
        'focus:ring-transparent focus:ring-offset-transparent',
        className
      )}
      {...restProps}
    />
  )
})
