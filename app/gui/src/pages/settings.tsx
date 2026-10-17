import React from 'react'
import { motion } from 'framer-motion'
import { useTranslation } from 'react-i18next'

import { ConfigContext } from '@/main/config'
import { useErrorPopup } from '@/main/error-popup'
import type { LocalizationKey } from '@/main/i18n'
import { LanguageSelect } from '@/main/language-select'

import {
  GetAppVersion,
  SaveBrowserSourcePort,
  SaveLocale,
  SaveLogFile,
  SaveSidebarMinimized,
  SaveTheme
} from '@cmd/CommandHandler'
import {
  BackupData,
  ChooseBackup,
  ChooseDataDirectory,
  GetDataLocation,
  HasSaveData,
  ImportBackup,
  OpenDataDirectory,
  SaveDataDirectory
} from '@cmd/DataHandler'
import { BrowserOpenURL } from '@runtime'

import * as Page from '@/ui/page'
import * as Dialog from '@/ui/dialog'
import { Button } from '@/ui/button'
import * as Select from '@/ui/select'
import { Switch } from '@/ui/switch'
import { Input } from '@/ui/input'
import { SettingButton, settingTextClass } from '@/ui/setting-button'
import { cn } from '@/helpers/cn'
import { model } from '@model'

export function SettingsPage() {
  const { t } = useTranslation()
  return (
    <Page.Root>
      <Page.Header>
        <Page.Title>{t('settings')}</Page.Title>
      </Page.Header>
      <motion.section
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className='overflow-y-scroll text-xl'
      >
        <div className='border-b-divider border-b-[1px] border-solid px-8 py-6'>
          <div className='mx-auto grid max-w-xl gap-3'>
            <h3 className='font-bold'>{t('generalSettings')}</h3>
            <div className='grid auto-rows-fr gap-2'>
              <ThemeSelect />
              <AppLanguageSelect />
              <BrowserSourcePortInput />
              <SideBarToggle />
              <LogFileToggle />
            </div>
          </div>
        </div>
        <div className='border-b-divider border-b-[1px] border-solid px-8 py-6'>
          <div className='mx-auto max-w-xl'>
            <SaveDataSettings />
          </div>
        </div>
        <div className='px-8 py-6'>
          <div className='mx-auto grid max-w-xl gap-3'>
            <div>
              <h3 className='font-bold'>{t('about')}</h3>
              <AppVersion />
            </div>
            <div className='flex flex-wrap gap-2'>
              <SettingButton
                icon='mdi:scroll'
                onClick={() => BrowserOpenURL('https://cfn.williamsjokvist.se/changelog')}
              >
                {t('changelog')}
              </SettingButton>
              <SettingButton
                icon='fa6-brands:github'
                onClick={() => BrowserOpenURL('https://github.com/williamsjokvist/cfn-tracker')}
              >
                {t('source')}
              </SettingButton>
              <SettingButton
                icon='fa6-brands:twitter'
                onClick={() => BrowserOpenURL('https://x.com/greensoap_')}
              >
                {t('follow')}
              </SettingButton>
            </div>
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
  return <p className={settingTextClass}>{t('appVersion', { appVersion: `v${appVersion}` })}</p>
}

function SideBarToggle() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [cfg, setCfg] = React.useContext(ConfigContext)
  return (
    <div className='flex w-full items-center justify-between'>
      <span className={settingTextClass}>{t('minimize')}</span>
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

function LogFileToggle() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [cfg, setCfg] = React.useContext(ConfigContext)
  return (
    <div className='flex w-full items-center justify-between'>
      <span className={settingTextClass}>{t('writeLogFile')}</span>
      <Switch
        checked={cfg.logFile}
        onCheckedChange={checked => {
          SaveLogFile(checked)
            .then(() => setCfg({ ...cfg, logFile: checked }))
            .catch(setError)
        }}
      />
    </div>
  )
}

function BrowserSourcePortInput() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [cfg, setCfg] = React.useContext(ConfigContext)
  const [port, setPort] = React.useState(String(cfg.browserSourcePort))
  const [saving, setSaving] = React.useState(false)

  // The config loads after the first render.
  React.useEffect(() => {
    setPort(String(cfg.browserSourcePort))
  }, [cfg.browserSourcePort])

  // Every save rebinds the server, so save on blur or Enter rather than on each keystroke.
  function save() {
    const newPort = Number(port)
    if (newPort === cfg.browserSourcePort) {
      return
    }
    // Wails rejects values that don't fit a uint16 before the command runs, with a generic error.
    if (port.trim() === '' || !Number.isInteger(newPort) || newPort < 1024 || newPort > 65535) {
      setPort(String(cfg.browserSourcePort))
      setError({
        localizationKey: model.ErrorLocalizationKey.errInvalidBrowserSourcePort,
        message: `port ${port}`
      })
      return
    }
    // Disabled until the save finishes, so two saves can't race to rebind the server.
    setSaving(true)
    SaveBrowserSourcePort(newPort)
      .then(() => setCfg({ ...cfg, browserSourcePort: newPort }))
      .catch(err => {
        setPort(String(cfg.browserSourcePort))
        setError(err)
      })
      .finally(() => setSaving(false))
  }

  return (
    <div className='flex w-full items-center justify-between'>
      <label htmlFor='browser-source-port' className={settingTextClass}>
        {t('browserSourcePort')}
      </label>
      <Input
        size='sm'
        id='browser-source-port'
        disabled={saving}
        type='number'
        inputMode='numeric'
        min={1024}
        max={65535}
        value={port}
        onValueChange={setPort}
        onBlur={save}
        onKeyDown={e => e.key === 'Enter' && e.currentTarget.blur()}
      />
    </div>
  )
}

