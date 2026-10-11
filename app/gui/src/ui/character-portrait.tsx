const portraits = Object.fromEntries(
  Object.entries(
    import.meta.glob<string>('../assets/characters/*/*.webp', { eager: true, import: 'default' })
  ).map(([path, url]) => [path.slice(path.lastIndexOf('/') + 1, -'.webp'.length), url])
)

export function CharacterPortrait({
  character,
  className,
  style
}: {
  character: string
  className?: string
  style?: React.CSSProperties
}) {
  const src = portraits[character.toLowerCase().replace(/[^a-z0-9]/g, '')]
  if (!src) {
    return null
  }
  return (
    <div className={className} style={style}>
      <img className='absolute inset-0 h-full w-full object-cover' src={src} alt='' />
    </div>
  )
}
