import React from 'react'
import { Icon } from '@iconify/react'
import { useTranslation } from 'react-i18next'
import * as DialogPrimitive from '@radix-ui/react-dialog'

import type { LocalizationKey } from '@/main/i18n'
import { cn } from '@/helpers/cn'

type DialogContentProps = {
  title: LocalizationKey
  description?: LocalizationKey
} & React.PropsWithChildren<DialogPrimitive.DialogContentProps>
export const Content = React.forwardRef<HTMLDivElement, DialogContentProps>((props, ref) => {
  const { title, className, description, ...restProps } = props
  const { t } = useTranslation()
  return (
    <DialogPrimitive.Portal>
      <DialogPrimitive.Content
        ref={ref}
        className={cn(
          'fixed z-[9999] h-[420] w-[450px] p-4',
          'left-[50%] top-[50%] translate-x-[-50%] translate-y-[-50%]',
          'text-lg text-white',
          'bg-black bg-opacity-60 backdrop-blur-xl',
          'rounded-3xl shadow-[0_3px_16px_rgba(0,0,0,.5)]',
          className
        )}
        {...restProps}
      >
        <header className='flex justify-between'>
          <div>
            <DialogPrimitive.Title className='text-2xl font-bold'>{t(title)}</DialogPrimitive.Title>
            {description && (
              <DialogPrimitive.Description className='text-lg leading-none'>
                {t(description)}
              </DialogPrimitive.Description>
            )}
          </div>
          <DialogPrimitive.Close
            aria-label='Close'
            className={cn('h-11 w-11 rounded-full', 'bg-[#202020] hover:bg-[#2b2a33]')}
          >
            <Icon icon='ci:close-big' width={28} className='mx-auto' />
          </DialogPrimitive.Close>
        </header>
        {props.children}
      </DialogPrimitive.Content>
      <DialogPrimitive.Overlay
        className={cn(
          'grid place-items-center overflow-y-auto',
          'fixed bottom-0 left-0 right-0 top-0',
          'bg-black bg-opacity-20'
        )}
      />
    </DialogPrimitive.Portal>
  )
})

export const Root = DialogPrimitive.Root
export const Trigger = DialogPrimitive.Trigger
