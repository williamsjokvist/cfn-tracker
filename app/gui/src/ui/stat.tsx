import { cn } from '@/helpers/cn'

export function BigStat({
  text,
  value,
  className
}: {
  text: string
  value: string | number
  className?: string
}) {
  return (
    <div
      className={cn(
        'mb-2 flex flex-1 justify-between gap-4 rounded-xl bg-slate-50/5 p-3 pb-1',
        className
      )}
    >
      <dt className='font-extralight tracking-wider'>{text}</dt>
      <dd className='text-4xl font-semibold'>{value}</dd>
    </div>
  )
}
