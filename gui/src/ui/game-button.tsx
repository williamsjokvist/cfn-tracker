import React from 'react'

type GameButtonProps = {
  logo: string;
  code: string;
  alt: string;
} & React.HTMLAttributes<HTMLButtonElement>;

export const GameButton = React.forwardRef((
  props: GameButtonProps, 
  ref: React.MutableRefObject<HTMLButtonElement>
) => (
  <button
    ref={ref}
    {...props}
    type="button"
    className="w-52 p-3 rounded-lg hover:bg-slate-50 hover:bg-opacity-5 transition-colors"
  >
    <img
      src={props.logo}
      alt={props.alt}
      className="pointer-events-none select-none"
    />
  </button>
));
