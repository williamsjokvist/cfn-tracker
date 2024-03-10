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
  const containerRef = React.useRef<HTMLUListElement>(null)
  const [themes, setThemes] = React.useState<model.Theme[]>([])
  const [isOpen, setOpen] = React.useState(false)
  const setError = useErrorPopup()

  React.useEffect(() => {
    GetThemes().then(setThemes).catch(setError)
  }, [])

  React.useEffect(() => {
    if (!isOpen || !containerRef.current) {
      return
    }

    for (const theme of themes) {
      const container = document.querySelector(`.${theme.name}-preview`)
      if (!container) {
        continue
      }
      const shadowRoot = container.attachShadow({ mode: 'open' })
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
  }, [isOpen, themes])

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
        <ul ref={containerRef} className='mt-2 grid h-80 w-full gap-4 overflow-y-scroll pr-2'>
          {themes.map(theme => (
            <li key={theme.name}>
              <div className='mb-4 flex px-2 py-1 hover:bg-white hover:bg-opacity-[.075]'>
                <Checkbox
                  id={`${theme.name}-checkbox`}
                  checked={theme.name === props.selectedTheme}
                  onChange={e => props.onSelect(theme.name)}
                />
                <label
                  htmlFor={`${theme.name}-checkbox`}
                  className='font-spartan w-full cursor-pointer text-lg font-bold capitalize'
                >
                  {theme.name}
                </label>
              </div>
              <div className={`${theme.name}-preview`}>
                <style>{theme.css.match(/@import.*;/g)}</style>
              </div>
            </li>
          ))}
        </ul>
      </Dialog.Content>
    </Dialog.Root>
  )
}
