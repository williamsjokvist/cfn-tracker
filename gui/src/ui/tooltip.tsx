import React from 'react'
import { motion } from 'framer-motion'

export function Tooltip(props: React.PropsWithChildren & {
  text: string
  disabled?: boolean
}) {
  return (
    <div
      className="group z-[9999] relative inline-flex justify-center"
    >
      {props.children}
      {!props.disabled &&
        <motion.span
          key='tooltip-text'
          className="select-none pointer-events-none z-50 group-hover:opacity-100 group-hover:visible opacity-0 invisible transition-all top-[-33px] absolute rounded-full border text-black bg-white px-2.5 pt-1 pb-0.5 text-sm font-semibold">
          {props.text}
        </motion.span>
      }
    </div>
  )
}
