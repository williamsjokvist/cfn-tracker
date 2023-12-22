import "i18next";
import type { locales } from "@@/go/models";

declare module "i18next" {
  interface CustomTypeOptions {
    // We don't have any namespaces
    defaultNS: "",
    nsSeparator: "",
    resources: {
      '': locales.Localization,
    };
  }
}
