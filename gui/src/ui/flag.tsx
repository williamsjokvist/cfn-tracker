import { FR, JP, GB, type FlagComponent, Props } from 'country-flag-icons/react/3x2'

export function Flag(props: { code: string } & Props) {
  const { code, ...restProps } = props

  let Flag: FlagComponent
  switch (props.code) {
    case 'fr-FR':
      Flag = FR
      break
    case 'ja-JP':
      Flag = JP
      break
    case 'en-GB':
    default:
      Flag = GB
  }

  return <Flag {...restProps} />
}
