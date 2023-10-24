import { Icon } from "@iconify/react";
import clsx from "clsx";
import { useAnimate } from "framer-motion";
import React from "react";

type ErrorMessageProps = {
  errorMessage?: string;
  onFadedOut: () => void;
};

export const ErrorMessage: React.FC<ErrorMessageProps> = ({
  errorMessage,
  onFadedOut,
}) => {
  const [scope, animate] = useAnimate();

  React.useEffect(() => {
    if (!errorMessage) {
      return;
    }

    animate("#error-message", { opacity: [0, 1] }).then(() => {
      animate("#error-message", { opacity: [1, 0] }, { delay: 3.5 }).then(() =>
        onFadedOut()
      );
    });
  }, [errorMessage]);

  return (
    <div ref={scope} className="absolute">
      <div
        id="error-message"
        className={clsx(
          "flex gap-6 items-center justify-around",
          "fixed bottom-2 z-50",
          "px-8 py-3 rounded-r-xl text-xl backdrop-blur-sm pointer-events-none",
          "bg-[rgba(255,0,0,.125)]"
        )}
        {...(errorMessage === null && {
          style: {
            opacity: 0,
          },
        })}
      >
        <Icon
          icon="material-symbols:warning-outline"
          className="text-[#ff6388] w-8 h-8 blink-pulse"
        />
        <span>{errorMessage}</span>
      </div>
    </div>
  );
};
