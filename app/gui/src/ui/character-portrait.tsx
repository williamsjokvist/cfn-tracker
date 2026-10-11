const portraits = Object.fromEntries(
  Object.entries(
    import.meta.glob<string>('../assets/characters/*.webp', { eager: true, import: 'default' })
  ).map(([path, url]) => [path.slice(path.lastIndexOf('/') + 1, -'.webp'.length), url])
)

export function CharacterPortrait({
  character,
  className
}: {
  character: string
  className?: string
}) {
  const src = portraits[character.toLowerCase().replace(/[^a-z0-9]/g, '')]
  if (!src) {
    return null
  }
  return <img className={className} src={src} alt='' />
}
