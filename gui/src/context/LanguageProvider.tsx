import React from "react"
import { useTranslation } from "react-i18next"
import { APP_LANGUAGES } from "@/i18n"

export const LanguageProvider: React.FC<React.PropsWithChildren> = ( { children } ) => {
  const { i18n } = useTranslation();
  
  React.useEffect(() => {
    const storedLang = window.localStorage.getItem('lng')
    if (APP_LANGUAGES.map(lng => lng.code).includes(storedLang))
      i18n.changeLanguage(storedLang)
  }, [])

  return <>{children}</>
}