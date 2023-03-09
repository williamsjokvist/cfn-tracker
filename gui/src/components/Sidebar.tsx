import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { FaTwitter, FaChevronLeft, FaArrowUp } from "react-icons/fa";
import { RxUpdate } from "react-icons/rx";
import {
  IoDocumentTextOutline,
  IoDocumentText,
} from "react-icons/io5";
import { RiSearch2Line, RiSearch2Fill } from "react-icons/ri";
import { VscChromeMinimize, VscChromeClose } from "react-icons/vsc";
import LanguageSelector from "./LanguageSelector";
import { useLocation } from "react-router-dom";
import { BrowserOpenURL, Quit, WindowMinimise } from "../../wailsjs/runtime";

import { GetAppVersion } from "../../wailsjs/go/core/App";
import { useAppStore } from "../store/use-app-store";

const Link = (props: {
  Icon: any;
  link: string;
  name: string;
  isSelected?: boolean;
  SelectedIcon?: any;
  showArrow: boolean;
}) => {
  const { Icon, link, name, isSelected, SelectedIcon, showArrow } = props;
  return (
    <li>
      <a
        href={"#/" + link}
        className="text-lg text-[#bfbcff] text-opacity-80 rounded py-2 px-1 group flex items-center justify-between hover:!text-white hover:bg-slate-50 hover:bg-opacity-5 transition-colors"
        style={{
          fontWeight: isSelected ? "600" : "200",
          color: isSelected ? "#d6d4ff" : "#bfbcff",
        }}
      >
        <span className="flex items-center justify-between">
          {isSelected && (
            <SelectedIcon className="text-[#f85961] transition-colors w-10 h-7 mr-1" />
          )}
          {!isSelected && (
            <Icon className="text-[#f85961] transition-colors w-10 h-7 mr-1" />
          )}
          {name}
        </span>
        <FaChevronLeft
          className="w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity"
          style={{ transform: "rotate(180deg)", display: showArrow ? 'block' : 'none' }}
        />
      </a>
    </li>
  );
};
const Sidebar = () => {
  const { t } = useTranslation();
  const location = useLocation();
  const [appVersion, setAppVersion] = useState("");

  const [isMinimized, setMinimized] = useState(false)

  const { newVersionAvailable, setNewVersionAvailable } = useAppStore();
  useEffect(() => {
    setMinimized(false)
    GetAppVersion().then((version: string) => {
      setAppVersion(version);
    });
  }, []);

  return (
    <aside
      className="select-none bg-[rgba(3,5,19,0.33)] border-r border-slate-50 border-opacity-10 backdrop-blur relative z-50 h-screen overflow-auto scrollbar-none grid grid-rows-[0fr_1fr_0fr] py-2 text-white whitespace-nowrap dark:border-slate-50 dark:border-opacity-10"
      style={{
        width: isMinimized ? '75px' : "175px",
        overflow: "visible",
        justifyContent: isMinimized ? 'center' : 'normal',
        paddingLeft: isMinimized ? '0' : '10px',
        paddingRight: isMinimized ? '0' : '10px'

      }}
    >
      <header
        style={
          {
            "--wails-draggable": "drag",
          } as React.CSSProperties
        }
      >
        <div className="flex justify-start">
          <div className="group mt-2 group ml-2 flex mb-3">
            <button
              aria-label="close"
              className="mr-[8px] p-[2px] w-[14px] h-[14px] group-hover:bg-red-500 bg-slate-600 flex items-center justify-center rounded-full transition-all"
              onClick={() => Quit()}
            >
              <VscChromeClose className="text-red-800 group-hover:opacity-100 opacity-0 transition-all" />
            </button>
            <button
              aria-label="close"
              className="p-[2px] w-[14px] h-[14px] group-hover:bg-yellow-500 bg-slate-600 flex items-center justify-center rounded-full transition-all"
              onClick={() => WindowMinimise()}
            >
              <VscChromeMinimize className="text-yellow-800 group-hover:opacity-100 opacity-0 transition-all" />
            </button>
          </div>
        </div>
      </header>
      <nav className="mt-5 w-full">
        <ul>
          <Link
            Icon={RiSearch2Line}
            link=""
            name={isMinimized ? '' : t("tracking")}
            isSelected={location.pathname == "/"}
            SelectedIcon={RiSearch2Fill}
            showArrow={!isMinimized}
          />
          <Link
            Icon={IoDocumentTextOutline}
            link="history"
            name={isMinimized ? '' : t("history")}
            isSelected={location.pathname == "/history"}
            SelectedIcon={IoDocumentText}
            showArrow={!isMinimized}
          />
        </ul>
      </nav>
      <footer className={`grid w-full text-xl px-2`}>
        <LanguageSelector isMinimized={isMinimized} />
        <a
          target="#"
          onClick={() => {
            BrowserOpenURL("https://twitter.com/greensoap_");
          }}
          className={`h-[28px] cursor-pointer w-full group font-extralight flex justify-between items-center mt-1 text-[#d6d4ff] hover:text-white transition-colors`}
        >
          <span className={`flex items-center justify-between`}>
            <FaTwitter className="text-[#49b3f5] w-4 h-4 mr-2 transition-colors group-hover:text-white" />
            {!isMinimized && 'greensoap_'}
          </span>
          <FaArrowUp
            className="relative right-[-8px] w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity"
            style={{ transform: "rotate(45deg)" }}
          />

        </a>
        <button className={`h-[28px] cursor-pointer w-full group font-extralight flex items-center mt-1 text-[#d6d4ff] hover:text-white transition-colors`}
          onClick={() => setMinimized(!isMinimized)}>
          <FaChevronLeft
            className="group-hover:text-white text-[#d6d4ff] w-4 h-4 transition-all"
            style={{ transform: isMinimized ? "rotate(-180deg)" : 'none' }}
          />
          {!isMinimized && (
            <span className="ml-2">{t('minimize')}</span>
          )}
        </button>
        <a
          target="#"
          className="text-sm mt-4 font-extralight cursor-pointer hover:underline"
          onClick={() => {
            BrowserOpenURL("https://github.com/GreenSoap/cfn-tracker/releases");
          }}
        >
          {isMinimized ? 'v' + appVersion : 'CFN Tracker v' + appVersion}
        </a>
        {newVersionAvailable && (
          <a
            className="group hover:bg-[rgba(0,0,0,.525)] text-[#bfbcff] hover:text-white transition-colors backdrop-blur cursor-pointer leading-5 bottom-2 absolute left-[107%] text-base py-2 px-3 rounded-lg bg-slate-900"
            onClick={() => {
              BrowserOpenURL(
                "https://github.com/GreenSoap/cfn-tracker/releases"
              );
              setNewVersionAvailable(false);
            }}
          >
            <RxUpdate className="group-hover:text-white inline text-[#49b3f5] transition-colors w-4 h-4 mr-2" />
            {t(`newVersionAvailable`)}
          </a>
        )}
      </footer>
    </aside>
  );
};

export default Sidebar;
