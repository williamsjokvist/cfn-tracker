import * as SwitchPrimitive from '@radix-ui/react-switch'

import { cn } from '@/helpers/cn'

export function Switch(props: SwitchPrimitive.SwitchProps) {
  const { className, ...restProps } = props
  return (
    <SwitchPrimitive.Root
      className={cn(
        'relative h-7 w-12 cursor-pointer',
        'rounded-full bg-white/25 backdrop-blur-xl',
        className
      )}
      {...restProps}
    >
      <SwitchPrimitive.Thumb
        className={cn(
          'block h-5 w-5 rounded-full bg-white',
          'translate-x-1 transition-transform data-[state=checked]:translate-x-6'
        )}
      >
        {props.children}
      </SwitchPrimitive.Thumb>
    </SwitchPrimitive.Root>
  )
}
