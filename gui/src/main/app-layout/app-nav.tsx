import React from 'react'
import { Icon } from '@iconify/react'
import { useTranslation } from 'react-i18next'
import { useLocation, Link } from 'react-router-dom'

import { cn } from '@/helpers/cn'
import { ConfigContext } from '@/main/config'
import { LocalizationKey } from '@/main/i18n'

const NavItems = [
  {
    icons: ['ri:search-line', 'ri:search-fill'],
    href: 'tracking'
  },
  {
    icons: ['ion:document-text-outline', 'ion:document-text'],
    href: 'sessions'
  },
  {
    icons: ['clarity:sign-out-line', 'clarity:sign-out-solid'],
    href: 'output'
  }
] as const

export function Nav() {
  const location = useLocation()
  return (
    <nav>
      {NavItems.map(({ href, icons }) => (
        <NavItem key={href} href={href} icons={icons} selected={location.pathname.includes(href)} />
      ))}
    </nav>
  )
}

function NavItem(
  props: {
    icons: readonly [string, string]
    href: LocalizationKey
    selected?: boolean
  } & React.PropsWithChildren
) {
  const { t } = useTranslation()
  const [cfg] = React.useContext(ConfigContext)

  return (
    <Link
      to={props.href}
      className={cn(
        'flex items-center justify-between',
        'group flex items-center justify-between',
        'transition-colors hover:bg-slate-50 hover:bg-opacity-5 hover:!text-white active:bg-[rgba(255,255,255,.075)]',
        'rounded px-1 py-2 text-lg text-[#bfbcff] text-opacity-80'
      )}
      style={{
        fontWeight: props.selected ? '600' : '200',
        color: props.selected ? '#d6d4ff' : '#bfbcff'
      }}
    >
      <div className='flex items-center justify-between'>
        <Icon
          icon={props.selected ? props.icons[1] : props.icons[0]}
          className='mr-1 h-7 w-10 text-[#f85961] transition-colors'
        />
        <span>{t(props.href)}</span>
      </div>
      {!cfg.sidebarMinified && (
        <Icon
          icon='fa6-solid:chevron-left'
          className='h-3 w-3 rotate-180 opacity-0 transition-opacity group-hover:opacity-100'
        />
      )}
    </Link>
  )
}
