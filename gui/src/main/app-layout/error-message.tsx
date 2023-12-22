import React from "react";
import clsx from "clsx";
import { Icon } from "@iconify/react";
import { useAnimate } from "framer-motion";
import { useTranslation } from "react-i18next";

import type { LocalizationKey } from "@/main/i18n";
import type { errorsx } from "@@/go/models";

type ErrorContextType = [error: errorsx.FrontEndError | null, setError: (error: errorsx.FrontEndError | null) => void]
const ErrorContext = React.createContext<ErrorContextType | null>(null)
export const useErrorMessage = () => React.useContext(ErrorContext)![1]

export const LocalizedErrorMessage: Record<number, LocalizationKey> = {
  401: "errUnauthorized",
  404: "errNotFound",
  500: "errInternalServerError",
};

export const ErrorMessageProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
  const { t } = useTranslation();
  const [scope, animate] = useAnimate();
  const [error, setError] = React.useState<errorsx.FrontEndError | null>(null)

  React.useEffect(() => {
    if (error === null) {
      return;
    }
    animate("#error-message", { opacity: [0, 1] }).then(() => {
      animate("#error-message", { opacity: [1, 0] }, { delay: 3.5 }).then(() =>
        setError(null)
      );
    });
  }, [error]);

  return (
    <ErrorContext.Provider value={[error, setError]}>
      <div ref={scope} className="absolute w-full flex justify-end">
        <div
          id="error-message"
          className={clsx(
            "flex gap-6 items-center justify-around",
            "fixed mx-auto z-50",
            "px-8 py-3 rounded-bl-xl text-xl backdrop-blur-sm pointer-events-none",
            "text-white bg-gradient-to-r from-[#870e65] to-[#6c086d]"
          )}
          style={{ opacity: 0 }}
        >
          <Icon
            icon="material-symbols:warning-outline"
            className="text-[#ff6388] w-8 h-8 animate-blink-pulse"
          />
          {error && (
            <span>
              {error?.code && LocalizedErrorMessage[error.code]
                ? t(LocalizedErrorMessage[error.code])
                : error.message}
            </span>
          )}
        </div>
      </div>
      {children}
    </ErrorContext.Provider>
  );
};
