import { motion } from 'framer-motion'

import { cn } from '@/helpers/cn'

export function Page(props: React.PropsWithChildren) {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.25 }}
      className='mx-4 mt-3 h-full max-h-[340px] overflow-y-scroll px-4 pb-4'
    >
      {props.children}
    </motion.div>
  )
}

export function Content(props: React.PropsWithChildren) {
  return <table className='w-full border-separate border-spacing-y-1'>{props.children}</table>
}

export function Tr(props: React.PropsWithChildren<React.TdHTMLAttributes<HTMLTableRowElement>>) {
  const { className, ...restProps } = props
  return (
    <tr
      className={cn(
        '[&>*:first-child]:rounded-l-xl [&>*:last-child]:rounded-r-xl',
        className
      )}
      {...restProps}
    >
      {props.children}
    </tr>
  )
}

export function Th(props: React.PropsWithChildren<React.ThHTMLAttributes<HTMLTableCellElement>>) {
  const { className, ...restProps } = props
  return (
    <th className={cn('whitespace-nowrap px-3 text-left', className)} {...restProps}>
      {props.children}
    </th>
  )
}

export function Td(
  props: React.PropsWithChildren<React.TdHTMLAttributes<HTMLTableCellElement>> & {
    interactive?: boolean
  }
) {
  const { className, interactive, ...restProps } = props
  return (
    <td
      className={cn(
        'whitespace-nowrap px-3 py-2 backdrop-blur-sm',
        'bg-white bg-opacity-5 group-hover:bg-opacity-10',
        'transition-colors',
        interactive && 'cursor-pointer hover:bg-opacity-20 active:bg-opacity-30',
        className
      )}
      {...restProps}
    >
      {props.children}
    </td>
  )
}
