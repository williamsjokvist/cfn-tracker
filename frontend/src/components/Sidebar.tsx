import { useState } from "react";
import { useTranslation } from "react-i18next";
import { FaTwitter, FaChevronLeft, FaArrowUp } from "react-icons/fa";
import { IoDocumentTextOutline, IoDocumentText, IoPlay, IoPlayOutline } from 'react-icons/io5'
import { RiSearch2Line, RiSearch2Fill } from 'react-icons/ri'

import appLogo from "../assets/logo.png";
import LanguageSelector from "./LanguageSelector";
import { useLocation } from "react-router-dom";

const Link = (props: { Icon: any, link: string; name: string, isSelected?: boolean, SelectedIcon?: any, }) => {
  const { Icon, link, name, isSelected, SelectedIcon } = props;
  return (
    <li className="">
      <a
        href={'#/'+link}
        className="text-lg text-slate-50 text-opacity-80 rounded py-2 px-1 mx-2 group flex items-center justify-between hover:!text-[#d6d4ff] hover:bg-slate-50 hover:bg-opacity-5 transition-colors"
        style={{
          fontWeight: isSelected ? '600' : '200',
          color: isSelected ? '#d6d4ff' : 'rgb(248 250 252)'
        }}
      >
        <span className="flex items-center justify-between">
          {isSelected && <SelectedIcon className="text-[#f85961] transition-colors w-10 h-7 mr-1" />}
          {!isSelected && <Icon className="text-[#f85961] transition-colors w-10 h-7 mr-1" />}
          {name}
        </span>
        <FaChevronLeft
          className="w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity"
          style={{ transform: "rotate(180deg)" }}
        />
      </a>
    </li>
  );
};
const Sidebar = () => {
  const [isOpen, setIsOpen] = useState(true);
  const { t } = useTranslation();
  const location = useLocation();

  return (
    <div className="relative grid z-50">
{/*       <button
        type="button"
        className="group flex items-center gap-2 absolute top-[45%] transition-all"
        style={{
          left: isOpen ? '97%' : '10px',
        }}
        onClick={() => {
          setIsOpen(!isOpen);
        }}
      >
        <FaChevronLeft
          className="text-[#d6d4ff] group-hover:!text-[#e8e6ff] text-2xl transition-all"
          style={{
            transform: isOpen ? "rotate(0deg)" : "rotate(-180deg)",
            color: isOpen ? '#6b6d93' : '#d6d4ff'
          }}
        />
        <span
          className="group-hover:opacity-100 group-hover:visible invisible opacity-0 text-white text-xs transition-all relative top-[1px]"
        >
          {isOpen ? 'CLOSE' : 'MENU'}
        </span>
      </button>*/}
      <aside
        className="h-full overflow-auto scrollbar-none grid grid-rows-[0fr_1fr_0fr] py-2 bg-[#222338] text-white whitespace-nowrap transition-all border-slate-900 border-opacity-10 dark:border-slate-50 dark:border-opacity-10"
        style={{
          width: isOpen ? "190px" : "0px",
          overflow: isOpen ? 'visible' : 'hidden'
        }}
      >
        <header className=''>
          <div className='pt-5 mb-5 --wails-draggable'/>
            <img
            src={appLogo}
            className="max-w-[65px] w-full relative left-[20px]"
            alt="CFN Tracker logo"
            />
        </header>
        <nav className="mt-5 w-full">
          <ul className="">
            <Link Icon={RiSearch2Line} link="tracking" name={t('tracking')} isSelected={(location.pathname == '/tracking')} SelectedIcon={RiSearch2Fill}/>
            <Link Icon={IoDocumentTextOutline} link="history" name={t('history')} isSelected={(location.pathname == '/history')} SelectedIcon={IoDocumentText}/>
          </ul>
        </nav>
        <footer className="grid px-5 w-full text-xl">
          <LanguageSelector />
          <a
            target='_blank'
            href="https://twitter.com/greensoap_"
            className="w-full group font-extralight flex items-center justify-between mt-1 text-[#d6d4ff] hover:text-white transition-colors"
          >
            <span className='flex items-center justify-between'>
              <FaTwitter className="text-[#49b3f5] w-4 h-4 mr-2" />
              greensoap_
            </span>
            <FaArrowUp
                className="relative right-[-8px] w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity"
                style={{ transform: "rotate(45deg)" }}
            />
          </a>
          <small className='text-sm mt-4 font-extralight'>CFN Tracker v2.1.0</small>
        </footer>
      </aside>
    </div>
  );
};

export default Sidebar;
