import React from 'react'

import { GetSupportedLanguages } from '@cmd/CommandHandler'

import { useErrorPopup } from '@/main/error-popup'
import * as Select from '@/ui/select'
import { Flag } from '@/ui/flag'

export function LanguageSelect(props: { value?: string; onValueChange: (code: string) => void }) {
  const [langs, setLangs] = React.useState<string[]>([])
  const setError = useErrorPopup()

  React.useEffect(() => {
    GetSupportedLanguages().then(setLangs).catch(setError)
  }, [])

  return (
    <Select.Root size='sm' onValueChange={props.onValueChange} value={props.value}>
      {langs.map(code => {
        const locale = new Intl.Locale(code).maximize()
        const sharesLanguage =
          langs.filter(c => new Intl.Locale(c).language === locale.language).length > 1
        const nativeName = new Intl.DisplayNames([code], {
          type: 'language',
          languageDisplay: 'standard',
          style: 'narrow',
          fallback: 'code'
        }).of(sharesLanguage ? `${locale.language}-${locale.script}` : locale.language)!

        return (
          <Select.Item
            key={code}
            value={code}
            className='flex items-center justify-between gap-2'
            {...(props.value === code && {
              style: {
                fontWeight: 600
              }
            })}
          >
            <Flag code={code} className='w-5 rounded-sm' />
            <span className='first-letter:uppercase'>{nativeName}</span>
          </Select.Item>
        )
      })}
    </Select.Root>
  )
}
