import React from 'react'
import { AnimatePresence, motion } from 'framer-motion'
import * as SelectPrimitive from '@radix-ui/react-select'
import { CheckIcon, ChevronDownIcon, ChevronUpIcon } from '@radix-ui/react-icons'

import { cn } from '@/helpers/cn'

export function Root(props: SelectPrimitive.SelectProps) {
  const { children, ...restProps } = props
  return (
    <SelectPrimitive.Root {...restProps}>
      <SelectPrimitive.Trigger
        className={cn(
          'transition-colors',
          'bg-white/5 hover:bg-white/10 active:bg-white/20',
          'rounded-2xl border-[1px] border-white/20',
          'inline-flex items-center justify-between rounded-lg text-lg',
          'w-48 px-4 py-2 font-normal',
          'data-[state=open]:[&>*:last-child]:rotate-180'
        )}
      >
        <SelectPrimitive.Value />
        <SelectPrimitive.Icon className='transition-transform'>
          <ChevronDownIcon className='h-6 w-6' />
        </SelectPrimitive.Icon>
      </SelectPrimitive.Trigger>
      <AnimatePresence>
        <SelectPrimitive.Portal>
          <SelectPrimitive.Content position='popper' side='bottom' className='animate- z-50'>
            <motion.div
              key='modal'
              initial={{ y: 15, opacity: 0 }}
              animate={{
                y: 5,
                opacity: 1,
                transition: { type: 'tween', duration: 0.3 }
              }}
              exit={{ y: 50, opacity: 0, transition: { duration: 0.1 } }}
              className={cn(
                'z-50 overflow-hidden text-white',
                'bg-black/60 backdrop-blur-xl',
                'rounded-xl shadow-[0_3px_16px_rgba(0,0,0,.5)]'
              )}
            >
              <SelectPrimitive.ScrollUpButton>
                <ChevronUpIcon />
              </SelectPrimitive.ScrollUpButton>
              <SelectPrimitive.Viewport>{children}</SelectPrimitive.Viewport>
              <SelectPrimitive.ScrollDownButton>
                <ChevronDownIcon />
              </SelectPrimitive.ScrollDownButton>
            </motion.div>
          </SelectPrimitive.Content>
        </SelectPrimitive.Portal>
      </AnimatePresence>
    </SelectPrimitive.Root>
  )
}

export function Item(props: SelectPrimitive.SelectItemProps & React.RefAttributes<HTMLDivElement>) {
  const { className, children, ...restProps } = props
  return (
    <SelectPrimitive.Item
      className={cn(
        'flex w-48 cursor-pointer justify-between px-4 py-3 text-lg',
        'hover:bg-[rgba(255,255,255,0.075)] data-[state=checked]:bg-[rgba(255,255,255,0.075)]'
      )}
      {...restProps}
    >
      <SelectPrimitive.ItemText>
        <div className={className}>{children}</div>
      </SelectPrimitive.ItemText>
      <SelectPrimitive.ItemIndicator>
        <CheckIcon className='inline h-6 w-6' />
      </SelectPrimitive.ItemIndicator>
    </SelectPrimitive.Item>
  )
}
