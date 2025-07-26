import * as HoverCardPrimitive from '@radix-ui/react-hover-card'
import { motion } from 'framer-motion'

export const Root = HoverCardPrimitive.Root
export const Trigger = HoverCardPrimitive.Trigger

export function Content(props: HoverCardPrimitive.HoverCardContentProps) {
  const { className, children, align = 'center', sideOffset = 4, ...restProps } = props
  return (
    <HoverCardPrimitive.Content
      className='z-50'
      align={align}
      sideOffset={sideOffset}
      {...restProps}
    >
      <motion.div
        key='modal'
        initial={{ y: 15, opacity: 0 }}
        animate={{
          y: 5,
          opacity: 1,
          transition: { type: 'tween', duration: 0.3 }
        }}
        exit={{ y: 50, opacity: 0, transition: { duration: 0.1 } }}
        className={className}
      >
        {children}
      </motion.div>
    </HoverCardPrimitive.Content>
  )
}
