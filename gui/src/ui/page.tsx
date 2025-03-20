import { motion } from 'framer-motion'

import { cn } from '@/helpers/cn'

import { EnthLogo } from './enth-logo'
import { EnthEye } from '@/ui/enth-eye'

export function Root(props: React.PropsWithChildren) {
  return (
    <main className='relative z-40 grid h-screen w-full grid-rows-[0fr_1fr] text-white'>
      {props.children}
      <EnthEye />
    </main>
  )
}

export function Header(props: React.PropsWithChildren) {
  return (
    <header
      className={cn(
        'relative flex items-center justify-between',
        'h-[53px] px-8 select-none',
        'border-b-divider border-b-[1px] border-solid'
      )}
      style={{ '--draggable': 'drag' } as React.CSSProperties}
    >
      {props.children}
      <EnthLogo
        className={cn('opacity-enth transition-opacity', 'absolute right-4 -z-10')}
        width={100}
        fill='rgba(149, 243, 246, 0.075)'
        strokeWidth='50'
      />
    </header>
  )
}

export function Title(props: React.PropsWithChildren) {
  return (
    <motion.h2
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.125 }}
      className='text-sm tracking-widest whitespace-nowrap uppercase'
    >
      {props.children}
    </motion.h2>
  )
}

export function LoadingIcon() {
  return (
    <motion.i
      animate={{ opacity: 1 }}
      aria-label='loading'
      className='text-highlight inline-block h-5 w-5 animate-spin rounded-full border-[3px] border-current border-t-transparent'
      initial={{ opacity: 0 }}
      role='status'
      transition={{ delay: 0.125 }}
    />
  )
}
