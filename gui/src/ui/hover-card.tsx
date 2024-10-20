import * as React from 'react'
import * as HoverCardPrimitive from '@radix-ui/react-hover-card'
import { motion } from 'framer-motion'

import { cn } from '@/helpers/cn'

const HoverCard = HoverCardPrimitive.Root

const HoverCardTrigger = HoverCardPrimitive.Trigger

const HoverCardContent = React.forwardRef<
  React.ElementRef<typeof HoverCardPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof HoverCardPrimitive.Content>
>(({ className, align = 'center', sideOffset = 4, ...props }, ref) => (
  <HoverCardPrimitive.Content
    className='z-50'
    ref={ref}
    align={align}
    sideOffset={sideOffset}
    {...props}
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
      className={cn(
        'overflow-hidden p-4 text-white',
        'bg-black bg-opacity-90 backdrop-blur-xl',
        'rounded-xl shadow-[0_3px_16px_rgba(0,0,0,.5)]',
        className
      )}
    >
      {props.children}
    </motion.div>
  </HoverCardPrimitive.Content>
))
HoverCardContent.displayName = HoverCardPrimitive.Content.displayName

export { HoverCard, HoverCardTrigger, HoverCardContent }
