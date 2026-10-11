import { Input as InputPrimitive } from '@base-ui/react/input'

import { cn } from '@/helpers/cn'

// Replaces the native size attribute, which sets the width in characters.
export function Input(
  props: Omit<InputPrimitive.Props, 'size'> & { className?: string; size?: 'md' | 'sm' }
) {
  const { className, size = 'md', ...restProps } = props
  return (
    <InputPrimitive
      className={cn(
        'w-48 rounded-lg border-[1px] border-white/20 font-normal text-white',
        size === 'sm' ? 'px-3 py-1.5 text-base' : 'px-4 py-2 text-lg',
        'bg-white/5 transition-colors outline-none',
        'hover:bg-white/10 data-[focused]:border-white/50 data-[focused]:bg-white/10',
        'data-[disabled]:cursor-not-allowed data-[disabled]:opacity-50',
        // Hide the native number spinners, which don't fit the theme.
        '[appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none',
        className
      )}
      {...restProps}
    />
  )
}
