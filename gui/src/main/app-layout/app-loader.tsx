export function AppLoader() {
  return (
    <main className='grid h-screen w-full items-center justify-center text-white'>
      <i
        aria-label='loading'
        className='inline-block h-12 w-12 animate-spin rounded-full border-[4px] border-current border-t-transparent text-white'
        role='status'
      />
    </main>
  )
}
