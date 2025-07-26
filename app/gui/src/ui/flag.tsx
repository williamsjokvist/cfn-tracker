import { FR, JP, GB, Props } from 'country-flag-icons/react/3x2'

const flagMap = {
  'fr-FR': FR,
  'ja-JP': JP,
  'en-GB': GB
}

export function Flag(props: Props & { code: string }) {
  const { code, ...restProps } = props
  const Flag = flagMap[code]
  return <Flag {...restProps} />
}
