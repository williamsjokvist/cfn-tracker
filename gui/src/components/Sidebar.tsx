import React from "react";
import { useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useAnimate } from "framer-motion";
import { Icon } from '@iconify/react';

import { LanguageSelector } from "./LanguageSelector";
import { BrowserOpenURL, Quit, WindowMinimise } from "@@/runtime";
import { GetAppVersion } from "@@/go/core/CommandHandler";

type SidebarLinkProps = {
  icon: any;
  link: string;
  name: string;
  isSelected?: boolean;
  selectedIcon?: any;
  isMinimized: boolean;
}
const SidebarLink: React.FC<SidebarLinkProps> = ( { icon, link, name, isSelected, isMinimized, selectedIcon } ) => {
  return (
    <a
      href={"#/" + link}
      className="text-lg text-[#bfbcff] text-opacity-80 rounded py-2 px-1 group flex items-center justify-between hover:!text-white hover:bg-slate-50 hover:bg-opacity-5 transition-colors"
      style={{
        fontWeight: isSelected ? "600" : "200",
        color: isSelected ? "#d6d4ff" : "#bfbcff",
      }}
    >
      <span className="flex items-center justify-between">
        <Icon icon={isSelected ? selectedIcon : icon} className="text-[#f85961] transition-colors w-10 h-7 mr-1" />
        <HideableText text={name} hide={isMinimized}/>
      </span>
      <Icon 
        icon='fa6-solid:chevron-left' 
        className="w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity"
        style={{ transform: "rotate(180deg)", display: isMinimized ? 'none' : 'block' }}
      />
    </a>
  );
};

export const Sidebar: React.FC<React.PropsWithChildren> = ( { children } ) => {
  const { t } = useTranslation();
  const location = useLocation();
  const [appVersion, setAppVersion] = React.useState('');
  const [isMinimized, setMinimized] = React.useState(window.localStorage.getItem('sidebar-minimized') == 'true')

  const [scope, animate] = useAnimate()

  React.useEffect(() => {
    animate('a, button', 
      { opacity: [0, 1] }, 
      { delay: 0.125 }
    )
  }, [])
  
  React.useEffect(() => {
    !appVersion && GetAppVersion().then(v => setAppVersion(v));
  }, [appVersion])

  React.useEffect(() => {
    if (isMinimized) {
      window.localStorage.setItem('sidebar-minimized', 'true')
    } else if (!isMinimized) {
      window.localStorage.removeItem('sidebar-minimized')
    }
  }, [isMinimized])

  return (
    <aside
      ref={scope}
      className="sidebar"
      style={{
        width: isMinimized ? '76px' : "175px",
      }}
    >
      <header style={{"--draggable": "drag"} as React.CSSProperties}>
        <div className="flex justify-start">
          <div className="group mt-2 group ml-2 flex mb-3">
            <button
              aria-label="close"
              className="mr-[8px] p-[2px] w-[14px] h-[14px] group-hover:bg-red-500 bg-slate-600 flex items-center justify-center rounded-full transition-all"
              onClick={Quit}
            >
              <Icon icon='ep:close-bold' className="text-red-800 group-hover:opacity-100 opacity-0 transition-all" />
            </button>
            <button
              aria-label="close"
              className="p-[2px] w-[14px] h-[14px] group-hover:bg-yellow-500 bg-slate-600 flex items-center justify-center rounded-full transition-all"
              onClick={WindowMinimise}
            >
              <Icon icon='mingcute:minimize-fill' className="text-yellow-800 group-hover:opacity-100 opacity-0 transition-all" />
            </button>
          </div>
        </div>
      </header>
      <nav className="mt-5 w-full">
        <ul>
          <li>
            <SidebarLink
              icon='ri:search-line'
              selectedIcon='ri:search-fill'
              link=""
              name={t("tracking")}
              isSelected={location.pathname == "/"}
              isMinimized={isMinimized}
            />
          </li>
          <li>
            <SidebarLink
              icon='ion:document-text-outline'
              selectedIcon='ion:document-text'
              link="history"
              name={t("history")}
              isSelected={location.pathname == "/history"}
              isMinimized={isMinimized}
            />
          </li>
          <li>
            <SidebarLink
              icon='clarity:sign-out-line'
              selectedIcon='clarity:sign-out-solid'
              link="output"
              name={'Output'}
              isSelected={location.pathname == "/output"}
              isMinimized={isMinimized}
            />
          </li>
        </ul>
      </nav>
      {children}
      <footer className={`grid w-full text-xl px-2`}>
        <LanguageSelector isMinimized={isMinimized} />

        {/* Twitter */}
        <a
          target="#"
          onClick={() => BrowserOpenURL("https://twitter.com/greensoap_")}
          className={`h-[28px] cursor-pointer w-full group font-extralight flex justify-between items-center mt-1 text-[#d6d4ff] hover:text-white transition-colors`}
        >
          <span className={`flex items-center justify-between lowercase`}>
            <Icon icon='fa6-brands:twitter' className="text-[#49b3f5] w-4 h-4 mr-2 transition-colors group-hover:text-white" />
            <HideableText text={t('follow')} hide={isMinimized}/>
          </span>
          <Icon icon='fa6-solid:arrow-up' className="relative right-[-8px] w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity" style={{ transform: "rotate(45deg)" }} />
        </a>

        {/* Minimize */}
        <button 
          type='button'
          className={`h-[28px] cursor-pointer w-full group font-extralight flex items-center mt-1 text-[#d6d4ff] hover:text-white transition-colors`}
          onClick={() => setMinimized(!isMinimized)}>
          <Icon 
            icon='fa6-solid:chevron-left' 
            className="group-hover:text-white text-[#d6d4ff] w-4 h-4 transition-all"
            style={{ transform: isMinimized ? "rotate(-180deg)" : 'none' }}
          />
          <HideableText className="ml-2" text={t('minimize')} hide={isMinimized}/>
        </button>

        {/* Version */}
        <a
          target="#"
          className="text-sm mt-4 font-extralight cursor-pointer hover:underline"
          onClick={() => BrowserOpenURL("https://github.com/GreenSoap/cfn-tracker/releases")}
        >
          {isMinimized ? `v${appVersion}` : `CFN Tracker v${appVersion}`}
        </a>
      </footer>
    </aside>
  );
};

type HideableTextProps = {
  text: string;
  hide: boolean;
  className?: string;
}
export const HideableText: React.FC<HideableTextProps> = ( { text, hide, className }) => 
  <span className={`transition-opacity ${hide ? 'opacity-0 w-0 overflow-hidden' : 'opacity-100'} ${className}`}>{text}</span>