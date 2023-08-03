type HideableTextProps = {
  text: string;
  hide: boolean;
} & React.ButtonHTMLAttributes<HTMLSpanElement>;

export const HideableText: React.FC<HideableTextProps> = ({
  text,
  hide,
  className,
}) => (
  <span
    className={`transition-opacity ${
      hide ? "opacity-0 w-0 overflow-hidden" : "opacity-100"
    } ${className}`}
  >
    {text}
  </span>
);
