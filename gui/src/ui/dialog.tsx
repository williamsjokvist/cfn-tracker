import React from 'react'
import { Icon } from '@iconify/react'
import { useTranslation } from 'react-i18next'
import * as DialogPrimitive from '@radix-ui/react-dialog'

import type { LocalizationKey } from '@/main/i18n'
import { cn } from '@/helpers/cn'
import { motion } from 'framer-motion'

type DialogContentProps = {
  title: LocalizationKey
  description?: LocalizationKey
} & React.PropsWithChildren<DialogPrimitive.DialogContentProps>
export const Content = React.forwardRef<HTMLDivElement, DialogContentProps>((props, ref) => {
  const { title, className, description, ...restProps } = props
  const { t } = useTranslation()
  return (
    <DialogPrimitive.Portal>
      <DialogPrimitive.Content asChild ref={ref} {...restProps}>
        <motion.div
          initial={{ opacity: 0, top: '55%' }}
          animate={{ opacity: 1, top: '50%' }}
          transition={{ delay: 0.025 }}
          className={cn(
            'fixed z-9999 h-[420px] w-[450px] p-4',
            'top-[50%] left-[50%] translate-x-[-50%] translate-y-[-50%]',
            'text-lg text-white',
            'bg-black/60 backdrop-blur-xl',
            'rounded-3xl shadow-[0_3px_16px_rgba(0,0,0,.5)]',
            className
          )}
        >
          <header className='flex justify-between'>
            <div>
              <DialogPrimitive.Title className='text-2xl font-bold'>
                {t(title)}
              </DialogPrimitive.Title>
              {description && (
                <DialogPrimitive.Description className='text-lg leading-none'>
                  {t(description)}
                </DialogPrimitive.Description>
              )}
            </div>
            <DialogPrimitive.Close
              aria-label='Close'
              className={cn('h-11 min-w-[44px] rounded-full', 'bg-[#202020] hover:bg-[#2b2a33]')}
            >
              <Icon icon='ci:close-big' width={28} className='mx-auto' />
            </DialogPrimitive.Close>
          </header>

          {props.children}
        </motion.div>
      </DialogPrimitive.Content>
      <DialogPrimitive.Overlay asChild>
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.025 }}
          className={cn(
            'grid place-items-center overflow-y-auto',
            'fixed top-0 right-0 bottom-0 left-0',
            'z-50 bg-black/40 backdrop-blur-xs'
          )}
        />
      </DialogPrimitive.Overlay>
    </DialogPrimitive.Portal>
  )
})

export const Root = DialogPrimitive.Root
export const Trigger = DialogPrimitive.Trigger
