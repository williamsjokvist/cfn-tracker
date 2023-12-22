import React from "react";
import { useRouteError } from "react-router-dom";
import { motion } from "framer-motion";
import { Icon } from "@iconify/react";

import { LocalizedErrorMessage } from "./error-message";
import { errorsx } from "@@/go/models";
import { PageHeader } from "@/ui/page-header";
import { useTranslation } from "react-i18next";
import { AppTitleBar } from "./app-titlebar";

const isAppError = (error: unknown) => error instanceof Object && 'message' in error && 'code' in error

type AppErrorBoundaryProps = {
  outer?: boolean
}
export const AppErrorBoundary: React.FC<AppErrorBoundaryProps> = ({ outer }) => {
  const { t } = useTranslation();
  const thrownError = useRouteError();
  const [err, setErr] = React.useState<errorsx.FrontEndError | null>(null)

  React.useEffect(() => {
    console.error(thrownError);
    if (thrownError instanceof Error) {
      setErr({ code: 500, message: thrownError.message });
    } else if (isAppError(thrownError)) {
      setErr({ code: 500, message: (thrownError as errorsx.AppError).message });
    }
  }, [thrownError])

  if (!err?.code) {
    return null
  }

  return (
    <motion.section
      className="text-white h-screen w-full"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.125 }}
    >
      {outer && <AppTitleBar />}
      {!outer && <PageHeader text={t(LocalizedErrorMessage[err.code])} />}
      <div className="mt-8 w-full text-center flex flex-col items-center justify-center pb-16 rounded-md">
        <Icon
          icon="material-symbols:warning-outline"
          className="text-[#ff6388] w-40 h-40"
        />
        <h1 className="text-2xl text-center font-bold">{t(LocalizedErrorMessage[err.code])}</h1>
        <p className="text-xl">{err.message}</p>
      </div>
    </motion.section>
  )
}
