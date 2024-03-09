import 'i18next'
import type { locales } from '@model'

declare module 'i18next' {
  interface CustomTypeOptions {
    // We don't have any namespaces
    defaultNS: ''
    nsSeparator: ''
    resources: {
      '': locales.Localization
    }
  }
}
