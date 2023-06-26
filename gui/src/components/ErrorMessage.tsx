import React from 'react'
import { Icon } from '@iconify/react'

type ErrorMessageProps = {
  message: string
}
export const ErrorMessage: React.FC<ErrorMessageProps> = ( { message } ) => {
  const [isOpen, setOpen] = React.useState(false)

  React.useEffect(() => {
    if (message == '') return
    setOpen(true)
    setTimeout(() => setOpen(false), 3500)
  }, [message])

  return (
    <div className={`transition-opacity z-50 fixed right-0 bottom-2 flex justify-around px-8 py-3 gap-8 rounded-l-2xl text-xl items-center bg-[rgba(255,0,0,.125)] backdrop-blur-sm pointer-events-none ` + `${isOpen ? `opacity-100` : `opacity-0`}`}>
      <Icon icon='material-symbols:warning-outline' className='text-[#ff6388] w-8 h-8 blink-pulse' />
      <span>{message}</span>
    </div>
  )
}