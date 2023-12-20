import type { LocalizationKey } from "@/main/i18n-config";
import { AuthMachineContext } from "@/machines/auth-machine";
import { TrackingMachineContext } from "@/machines/tracking-machine";
import { errorsx } from "@@/go/models";
import { Icon } from "@iconify/react";
import { useSelector } from "@xstate/react";
import clsx from "clsx";
import { useAnimate } from "framer-motion";
import React from "react";
import { useTranslation } from "react-i18next";

const LocalizedErrorMessage: Record<number, LocalizationKey> = {
  401: "errUnauthorized",
  404: "errNotFound",
  500: "errInternalServerError",
};

export const ErrorMessage: React.FC= () => {
  const { t } = useTranslation();
  const [scope, animate] = useAnimate();
  
  const authActor = AuthMachineContext.useActorRef()
  const trackingActor = TrackingMachineContext.useActorRef()

  const authError = useSelector(authActor, ({ context }) => context.error)
  const trackingError = useSelector(trackingActor, ({ context }) => context.error)

  const [error, setError] = React.useState<errorsx.FrontEndError>(null)

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

  React.useEffect(() => {
    authError && setError(authError)
  }, [authError])

  React.useEffect(() => {
    trackingError && setError(trackingError)
  }, [trackingError])

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
        style={{ opacity: 0 }}
      >
        <Icon
          icon="material-symbols:warning-outline"
          className="text-[#ff6388] w-8 h-8 animate-blink-pulse"
        />
        {error && (
          <span>
            {LocalizedErrorMessage[error.code]
              ? t(LocalizedErrorMessage[error.code])
              : error.message}
          </span>
        )}
      </div>
    </div>
  );
};
