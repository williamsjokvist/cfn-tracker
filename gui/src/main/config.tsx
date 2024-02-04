import React from 'react'
import { useTranslation } from 'react-i18next'

import type { model } from '@@/go/models'
import { GetGuiConfig } from '@@/go/core/CommandHandler'

const initialConfig: model.GuiConfig = {
  locale: 'en-GB',
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
    })()
  }, [])

  return <ConfigContext.Provider value={[cfg, setCfg]}>{children}</ConfigContext.Provider>
}
