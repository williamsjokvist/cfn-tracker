import React from "react";
import { useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useAnimate } from "framer-motion";

import { FaTwitter, FaChevronLeft, FaArrowUp } from "react-icons/fa";
import { RxUpdate } from "react-icons/rx";
import { IoDocumentTextOutline, IoDocumentText } from "react-icons/io5";
import { RiSearch2Line, RiSearch2Fill } from "react-icons/ri";
import { VscChromeMinimize, VscChromeClose } from "react-icons/vsc";

import { LanguageSelector } from "./LanguageSelector";

import { useAppStore } from "@/store/use-app-store";
import { BrowserOpenURL, Quit, WindowMinimise } from "@@/runtime";
import { GetAppVersion } from "@@/go/core/CommandHandler";

type SidebarLinkProps = {
  Icon: any;
  link: string;
  name: string;
  isSelected?: boolean;
  SelectedIcon?: any;
  isMinimized: boolean;
}
const SidebarLink: React.FC<SidebarLinkProps> = ( { Icon, link, name, isSelected, isMinimized, SelectedIcon } ) => {
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
        {isSelected ? 
          <SelectedIcon className="text-[#f85961] transition-colors w-10 h-7 mr-1" /> : 
          <Icon className="text-[#f85961] transition-colors w-10 h-7 mr-1" />
        }
        <HideableText text={name} hide={isMinimized}/>
      </span>
      <FaChevronLeft
        className="w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity"
        style={{ transform: "rotate(180deg)", display: isMinimized ? 'none' : 'block' }}
      />
    </a>
  );
};

export const Sidebar: React.FC = () => {
  const { t } = useTranslation();
  const location = useLocation();
  const [appVersion, setAppVersion] = React.useState('');
  const [isMinimized, setMinimized] = React.useState(window.localStorage.getItem('sidebar-minimized') == 'true')
  const { newVersionAvailable, setNewVersionAvailable } = useAppStore();

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

  const newVersionPrompt = (
    <a
      className="group hover:bg-[rgba(0,0,0,.525)] text-[#bfbcff] hover:text-white transition-colors backdrop-blur cursor-pointer leading-5 bottom-2 absolute left-[107%] text-base py-2 px-3 rounded-lg bg-slate-900"
      onClick={() => {
        BrowserOpenURL("https://github.com/GreenSoap/cfn-tracker/releases");
        setNewVersionAvailable(false);
      }}
    >
      <RxUpdate className="group-hover:text-white inline text-[#49b3f5] transition-colors w-4 h-4 mr-2" />
      {t(`newVersionAvailable`)}
    </a>
  )

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
              <VscChromeClose className="text-red-800 group-hover:opacity-100 opacity-0 transition-all" />
            </button>
            <button
              aria-label="close"
              className="p-[2px] w-[14px] h-[14px] group-hover:bg-yellow-500 bg-slate-600 flex items-center justify-center rounded-full transition-all"
              onClick={WindowMinimise}
            >
              <VscChromeMinimize className="text-yellow-800 group-hover:opacity-100 opacity-0 transition-all" />
            </button>
          </div>
        </div>
      </header>
      <nav className="mt-5 w-full">
        <ul>
          <li>
            <SidebarLink
              Icon={RiSearch2Line}
              link=""
              name={t("tracking")}
              isSelected={location.pathname == "/"}
              SelectedIcon={RiSearch2Fill}
              isMinimized={isMinimized}
            />
          </li>
          <li>
            <SidebarLink
              Icon={IoDocumentTextOutline}
              link="history"
              name={t("history")}
              isSelected={location.pathname == "/history"}
              SelectedIcon={IoDocumentText}
              isMinimized={isMinimized}
            />
          </li>
        </ul>
      </nav>
      <footer className={`grid w-full text-xl px-2`}>
        <LanguageSelector isMinimized={isMinimized} />

        {/* Twitter */}
        <a
          target="#"
          onClick={() => BrowserOpenURL("https://twitter.com/greensoap_")}
          className={`h-[28px] cursor-pointer w-full group font-extralight flex justify-between items-center mt-1 text-[#d6d4ff] hover:text-white transition-colors`}
        >
          <span className={`flex items-center justify-between lowercase`}>
            <FaTwitter className="text-[#49b3f5] w-4 h-4 mr-2 transition-colors group-hover:text-white" />
            <HideableText text={t('follow')} hide={isMinimized}/>
          </span>
          <FaArrowUp className="relative right-[-8px] w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity" style={{ transform: "rotate(45deg)" }} />
        </a>

        {/* Minimize */}
        <button 
          type='button'
          className={`h-[28px] cursor-pointer w-full group font-extralight flex items-center mt-1 text-[#d6d4ff] hover:text-white transition-colors`}
          onClick={() => setMinimized(!isMinimized)}>
          <FaChevronLeft
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
        {newVersionAvailable && newVersionPrompt}
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