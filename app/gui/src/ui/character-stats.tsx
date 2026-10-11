import React from 'react'

import { CharacterPortrait } from '@/ui/character-portrait'

export function CharacterStatsCard({
  character,
  heading,
  children
}: React.PropsWithChildren<{ character: string; heading: React.ReactNode }>) {
  // The portrait is a square as tall as the stats, which CSS can't size from a sibling.
  const statsRef = React.useRef<HTMLDivElement>(null)
  const [size, setSize] = React.useState(0)
  React.useLayoutEffect(() => {
    const observer = new ResizeObserver(([entry]) => setSize(entry.borderBoxSize[0].blockSize))
    observer.observe(statsRef.current!)
    return () => observer.disconnect()
  }, [])

  return (
    <section className='flex items-start gap-4'>
      <CharacterPortrait
        character={character}
        className='relative shrink-0 overflow-hidden rounded-xl bg-slate-50/5'
        style={{ width: size, height: size }}
      />
      <div ref={statsRef} className='flex-1'>
        {heading}
        <dl className='grid grid-cols-3 gap-2 [&>div]:mb-0'>{children}</dl>
      </div>
    </section>
  )
}
