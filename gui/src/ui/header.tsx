import { motion } from "framer-motion";

type PageHeaderProps = {
  text: string;
  showSpinner?: boolean;
}
export const PageHeader: React.FC<React.PropsWithChildren<PageHeaderProps>> = ( { text, showSpinner, children } ) => {
  return (
    <header className="page-header" style={{ "--draggable": "drag" } as React.CSSProperties}>
      <motion.h2 
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className="whitespace-nowrap uppercase text-sm tracking-widest"
      >
        {text}
      </motion.h2>
      {showSpinner && (
        <motion.i
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className="animate-spin inline-block w-5 h-5 border-[3px] border-current border-t-transparent text-pink-600 rounded-full"
          role="status"
          aria-label="loading"
        />
      )}
      {children}
    </header>
  )
}