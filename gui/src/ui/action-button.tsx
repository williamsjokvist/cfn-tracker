import React from "react";
import clsx from "clsx";

export const ActionButton = React.forwardRef((
  props: React.PropsWithChildren<React.ButtonHTMLAttributes<HTMLButtonElement>>, 
  ref: React.MutableRefObject<HTMLButtonElement>
) => (
  <button
    ref={ref}
    {...props}
    className={clsx(
      "flex items-center justify-between",
      "whitespace-nowrap  font-semibold text-md",
      "rounded-[18px] px-5 py-3 border-[1px]",
      "bg-[rgba(255,10,10,.1)] border-[#FF3D51] hover:bg-[#FF3D51] active:bg-[#ff6474] transition-colors",
      props.className
    )}
    {...props.disabled && {
      style: { filter: "saturate(0)" }
    }}
  >
    {props.children}
  </button>
));
