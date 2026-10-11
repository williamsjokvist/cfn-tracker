import { Input as InputPrimitive } from '@base-ui/react/input'

import { cn } from '@/helpers/cn'

export function Input(props: InputPrimitive.Props & { className?: string }) {
  const { className, ...restProps } = props
  return (
    <InputPrimitive
      className={cn(
        'w-48 rounded-lg border-[1px] border-white/20 px-4 py-2 text-lg font-normal text-white',
        'bg-white/5 transition-colors outline-none',
        'hover:bg-white/10 data-[focused]:border-white/50 data-[focused]:bg-white/10',
        // Hide the native number spinners, which don't fit the theme.
        '[appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none',
        className
      )}
      {...restProps}
    />
  )
}
