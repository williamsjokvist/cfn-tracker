import React from 'react'

export const ActionButton: React.FC<React.PropsWithChildren<React.HTMLAttributes<HTMLButtonElement>>> = (
  { onClick, className, style,  children }
) => {
  return (
    <button 
      type='button'
      onClick={onClick}
      style={style}
      className={`
        flex
        items-center
        justify-between
        bg-[rgba(255,10,10,.1)]
        rounded-[18px]
        px-5 py-3 
        border-[1px] border-[#FF3D51] hover:bg-[#FF3D51] 
        whitespace-nowrap
        transition-colors font-semibold text-md ` + className
      }
    >
      {children}
    </button>
  )
}
