import React from 'react'
import { CheckIcon, ChevronDownIcon, ChevronUpIcon } from '@radix-ui/react-icons'

import * as SelectPrimitive from '@radix-ui/react-select'

import { AnimatePresence, motion } from 'framer-motion'
import { cn } from '@/helpers/cn'

export const Root = React.forwardRef<HTMLButtonElement, SelectPrimitive.SelectProps>(
  ({ children, ...props }, forwardedRef) => {
    return (
      <SelectPrimitive.Root {...props}>
        <SelectPrimitive.Trigger
          className={cn(
            'transition-colors',
            'bg-white bg-opacity-5 hover:bg-opacity-20 active:bg-opacity-30',
            'rounded-2xl border-[1px] border-white border-opacity-20',
            'inline-flex items-center justify-between rounded-lg text-lg',
            'w-48 px-4 py-2 font-normal',
            '[&>*:last-child]:data-[state=open]:rotate-180'
          )}
          ref={forwardedRef}
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
                  'bg-black bg-opacity-60 backdrop-blur-xl',
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
)

export const Item = React.forwardRef<HTMLDivElement, SelectPrimitive.SelectItemProps>(
  ({ children, ...props }, forwardedRef) => {
    const { className, ...restProps } = props
    return (
      <SelectPrimitive.Item
        className={cn(
          'text-lg flex w-48 cursor-pointer justify-between px-4 py-3',
          'hover:bg-[rgba(255,255,255,0.075)] data-[state=checked]:bg-[rgba(255,255,255,0.075)]'
        )}
        ref={forwardedRef}
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
)
