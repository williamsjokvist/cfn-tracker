import React from 'react'
import { AnimatePresence, motion } from 'framer-motion'
import { Icon } from '@iconify/react'

import { cn } from '@/helpers/cn'

const variants = {
  error: {
    icon: (
      <Icon
        icon='material-symbols:warning-outline'
        className='animate-blink h-8 w-8 text-[#ff6388]'
      />
    ),
    className:
      'end-0 top-0 rounded-es-xl bg-linear-to-r from-[#870e65] to-[#6c086d] backdrop-blur-xs'
  },
  loading: {
    icon: (
      <i className='text-highlight inline-block h-6 w-6 animate-spin rounded-full border-[3px] border-current border-t-transparent' />
    ),
    className:
      'start-1/2 bottom-10 -translate-x-1/2 rounded-xl bg-black/60 backdrop-blur-xl rtl:translate-x-1/2'
  }
}

export function Banner({
  open,
  variant,
  onClick,
  children
}: React.PropsWithChildren<{
  open: boolean
  variant: keyof typeof variants
  onClick?: () => void
}>) {
  const { icon, className } = variants[variant]
  return (
    <AnimatePresence>
      {open && (
        <motion.div
          role={variant === 'error' ? 'alert' : 'status'}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className={cn(
            'flex items-center justify-around gap-6',
            'fixed z-50 px-8 py-3 text-xl whitespace-nowrap',
            'text-white',
            className,
            onClick ? 'cursor-pointer' : 'pointer-events-none'
          )}
          onClick={onClick}
        >
          {icon}
          {children}
        </motion.div>
      )}
    </AnimatePresence>
  )
}