function SaveDataSettings() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [location, setLocation] = React.useState<model.DataLocation | null>(null)
  const [busy, setBusy] = React.useState(false)
  const [folderWithSaveData, setFolderWithSaveData] = React.useState<string | null>(null)
  const [backupToImport, setBackupToImport] = React.useState<string | null>(null)

  function refresh() {
    GetDataLocation().then(setLocation).catch(setError)
  }
  React.useEffect(refresh, [])

  // Moving and importing copy the whole database, so disable the buttons until they finish.
  function run<T>(action: () => Promise<T>) {
    setBusy(true)
    return action()
      .catch(err => {
        setError(err)
        return undefined
      })
      .finally(() => setBusy(false))
  }

  function moveTo(dir: string) {
    if (!dir) {
      return
    }
    HasSaveData(dir)
      .then(hasSaveData =>
        hasSaveData ? setFolderWithSaveData(dir) : saveDataDirectory(dir, false)
      )
      .catch(setError)
  }

  function saveDataDirectory(dir: string, useExisting: boolean) {
    setFolderWithSaveData(null)
    run(() => SaveDataDirectory(dir, useExisting)).then(refresh)
  }

  function importBackup(path: string) {
    setBackupToImport(null)
    run(() => ImportBackup(path)).then(refresh)
  }

  return (
    <div className='grid gap-3'>
      <div>
        <h3 className='font-bold'>{t('saveData')}</h3>
        {/* Paths read left to right in every locale. */}
        <p dir='ltr' title={location?.dir} className={cn('truncate text-start', settingTextClass)}>
          {location?.dir}
        </p>
      </div>
      <div className='flex flex-wrap gap-2'>
        <SettingButton
          icon='mdi:folder-move-outline'
          disabled={busy}
          onClick={() => ChooseDataDirectory().then(moveTo).catch(setError)}
        >
          {t('changeFolder')}
        </SettingButton>
        <SettingButton
          icon='mdi:folder-open-outline'
          onClick={() => OpenDataDirectory().catch(setError)}
        >
          {t('openFolder')}
        </SettingButton>
        <SettingButton
          icon='mdi:restore'
          disabled={busy || !location || location.dir === location.defaultDir}
          onClick={() => location && moveTo(location.defaultDir)}
        >
          {t('resetToDefault')}
        </SettingButton>
      </div>
      <div className='flex flex-wrap gap-2'>
        <SettingButton
          icon='mdi:content-save-outline'
          disabled={busy}
          onClick={() => run(BackupData)}
        >
          {t('backUp')}
        </SettingButton>
        <SettingButton
          icon='mdi:database-import-outline'
          disabled={busy}
          onClick={() =>
            ChooseBackup()
              .then(path => path && setBackupToImport(path))
              .catch(setError)
          }
        >
          {t('importBackup')}
        </SettingButton>
      </div>
      {location?.restartRequired && (
        <p className='text-base text-yellow-300'>{t('restartToApply')}</p>
      )}

      <ConfirmDialog
        open={folderWithSaveData !== null}
        onClose={() => setFolderWithSaveData(null)}
        title='folderHasSaveData'
        text={t('folderHasSaveDataDescription')}
      >
        <Button onClick={() => saveDataDirectory(folderWithSaveData!, true)}>
          {t('useExisting')}
        </Button>
        <Button onClick={() => saveDataDirectory(folderWithSaveData!, false)}>
          {t('replace')}
        </Button>
      </ConfirmDialog>
      <ConfirmDialog
        open={backupToImport !== null}
        onClose={() => setBackupToImport(null)}
        title='importBackup'
        text={t('importBackupDescription')}
      >
        <Button onClick={() => importBackup(backupToImport!)}>{t('import')}</Button>
      </ConfirmDialog>
    </div>
  )
}

