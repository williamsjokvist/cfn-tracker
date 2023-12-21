import "i18next";
import type { locales } from "@@/go/models";

declare module "i18next" {
  interface CustomTypeOptions {
    defaultNS: "",
    resources: {
      '': locales.Localization,
    };
  }
}
