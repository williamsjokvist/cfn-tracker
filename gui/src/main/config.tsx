import React from "react"
import { useTranslation } from "react-i18next"

import type { model } from "@@/go/models"
import { GetGuiConfig } from "@@/go/core/CommandHandler"

import { useErrorMessage } from "./app-layout/error-message"

type ConfigContextType = [cfg: model.GuiConfig | null, setCfg: (cfg: model.GuiConfig | null) => void]

const ConfigContext = React.createContext<ConfigContextType | null>(null)

export const useConfig = () => React.useContext(ConfigContext)

export const ConfigProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
  const { i18n } = useTranslation();
  const setErrorMessage = useErrorMessage();
  const configState = React.useState<model.GuiConfig | null>(null)
  const [_, setCfg] = configState;

  React.useEffect(()=> {
    GetGuiConfig().then(cfg => {
      i18n.changeLanguage(cfg.locale)
      setCfg(cfg)
      console.log("App config:", cfg)
    }).catch(err => setErrorMessage(err))
  }, [])

  return (
    <ConfigContext.Provider value={configState}>
      {children}
    </ConfigContext.Provider>
  )
}
