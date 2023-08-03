import React from "react";
import { Icon } from "@iconify/react";
import clsx from "clsx";

type ErrorMessageProps = {
  message?: string;
};

export const ErrorMessage: React.FC<ErrorMessageProps> = ({ message }) => {
  const [isOpen, setOpen] = React.useState(false);

  React.useEffect(() => {
    if (message) {
      setOpen(true);
      setTimeout(() => setOpen(false), 3500);
    }
  }, [message]);

  return (
    <div
      className={clsx(
        "flex gap-8 items-center justify-around",
        "fixed right-0 bottom-2 z-50",
        "px-8 py-3 rounded-l-2xl text-xl backdrop-blur-sm pointer-events-none",
        "bg-[rgba(255,0,0,.125)] transition-opacity",
        `${isOpen ? `opacity-100` : `opacity-0`}`
      )}
    >
      <Icon
        icon="material-symbols:warning-outline"
        className="text-[#ff6388] w-8 h-8 blink-pulse"
      />
      <span>{message}</span>
    </div>
  );
};
