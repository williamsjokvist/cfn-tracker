import React from 'react'

export function CharacterPortrait({
  character,
  className
}: {
  character: string
  className?: string
}) {
  const src = `https://www.streetfighter.com/6/buckler/assets/images/material/character/character_${character
    .toLowerCase()
    .replace(/\s/g, '')
    .replace('.', '')}_r.png`
  const [failedSrc, setFailedSrc] = React.useState<string>()
  if (!character || failedSrc === src) {
    return null
  }
  return <img className={className} src={src} alt='' onError={() => setFailedSrc(src)} />
}
