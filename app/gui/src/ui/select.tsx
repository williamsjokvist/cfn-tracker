import React from 'react'
import * as SelectPrimitive from '@radix-ui/react-select'
import { CheckIcon, ChevronDownIcon, ChevronUpIcon } from '@radix-ui/react-icons'

import { cn } from '@/helpers/cn'

type Size = 'md' | 'sm'

// Items render in a portal, so they get the size from context rather than a prop.
const SizeContext = React.createContext<Size>('md')

const triggerClassName = (size: Size) =>
  cn(
    'inline-flex w-48 items-center justify-between rounded-lg border-[1px] font-normal',
    size === 'sm' ? 'px-3 py-1.5 text-base' : 'px-4 py-2 text-lg'
  )

export function Root(
  props: SelectPrimitive.SelectProps & { size?: Size; className?: string; loading?: boolean }
) {
  const { children, size = 'md', className, loading = false, ...restProps } = props
  const icon = <ChevronDownIcon className={size === 'sm' ? 'h-5 w-5' : 'h-6 w-6'} />
  if (loading) {
    // Laid out like the trigger, with its content hidden, so nothing moves once options load
    return (
      <div
        aria-busy
        className={cn(
          triggerClassName(size),
          'animate-pulse border-white/20 bg-white/10',
          className
        )}
      >
        <span className='invisible w-40'>&nbsp;</span>
        <span className='invisible'>{icon}</span>
      </div>
    )
  }
  return (
    <SelectPrimitive.Root {...restProps}>
      <SelectPrimitive.Trigger
        className={cn(
          triggerClassName(size),
          'transition-colors',
          'bg-white/5 hover:bg-white/10 active:bg-white/20',
          'border-white/20',
          'data-[state=open]:[&>*:last-child]:rotate-180',
          className
        )}
      >
        {/* min-w-0 lets a long value truncate instead of pushing past the trigger */}
        <span className='min-w-0 flex-1 overflow-hidden text-start'>
          <SelectPrimitive.Value />
        </span>
        <SelectPrimitive.Icon className='shrink-0 transition-transform'>
          {icon}
        </SelectPrimitive.Icon>
      </SelectPrimitive.Trigger>
      <SelectPrimitive.Portal>
        {/* CSS animations rather than framer-motion, since Radix only waits for CSS ones before unmounting */}
        <SelectPrimitive.Content
          position='popper'
          side='bottom'
          sideOffset={5}
          collisionPadding={16}
          className={cn(
            'data-[state=closed]:animate-select-out data-[state=open]:animate-select-in',
            // Cap to the space Radix measured so the viewport scrolls instead of overflowing the window.
            'z-50 flex max-h-[var(--radix-select-content-available-height)] flex-col overflow-hidden text-white',
            'max-w-[var(--radix-select-content-available-width)] min-w-[var(--radix-select-trigger-width)]',
            'bg-black/60 backdrop-blur-xl',
            'rounded-xl shadow-[0_3px_16px_rgba(0,0,0,.5)]'
          )}
        >
          <SelectPrimitive.ScrollUpButton className='flex justify-center py-1'>
            <ChevronUpIcon className='h-5 w-5' />
          </SelectPrimitive.ScrollUpButton>
          <SelectPrimitive.Viewport>
            <SizeContext.Provider value={size}>{children}</SizeContext.Provider>
          </SelectPrimitive.Viewport>
          <SelectPrimitive.ScrollDownButton className='flex justify-center py-1'>
            <ChevronDownIcon className='h-5 w-5' />
          </SelectPrimitive.ScrollDownButton>
        </SelectPrimitive.Content>
      </SelectPrimitive.Portal>
    </SelectPrimitive.Root>
  )
}

export function Item(props: SelectPrimitive.SelectItemProps & React.RefAttributes<HTMLDivElement>) {
  const { className, children, ...restProps } = props
  const size = React.useContext(SizeContext)
  return (
    <SelectPrimitive.Item
      className={cn(
        'flex min-w-48 cursor-pointer justify-between gap-4',
        size === 'sm' ? 'px-3 py-2 text-base' : 'px-4 py-3 text-lg',
        'outline-none data-[highlighted]:bg-[rgba(255,255,255,0.075)] data-[state=checked]:bg-[rgba(255,255,255,0.075)]'
      )}
      {...restProps}
    >
      <SelectPrimitive.ItemText className='min-w-0'>
        <div className={className}>{children}</div>
      </SelectPrimitive.ItemText>
      <SelectPrimitive.ItemIndicator>
        <CheckIcon className={cn('inline shrink-0', size === 'sm' ? 'h-5 w-5' : 'h-6 w-6')} />
      </SelectPrimitive.ItemIndicator>
    </SelectPrimitive.Item>
  )
}
