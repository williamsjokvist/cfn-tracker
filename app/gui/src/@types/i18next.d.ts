import 'i18next'
import type { model } from '@model'

declare module 'i18next' {
  interface CustomTypeOptions {
    // We don't have any namespaces
    defaultNS: ''
    nsSeparator: ''
    resources: {
      '': model.Localization
    }
  }
}
