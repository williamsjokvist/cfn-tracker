import { cn } from '@/helpers/cn'
import eye from './enth-eye.jpg'

export function EnthEye(props: React.ImgHTMLAttributes<HTMLImageElement>) {
  const { className, ...restProps } = props
  return (
    <div className='absolute top-56 -right-20 -z-50'>
      <span
        className={cn(
          'relative inline-block',
          'rounded-full',
          'w-[624px] origin-[bottom_left] -rotate-90'
        )}
        style={{
          filter: 'opacity(40%) sepia(100%) hue-rotate(200deg)'
        }}
      >
        <img
          src={eye}
          alt='Enth'
          style={{
            filter: 'brightness(20%)',
            WebkitMaskImage: `url('data:image/svg+xml;utf8,<svg preserveAspectRatio="none" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0" y="0" viewBox="0 0 100 100" xml:space="preserve"><style type="text/css">.blur{filter:url(%23softedge);}</style><filter id="softedge"><feGaussianBlur stdDeviation="5"></feGaussianBlur></filter><g class="blur"><rect x="10" y="10" width="87" height="70"/></g></svg>')`,
            WebkitMaskSize: 'cover'
          }}
          className={cn(
            'opacity-enth rounded-full transition-opacity',
            'pointer-events-none -z-10'
          )}
          {...restProps}
        />
      </span>
    </div>
  )
}
