import { motion } from 'framer-motion'

import { cn } from '@/helpers/cn'

export function Root(props: React.PropsWithChildren) {
  return (
    <main className='relative z-40 grid h-screen w-full grid-rows-[0fr_1fr] text-white'>
      {props.children}
    </main>
  )
}

export function Header(props: React.PropsWithChildren) {
  return (
    <header
      className={cn(
        'flex items-center justify-between',
        'h-[53px] select-none px-8',
        'border-b-[1px] border-solid border-b-[rgba(255,255,255,.125)]'
      )}
      style={{ '--draggable': 'drag' } as React.CSSProperties}
    >
      {props.children}
    </header>
  )
}

export function Title(props: React.PropsWithChildren) {
  return (
    <motion.h2
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.125 }}
      className='whitespace-nowrap text-sm uppercase tracking-widest'
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
      className='inline-block h-5 w-5 animate-spin rounded-full border-[3px] border-current border-t-transparent text-pink-600'
      initial={{ opacity: 0 }}
      role='status'
      transition={{ delay: 0.125 }}
    />
  )
}
