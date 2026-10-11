import {
  FR,
  JP,
  GB,
  ES,
  SA,
  PT,
  KR,
  CN,
  DE,
  IT,
  RU,
  TW,
  TH,
  PL,
  SE,
  IN,
  ID,
  TR,
  VN,
  UA,
  NL,
  Props
} from 'country-flag-icons/react/3x2'

const flagMap = {
  'fr-FR': FR,
  'ja-JP': JP,
  'en-GB': GB,
  'es-ES': ES,
  'ar-SA': SA,
  'pt-BR': PT,
  'ko-KR': KR,
  'zh-CN': CN,
  'de-DE': DE,
  'it-IT': IT,
  'ru-RU': RU,
  'zh-TW': TW,
  'th-TH': TH,
  'pl-PL': PL,
  'sv-SE': SE,
  'hi-IN': IN,
  'id-ID': ID,
  'tr-TR': TR,
  'vi-VN': VN,
  'uk-UA': UA,
  'nl-NL': NL
}

export function Flag(props: Props & { code: string }) {
  const { code, ...restProps } = props
  const Flag = flagMap[code]
  return <Flag {...restProps} />
}
