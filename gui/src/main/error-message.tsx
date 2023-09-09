import { Icon } from "@iconify/react"
import clsx from "clsx"
import React from "react"

const ErrorMessageContext = React.createContext<{
  errorMessage: string|null ,
  setErrorMessage: (message: string) => void
}>({
  errorMessage: null,
  setErrorMessage: () => {}
})

export const useErrorMessage = () => React.useContext(ErrorMessageContext)

export const ErrorMessageProvider: React.FC<React.PropsWithChildren> = ( { children }) => {
  const [_, setErrorMessage] = React.useState<string | null>(null)

  return (
    <ErrorMessageContext.Provider value={{ errorMessage: null, setErrorMessage }}>
      {children}
    </ErrorMessageContext.Provider>
  )
}

export const ErrorMessage: React.FC = () => {
  const [isOpen, setOpen] = React.useState(false);
  const { errorMessage } = useErrorMessage()

  React.useEffect(() => {
    if (errorMessage) {
      setOpen(true);
      setTimeout(() => setOpen(false), 3500);
    }
  }, [errorMessage]);

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
      <span>{errorMessage}</span>
    </div>
  );
};
