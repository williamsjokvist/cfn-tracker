import { JSX } from 'react'
import { cn } from '@/helpers/cn'

export function Page(props: JSX.IntrinsicElements['div']) {
  const { className, ...restProps } = props
  return (
    <div
      className={cn('mx-4 mt-3 h-full max-h-[340px] overflow-y-scroll px-4 pb-4', className)}
      {...restProps}
    />
  )
}

export function Content(props: JSX.IntrinsicElements['table']) {
  const { className, ...restProps } = props
  return (
    <table className={cn('w-full border-separate border-spacing-y-1', className)} {...restProps} />
  )
}

export function Tr(props: JSX.IntrinsicElements['tr']) {
  const { className, ...restProps } = props
  return (
    <tr
      className={cn('[&>*:first-child]:rounded-l-xl [&>*:last-child]:rounded-r-xl', className)}
      {...restProps}
    />
  )
}

export function Th(props: JSX.IntrinsicElements['th']) {
  const { className, ...restProps } = props
  return <th className={cn('px-3 text-left whitespace-nowrap', className)} {...restProps} />
}

export function Td(
  props: JSX.IntrinsicElements['td'] & {
    interactive?: boolean
  }
) {
  const { className, interactive, ...restProps } = props
  return (
    <td
      className={cn(
        'px-3 py-2 whitespace-nowrap backdrop-blur-xs',
        'bg-white/5 group-hover:bg-white/10',
        'transition-colors',
        interactive ? 'cursor-pointer hover:bg-white/20 active:bg-white/30' : '',
        className
      )}
      {...restProps}
    />
  )
}
