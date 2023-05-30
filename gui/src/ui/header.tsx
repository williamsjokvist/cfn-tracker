type PageHeaderProps = {
  text: string;
  showSpinner?: boolean;
}
export const PageHeader: React.FC<React.PropsWithChildren<PageHeaderProps>> = ( { text, showSpinner, children } ) => {
  return (
    <header className="border-b border-slate-50 backdrop-blur border-opacity-10 select-none" style={{ "--wails-draggable": "drag" } as React.CSSProperties}>
      <h2 className="pt-4 px-8 flex items-center justify-between gap-5 uppercase text-sm tracking-widest mb-4">
        {text}
        {showSpinner && (
          <i
            className="animate-spin inline-block w-5 h-5 border-[3px] border-current border-t-transparent text-pink-600 rounded-full"
            role="status"
            aria-label="loading"
          />
        )}
      </h2>
      {children}
    </header>
  )
}