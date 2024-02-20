import { cn } from '@/helpers/cn'
import eye from './enth-eye.jpg'

export function EnthEye(props: React.ImgHTMLAttributes<HTMLImageElement>) {
  const { className, ...restProps } = props
  return (
    <img
      src={eye}
      alt='Enth'
      className={cn('opacity-enth transition-opacity', 'pointer-events-none -z-10', className)}
      {...restProps}
    />
  )
}
