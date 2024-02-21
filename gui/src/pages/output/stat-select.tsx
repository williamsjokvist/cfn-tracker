import { useTranslation } from 'react-i18next'
import { Icon } from '@iconify/react'

import * as Dialog from '@/ui/dialog'
import { Checkbox } from '@/ui/checkbox'
import { Button } from '@/ui/button'

import type { StatOptions } from '.'

export function StatSelect(props: {
  options: Omit<StatOptions, 'theme'>
  onSelect: (option: string, checked: boolean) => void
}) {
  const { t } = useTranslation()
  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <Button style={{ filter: 'hue-rotate(-45deg)' }}>
          <Icon icon='bx:stats' className='mr-3 h-6 w-6' />
          {t('displayStats')}
        </Button>
      </Dialog.Trigger>
      <Dialog.Content title='displayStats' description='statsWillBeDisplayed'>
        <ul className='mt-4 h-72 overflow-y-scroll'>
          {Object.entries(props.options).map(([opt, checked]) => {
            return (
              <li key={opt}>
                <button
                  className='flex w-full cursor-pointer items-center px-2 py-1 text-lg hover:bg-[rgba(255,255,255,0.075)]'
                  onClick={() => props.onSelect(opt, !checked)}
                >
                  <Checkbox checked={props.options[opt] === true} readOnly />
                  <span className='ml-2 cursor-pointer text-center capitalize'>{opt}</span>
                </button>
              </li>
            )
          })}
        </ul>
      </Dialog.Content>
    </Dialog.Root>
  )
}
