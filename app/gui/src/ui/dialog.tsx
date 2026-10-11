import { Icon } from '@iconify/react'
import { useTranslation } from 'react-i18next'
import * as DialogPrimitive from '@radix-ui/react-dialog'

import type { LocalizationKey } from '@/main/i18n'
import { cn } from '@/helpers/cn'

type DialogContentProps = {
  title: LocalizationKey
  description?: LocalizationKey
  dismissible?: boolean
} & DialogPrimitive.DialogContentProps

export const Root = DialogPrimitive.Root
export const Trigger = DialogPrimitive.Trigger

export function Content(props: DialogContentProps) {
  const { title, description, dismissible = true, className, children, ...restProps } = props
  const { t } = useTranslation()
  const preventUnlessDismissible = (e: Event) => !dismissible && e.preventDefault()
  return (
    <DialogPrimitive.Portal>
      <DialogPrimitive.Content
        asChild
        {...restProps}
        onEscapeKeyDown={e => {
          preventUnlessDismissible(e)
          restProps.onEscapeKeyDown?.(e)
        }}
        onInteractOutside={e => {
          preventUnlessDismissible(e)
          restProps.onInteractOutside?.(e)
        }}
      >
        {/* CSS animations rather than framer-motion, since Radix only waits for CSS ones before unmounting */}
        <div
          className={cn(
            'fixed z-9999 h-[420px] w-[450px] p-4',
            'data-[state=closed]:animate-dialog-out data-[state=open]:animate-dialog-in',
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
            {dismissible && (
              <DialogPrimitive.Close
                aria-label='Close'
                className={cn('h-11 min-w-[44px] rounded-full', 'bg-[#202020] hover:bg-[#2b2a33]')}
              >
                <Icon icon='ci:close-big' width={28} className='mx-auto' />
              </DialogPrimitive.Close>
            )}
          </header>

          {children}
        </div>
      </DialogPrimitive.Content>
      <DialogPrimitive.Overlay asChild>
        <div
          className={cn(
            'grid place-items-center overflow-y-auto',
            'data-[state=closed]:animate-fade-out data-[state=open]:animate-fade-in',
            'fixed top-0 right-0 bottom-0 left-0',
            'z-50 bg-black/40 backdrop-blur-xs'
          )}
        />
      </DialogPrimitive.Overlay>
    </DialogPrimitive.Portal>
  )
}
