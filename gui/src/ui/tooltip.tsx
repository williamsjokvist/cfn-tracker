import React from 'react'
import { motion } from 'framer-motion'

export function Tooltip(
  props: React.PropsWithChildren & {
    text: string
    disabled?: boolean
  }
) {
  return (
    <div className='group relative z-[9999] inline-flex justify-center'>
      {props.children}
      {!props.disabled && (
        <motion.span
          key='tooltip-text'
          className='pointer-events-none invisible absolute top-[-33px] z-50 select-none rounded-full border bg-white px-2.5 pb-0.5 pt-1 text-sm font-semibold text-black opacity-0 transition-all group-hover:visible group-hover:opacity-100'
        >
          {props.text}
        </motion.span>
      )}
    </div>
  )
}
