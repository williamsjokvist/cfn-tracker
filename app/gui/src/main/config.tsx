import React from 'react'
import { useTranslation } from 'react-i18next'

import { model } from '@model'
import { GetGuiConfig } from '@cmd/CommandHandler'
import { AuthMachineContext } from '@/state/auth-machine'

const initialConfig: model.GUIConfig = {
  locale: 'en-GB',
  theme: model.ThemeName.DEFAULT,
  sidebar: false,
  logFile: true,
  browserSourcePort: 4242,
  lastGame: '' as model.GameType,
  lastUsers: {}
}

export const ConfigContext = React.createContext<
  [cfg: model.GUIConfig, setCfg: React.Dispatch<model.GUIConfig>]
>([
  initialConfig,
  () => {
    return
  }
])

export function ConfigProvider({ children }: React.PropsWithChildren) {
  const { i18n } = useTranslation()
  const [cfg, setCfg] = React.useState(initialConfig)
  const authActor = AuthMachineContext.useActorRef()

  React.useEffect(() => {
    ;(async function () {
      const cfg = await GetGuiConfig()
      setCfg(cfg)
      i18n.changeLanguage(cfg.locale)
      document.body.setAttribute('data-theme', cfg.theme)
      if (cfg.lastGame) {
        authActor.send({ type: 'submit', game: cfg.lastGame })
      }
    })()
  }, [])

  return <ConfigContext.Provider value={[cfg, setCfg]}>{children}</ConfigContext.Provider>
}
