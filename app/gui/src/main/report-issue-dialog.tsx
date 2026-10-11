import React from 'react'
import { Icon } from '@iconify/react'
import { useTranslation } from 'react-i18next'

import { ReportIssue } from '@cmd/CommandHandler'
import { BrowserOpenURL } from '@runtime'
import { model } from '@model'

import * as Dialog from '@/ui/dialog'
import { Button } from '@/ui/button'
import { Input } from '@/ui/input'
import { cn } from '@/helpers/cn'

type Status =
  | { kind: 'editing' }
  | { kind: 'sending' }
  | { kind: 'sent'; issueUrl: string }
  | { kind: 'failed' }

// The sidebar is never unmounted, so the outcome of a report stays until the app is closed.
export function ReportIssueDialog() {
  const { t } = useTranslation()
  const [open, setOpen] = React.useState(false)
  const [status, setStatus] = React.useState<Status>({ kind: 'editing' })
  const [report, setReport] = React.useState<model.IssueReport>({
    title: '',
    description: '',
    contact: ''
  })

  const sending = status.kind === 'sending'
  const canSend = report.title.trim() !== '' && report.description.trim() !== ''

  function submit(e: React.FormEvent) {
    e.preventDefault()
    // Button doesn't forward `disabled` to the element, so it still submits the form.
    if (!canSend) {
      return
    }
    setStatus({ kind: 'sending' })
    ReportIssue(report)
      .then(issueUrl => setStatus({ kind: 'sent', issueUrl }))
      .catch(() => setStatus({ kind: 'failed' }))
  }

  return (
    <Dialog.Root open={open} onOpenChange={open => !sending && setOpen(open)}>
      <Dialog.Trigger asChild>
        <button
          aria-label={t('reportIssue')}
          title={t('reportIssue')}
          className={cn(
            'grid h-9 w-9 place-items-center rounded-sm text-[#bfbcff]/60',
            'transition-colors hover:bg-slate-50/5 hover:text-white active:bg-white/[0.075]'
          )}
        >
          <Icon icon='mdi:alarm-light-outline' className='h-5 w-5' />
        </button>
      </Dialog.Trigger>
      <Dialog.Content
        title='reportIssue'
        description='reportIssueDescription'
        dismissible={!sending}
        className='flex flex-col'
      >
        {status.kind === 'editing' ? (
          <form onSubmit={submit} className='mt-4 flex flex-1 flex-col gap-3'>
            <Input
              aria-label={t('issueTitle')}
              placeholder={t('issueTitle')}
              maxLength={120}
              className='w-full'
              value={report.title}
              onValueChange={title => setReport({ ...report, title })}
            />
            <textarea
              aria-label={t('issueDescription')}
              placeholder={t('issueDescription')}
              maxLength={5000}
              className={cn(
                'min-h-0 w-full flex-1 resize-none rounded-lg border-[1px] border-white/20 px-4 py-2 text-lg text-white',
                'bg-white/5 transition-colors outline-none',
                'hover:bg-white/10 focus:border-white/50 focus:bg-white/10'
              )}
              value={report.description}
              onChange={e => setReport({ ...report, description: e.target.value })}
            />
            <Input
              aria-label={t('issueContact')}
              placeholder={t('issueContact')}
              maxLength={200}
              className='w-full'
              value={report.contact}
              onValueChange={contact => setReport({ ...report, contact })}
            />
            <Button type='submit' disabled={!canSend} className='self-end'>
              {t('send')}
            </Button>
          </form>
        ) : (
          <div className='flex flex-1 flex-col items-center justify-center gap-4 text-center'>
            {status.kind === 'sending' && (
              <i
                aria-label='loading'
                role='status'
                className='text-highlight inline-block h-12 w-12 animate-spin rounded-full border-[4px] border-current border-t-transparent'
              />
            )}
            {status.kind === 'sent' && (
              <>
                <Icon icon='mdi:check-circle-outline' className='h-12 w-12 text-green-400' />
                <p>{t('issueReported')}</p>
                <Button onClick={() => BrowserOpenURL(status.issueUrl)}>
                  {t('viewIssue')}
                  <Icon icon='fa6-solid:arrow-up' className='ms-2 h-3 w-3 rotate-45' />
                </Button>
              </>
            )}
            {status.kind === 'failed' && (
              <>
                <Icon icon='mdi:alert-circle-outline' className='h-12 w-12 text-[#ff6388]' />
                <p>{t('errReportIssue')}</p>
              </>
            )}
          </div>
        )}
      </Dialog.Content>
    </Dialog.Root>
  )
}
