import React from "react"

import type { model } from "@@/go/models"
import { GetGuiConfig } from "@@/go/core/CommandHandler"
import { useErrorMessage } from "./app-layout/error-message"
import { useTranslation } from "react-i18next"

type ConfigContextType = [cfg: model.GuiConfig | null, setConfig: (cfg: model.GuiConfig | null) => void]

const ConfigContext = React.createContext<ConfigContextType | null>(null)

export const useConfig = () => React.useContext(ConfigContext)

export const ConfigProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
  const { i18n } = useTranslation();
  const setErrorMessage = useErrorMessage();
  const [cfg, setCfg] = React.useState<model.GuiConfig | null>()
  
  React.useEffect(()=> {
    GetGuiConfig().then(cfg => {
      i18n.changeLanguage(cfg.locale)
      setCfg(cfg)
      console.log("App config:", cfg)
    }).catch(err => setErrorMessage(err))
  }, [])

  return (
    <ConfigContext.Provider value={[cfg, setCfg]}>
      {children}
    </ConfigContext.Provider>
  )
}
