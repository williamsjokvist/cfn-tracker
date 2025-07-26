import React from 'react'
import { Icon } from '@iconify/react'
import { motion } from 'framer-motion'
import { useTranslation } from 'react-i18next'

import { ConfigContext } from '@/main/config'
import { useErrorPopup } from '@/main/error-popup'

import {
  GetAppVersion,
  GetSupportedLanguages,
  SaveLocale,
  SaveSidebarMinimized,
  SaveTheme
} from '@cmd/CommandHandler'
import { BrowserOpenURL } from '@runtime'

import * as Page from '@/ui/page'
import * as Select from '@/ui/select'
import { Switch } from '@/ui/switch'
import { Flag } from '@/ui/flag'
import { cn } from '@/helpers/cn'
import { model } from '@model'

export function SettingsPage() {
  const { t } = useTranslation()
  return (
    <Page.Root>
      <Page.Header>
        <div className='mx-auto w-full max-w-xl text-left'>
          <Page.Title>{t('settings')}</Page.Title>
        </div>
      </Page.Header>
      <motion.section
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='text-xl'
      >
        <div className='border-b-divider border-b-[1px] border-solid px-8 py-6'>
          <div className='mx-auto grid max-w-xl gap-4'>
            <ThemeSelect />
            <LanguageSelect />
            <SideBarToggle />
          </div>
        </div>
        <div className='relative mx-auto grid max-w-xl gap-4 py-6'>
          <h3 className='font-bold'>{t('about')}</h3>
          <AppVersion />
          <div className='flex gap-8'>
            <Link
              icon={
                <Icon
                  icon='mdi:scroll'
                  className='mr-2 h-6 w-6 text-stone-300 transition-colors group-hover:text-white'
                />
              }
              text={t('changelog')}
              url='https://cfn.williamsjokvist.se/changelog'
            />
            <Link
              icon={
                <Icon
                  icon='fa6-brands:github'
                  className='mr-2 h-6 w-6 text-orange-300 transition-colors group-hover:text-white'
                />
              }
              text={t('source')}
              url='https://github.com/williamsjokvist/cfn-tracker'
            />
            <Link
              icon={
                <Icon
                  icon='fa6-brands:twitter'
                  className='mr-2 h-6 w-6 text-[#49b3f5] transition-colors group-hover:text-white'
                />
              }
              text={t('follow')}
              url='https://x.com/greensoap_'
            />
          </div>
        </div>
      </motion.section>
    </Page.Root>
  )
}

function AppVersion() {
  const { t } = useTranslation()
  const [appVersion, setAppVersion] = React.useState('')
  const setError = useErrorPopup()
  React.useEffect(() => {
    !appVersion &&
      GetAppVersion()
        .then(v => setAppVersion(v))
        .catch(setError)
  }, [appVersion])
  return <span>{t('appVersion', { appVersion: `v${appVersion}` })}</span>
}

function SideBarToggle() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [cfg, setCfg] = React.useContext(ConfigContext)
  return (
    <div className='flex w-full justify-between'>
      <h3 className='font-bold'>{t('minimize')}</h3>
      <Switch
        checked={cfg.sidebar}
        onCheckedChange={checked => {
          SaveSidebarMinimized(checked)
            .then(() => setCfg({ ...cfg, sidebar: checked }))
            .catch(setError)
        }}
      />
    </div>
  )
}

const Themes = [
  {
    name: model.ThemeName.DEFAULT,
    colors: ['#901169', '#330083']
  },
  {
    name: model.ThemeName.ENTH,
    colors: ['#95F3F6', '#0E254D']
  },
  {
    name: model.ThemeName.TEKKEN,
    colors: ['#dd1d5b', '#1e3c52']
  }
]

function ThemeSelect() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [cfg, setCfg] = React.useContext(ConfigContext)

  return (
    <div className='flex items-center justify-between'>
      <h3 className='font-bold'>{t('theme')}</h3>
      <Select.Root
        value={cfg.theme}
        onValueChange={(theme: model.ThemeName) => {
          SaveTheme(theme)
            .then(() => {
              document.body.setAttribute('data-theme', theme)
              setCfg({ ...cfg, theme: theme })
            })
            .catch(setError)
        }}
      >
        {Themes.map(theme => (
          <Select.Item
            key={theme.name}
            value={theme.name}
            className='flex items-center justify-between gap-2'
          >
            {theme.colors.map(color => (
              <i key={color} style={{ background: color }} className={`h-4 w-3 rounded-md`} />
            ))}
            <span className='first-letter:uppercase'>{theme.name}</span>
          </Select.Item>
        ))}
      </Select.Root>
    </div>
  )
}

function LanguageSelect(props: React.PropsWithChildren) {
  const { i18n, t } = useTranslation()
  const [langs, setLangs] = React.useState<string[]>([])
  const setError = useErrorPopup()

  React.useEffect(() => {
    GetSupportedLanguages().then(setLangs).catch(setError)
  }, [])

  return (
    <div className='flex w-full items-center justify-between'>
      <h3 className='font-bold'>{t('language')}</h3>
      <Select.Root
        onValueChange={code => {
          i18n
            .changeLanguage(code)
            .then(() => SaveLocale(code).catch(setError))
            .catch(setError)
        }}
        value={i18n.resolvedLanguage}
      >
        {langs.map(code => {
          const nativeName = new Intl.DisplayNames([code], {
            type: 'language',
            languageDisplay: 'standard',
            style: 'narrow',
            fallback: 'code'
          })
            .of(code)!
            .split(/[ |（]/)[0]

          return (
            <Select.Item
              key={code}
              value={code}
              className='flex items-center justify-between gap-2'
              {...(i18n.resolvedLanguage === code && {
                style: {
                  fontWeight: 600
                }
              })}
            >
              <Flag code={code} className='w-8 rounded-md' />
              <span className='first-letter:uppercase'>{nativeName}</span>
            </Select.Item>
          )
        })}
      </Select.Root>
    </div>
  )
}

function Link(props: {
  url: string
  icon: React.ReactNode
  text: React.ReactNode
  className?: string
}) {
  return (
    <button
      onClick={() => BrowserOpenURL(props.url)}
      className={cn(
        'group mt-1 flex h-[28px] items-center justify-between',
        'font-extralight text-[#d6d4ff] transition-colors',
        'hover:text-white'
      )}
    >
      <div className='flex items-center justify-between whitespace-nowrap lowercase'>
        {props.icon}
        <span>{props.text}</span>
      </div>
      <Icon
        icon='fa6-solid:arrow-up'
        className='relative right-[-8px] h-3 w-3 opacity-0 transition-opacity group-hover:opacity-100'
        style={{ transform: 'rotate(45deg)' }}
      />
    </button>
  )
}
