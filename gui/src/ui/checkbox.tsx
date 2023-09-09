import React from "react";
import clsx from "clsx";

export const Checkbox = React.forwardRef((
  props: React.InputHTMLAttributes<HTMLInputElement>,
  ref: React.MutableRefObject<HTMLInputElement>
) => (
  <input
    ref={ref}
    type="checkbox"
    className={clsx(
      "w-7 h-7 rounded-md cursor-pointer mr-4 bg-transparent text-transparent",
      "border-2 border-[rgba(255,255,255,.25)] focus:border-2",
      "checked:border-2 checked:border-[rgba(255,255,255,.25)] checked:focus:border-[rgba(255,255,255,.25)] checked:hover:border-[rgba(255,255,255,.25)]",
      "focus:ring-offset-transparent focus:ring-transparent"
    )}
    {...props}
  />
));