function ConfirmDialog(
  props: React.PropsWithChildren<{
    open: boolean
    onClose: () => void
    title: LocalizationKey
    text: string
  }>
) {
  return (
    <Dialog.Root open={props.open} onOpenChange={open => !open && props.onClose()}>
      {/* A style, as the class can't reliably override the dialog's fixed height. */}
      <Dialog.Content
        title={props.title}
        className='flex flex-col gap-6'
        style={{ height: 'auto' }}
      >
        <p>{props.text}</p>
        <div className='flex justify-end gap-3'>{props.children}</div>
      </Dialog.Content>
    </Dialog.Root>
  )
}

const Themes: { name: model.ThemeName; label: LocalizationKey; colors: string[] }[] = [
  {
    name: model.ThemeName.DEFAULT,
    label: 'themeDefault',
    colors: ['#901169', '#330083']
  },
  {
    name: model.ThemeName.ENTH,
    label: 'themeEnth',
    colors: ['#95F3F6', '#0E254D']
  },
  {
    name: model.ThemeName.TEKKEN,
    label: 'themeTekken',
    colors: ['#dd1d5b', '#1e3c52']
  },
  {
    name: model.ThemeName.JADE,
    label: 'themeJade',
    colors: ['#0c8f62', '#023a4d']
  },
  {
    name: model.ThemeName.EMBER,
    label: 'themeEmber',
    colors: ['#b8321f', '#3a0b2e']
  },
  {
    name: model.ThemeName.TIDE,
    label: 'themeTide',
    colors: ['#1767b8', '#0b1a4f']
  },
  {
    name: model.ThemeName.DUSK,
    label: 'themeDusk',
    colors: ['#5b3fd1', '#1a1240']
  },
  {
    name: model.ThemeName.SAKURA,
    label: 'themeSakura',
    colors: ['#d14f8a', '#3b1450']
  }
]

function ThemeSelect() {
  const { t } = useTranslation()
  const setError = useErrorPopup()
  const [cfg, setCfg] = React.useContext(ConfigContext)

  return (
    <div className='flex items-center justify-between'>
      <span className={settingTextClass}>{t('theme')}</span>
      <Select.Root
        size='sm'
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
          <Select.Item key={theme.name} value={theme.name} className='flex items-center gap-2'>
            {/* The ring keeps the darker half visible against the dark dropdown. */}
            <i
              style={{
                background: `linear-gradient(135deg, ${theme.colors[0]} 50%, ${theme.colors[1]} 50%)`
              }}
              className='h-5 w-5 shrink-0 rounded-full ring-1 ring-white/30'
            />
            <span>{t(theme.label)}</span>
          </Select.Item>
        ))}
      </Select.Root>
    </div>
  )
}

function AppLanguageSelect() {
  const { i18n, t } = useTranslation()
  const setError = useErrorPopup()
  return (
    <div className='flex w-full items-center justify-between'>
      <span className={settingTextClass}>{t('language')}</span>
      <LanguageSelect
        value={i18n.resolvedLanguage}
        onValueChange={code => {
          i18n
            .changeLanguage(code)
            .then(() => SaveLocale(code).catch(setError))
            .catch(setError)
        }}
      />
    </div>
  )
}
