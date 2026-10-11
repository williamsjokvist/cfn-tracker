import { FR, JP, GB, ES, SA, BR, KR, CN, DE, Props } from 'country-flag-icons/react/3x2'

const flagMap = {
  'fr-FR': FR,
  'ja-JP': JP,
  'en-GB': GB,
  'es-ES': ES,
  'ar-SA': SA,
  'pt-BR': BR,
  'ko-KR': KR,
  'zh-CN': CN,
  'de-DE': DE
}

export function Flag(props: Props & { code: string }) {
  const { code, ...restProps } = props
  const Flag = flagMap[code]
  return <Flag {...restProps} />
}
