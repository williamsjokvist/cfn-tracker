import React from 'react'
import { createRoot } from 'react-dom/client'
import { Icon } from '@iconify/react'

import { useErrorPopup } from '@/main/error-popup'

import * as Dialog from '@/ui/dialog'
import { Checkbox } from '@/ui/checkbox'
import { Button } from '@/ui/button'

import { GetThemes } from '@cmd'
import type { model } from '@model'

export function ThemeSelect(props: { selectedTheme: string; onSelect: (theme: string) => void }) {
  const containerRef = React.useRef<HTMLDivElement>(null)
  const [themes, setThemes] = React.useState<model.Theme[]>([])
  const [isOpen, setOpen] = React.useState(false)
  const setError = useErrorPopup()

  React.useEffect(() => {
    GetThemes().then(setThemes).catch(setError)
  }, [])

  React.useEffect(() => {
    if (!isOpen || !containerRef.current) return

    for (const theme of themes) {
      const container = document.createElement('div')
      containerRef.current.appendChild(container)
      const shadowRoot = container.attachShadow({ mode: 'open' })

      createRoot(container).render(
        <div className='mb-4'>
          <style>{theme.css.match(/@import.*;/g)}</style>
          <Checkbox
            id={`${theme.name}-checkbox`}
            checked={theme.name === props.selectedTheme}
            onChange={e => props.onSelect(theme.name)}
          />
          <label
            htmlFor={`${theme.name}-checkbox`}
            className='font-spartan cursor-pointer text-lg font-bold capitalize'
          >
            {theme.name}
          </label>
        </div>
      )

      createRoot(shadowRoot).render(
        <>
          <slot />
          <div className='stat-list'>
            <style>{theme.css}</style>
            <div className='stat-item'>
              <span className='stat-title'>MR</span>
              <span className='stat-value'>444</span>
            </div>
          </div>
        </>
      )
    }

    return () => {
      if (containerRef.current) {
        containerRef.current.innerHTML = ''
      }
    }
  }, [isOpen, themes, props.selectedTheme])

  return (
    <Dialog.Root onOpenChange={setOpen}>
      <Dialog.Trigger asChild>
        <Button
          className='capitalize'
          style={{ filter: 'hue-rotate(-180deg)', justifyContent: 'center' }}
        >
          <Icon icon='ph:paint-bucket-fill' className='mr-3 h-6 w-6' />
          {props.selectedTheme}
        </Button>
      </Dialog.Trigger>
      <Dialog.Content title='selectTheme'>
        <div ref={containerRef} className='mt-2 grid h-80 w-full gap-4 overflow-y-scroll pr-2' />
      </Dialog.Content>
    </Dialog.Root>
  )
}
