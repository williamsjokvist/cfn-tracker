import React from 'react'
import { useTranslation } from 'react-i18next'

import type { model } from '@model'
import { GetGuiConfig } from '@cmd/CommandHandler'

const initialConfig: model.GuiConfig = {
  locale: 'en-GB',
  theme: 'default',
  sidebarMinified: false
}

export const ConfigContext = React.createContext<
  [cfg: model.GuiConfig, setCfg: React.Dispatch<model.GuiConfig>]
>([
  initialConfig,
  () => {
    return
  }
])

export function ConfigProvider({ children }: React.PropsWithChildren) {
  const { i18n } = useTranslation()
  const [cfg, setCfg] = React.useState(initialConfig)

  React.useEffect(() => {
    ;(async function () {
      const cfg = await GetGuiConfig()
      setCfg(cfg)
      i18n.changeLanguage(cfg.locale)
      document.body.setAttribute('data-theme', cfg.theme)
    })()
  }, [])

  return <ConfigContext.Provider value={[cfg, setCfg]}>{children}</ConfigContext.Provider>
}
