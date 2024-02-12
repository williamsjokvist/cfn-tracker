import React from 'react'
import { Icon } from '@iconify/react'
import { useAnimate } from 'framer-motion'
import { useTranslation } from 'react-i18next'

import { cn } from '@/helpers/cn'
import { BrowserOpenURL } from '@@/runtime'
import {
  GetAppVersion,
  GetSupportedLanguages,
  SaveLocale,
  SaveSidebarMinimized
} from '@@/go/core/CommandHandler'

import { ConfigContext } from '../config'
import { useErrorPopup } from '../error-popup'
import { AppTitleBar } from './app-titlebar'
import { Nav } from './app-nav'

export function AppSidebar() {
  const [cfg] = React.useContext(ConfigContext)

  const [scope, animate] = useAnimate()
  React.useEffect(() => {
    animate('a, button', { opacity: [0, 1] }, { delay: 0.25 })
  }, [])

  React.useEffect(() => {
    animate(
      'span',
      { opacity: +!cfg.sidebarMinified, display: cfg.sidebarMinified ? 'none' : 'block' },
      { duration: 0.175, ease: 'circIn' }
    )
  }, [cfg.sidebarMinified])

  return (
    <div
      ref={scope}
      className={cn(
        'relative z-50 grid grid-rows-[0fr_1fr_0fr] gap-5',
        'select-none overflow-visible whitespace-nowrap px-[10px] py-3',
        'bg-[rgba(3,5,19,0.33)] text-white',
        'transition-[width_250ms_ease-out]'
      )}
      style={{
        width: cfg.sidebarMinified ? '76px' : '175px'
      }}
    >
      <AppTitleBar />
      <Nav />
      <div className='flex flex-col items-start px-2 text-xl'>
        <LanguageSelector />
        <TwitterLink />
        <MinimizeButton />
        <AppVersion />
      </div>
    </div>
  )
}

function AppVersion() {
  const [appVersion, setAppVersion] = React.useState('')
  React.useEffect(() => {
    !appVersion && GetAppVersion().then(v => setAppVersion(v))
  }, [appVersion])
  return (
    <button
      className='mt-4 cursor-pointer text-sm font-extralight hover:underline'
      onClick={() => BrowserOpenURL('https://github.com/GreenSoap/cfn-tracker/releases')}
    >
      {`v${appVersion}`}
    </button>
  )
}

function LanguageSelector(props: React.PropsWithChildren) {
  const { i18n, t } = useTranslation()
  const [langs, setLangs] = React.useState<string[]>([])
  const setError = useErrorPopup()

  React.useEffect(() => {
    GetSupportedLanguages().then(setLangs).catch(setError)
  }, [])

  const changeLanguage = (code: string) => {
    i18n.changeLanguage(code)
    SaveLocale(code).catch(setError)
  }

  return (
    <div className='group flex w-full'>
      <button
        type='button'
        className='relative left-0 flex h-[28px] w-full cursor-default items-center font-thin lowercase text-[#d6d4ff] transition-colors group-hover:text-white'
      >
        <Icon
          icon='fa6-solid:globe'
          className='mr-2 h-4 w-4 text-[#d6d4ff] transition-all group-hover:text-white'
        />
        <span>{t('language')}</span>
      </button>
      <div className='invisible absolute left-[98%] flex opacity-0 transition-all group-hover:visible group-hover:opacity-100'>
        <Icon
          icon='fa6-solid:chevron-left'
          className='relative right-4 top-2 h-3 w-3 rotate-180 text-white'
        />
        <ul className='group relative bottom-1 flex w-[195px] justify-between gap-2 rounded-lg bg-[rgba(0,0,0,.625)] px-3 py-2 text-[16px] text-base uppercase leading-5 text-[#bfbcff] backdrop-blur transition-colors hover:bg-[rgba(0,0,0,.525)]'>
          {langs.map(code => (
            <li key={code}>
              <button
                onClick={() => changeLanguage(code)}
                type='button'
                className='cursor-pointer transition-colors first-letter:uppercase hover:!text-white'
                {...(i18n.resolvedLanguage === code && {
                  style: {
                    fontWeight: 600
                  }
                })}
              >
                {
                  new Intl.DisplayNames([code], {
                    type: 'language',
                    languageDisplay: 'standard',
                    style: 'narrow',
                    fallback: 'code'
                  })
                    .of(code)!
                    .split(/[ |（]/)[0]
                }
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}

function TwitterLink() {
  const { t } = useTranslation()
  return (
    <button
      onClick={() => BrowserOpenURL('https://twitter.com/greensoap_')}
      className='group mt-1 flex h-[28px] w-full items-center justify-between font-extralight text-[#d6d4ff] transition-colors hover:text-white'
    >
      <div className='flex items-center justify-between lowercase'>
        <Icon
          icon='fa6-brands:twitter'
          className='mr-2 h-4 w-4 text-[#49b3f5] transition-colors group-hover:text-white'
        />
        <span>{t('follow')}</span>
      </div>
      <Icon
        icon='fa6-solid:arrow-up'
        className='relative right-[-8px] h-3 w-3 opacity-0 transition-opacity group-hover:opacity-100'
        style={{ transform: 'rotate(45deg)' }}
      />
    </button>
  )
}

function MinimizeButton() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [cfg, setCfg] = React.useContext(ConfigContext)
  return (
    <button
      type='button'
      className={cn(
        'group flex items-center',
        'text-[#d6d4ff] transition-colors hover:text-white',
        'h-[28px] w-full cursor-pointer',
        'mt-1 font-extralight'
      )}
      onClick={() => {
        SaveSidebarMinimized(!cfg.sidebarMinified)
          .then(() => setCfg({ ...cfg, sidebarMinified: !cfg.sidebarMinified }))
          .catch(setError)
      }}
    >
      <Icon
        icon='fa6-solid:chevron-left'
        className='h-4 w-4 text-[#d6d4ff] transition-all group-hover:text-white'
        style={{ transform: cfg.sidebarMinified ? 'rotate(-180deg)' : 'none' }}
      />
      <span className='ml-2'>{t('minimize')}</span>
    </button>
  )
}
