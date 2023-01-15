import { useState } from "react";
import { useTranslation } from "react-i18next";
import { FaTwitter, FaChevronLeft, FaScroll } from "react-icons/fa";
import { GiMagnifyingGlass } from "react-icons/gi";

import appLogo from "../assets/logo.png";
import LanguageSelector from "./LanguageSelector";

const Link = (props: { Icon: any; link: string; name: string }) => {
  const { Icon, link, name } = props;
  return (
    <li className="odd:border-b-[1px] border-[#181a3b]">
      <a
        href={'#/'+link}
        className="py-2 pl-2 pr-1 group flex items-center justify-between gap-1 text-[#d6d4ff] hover:text-white hover:bg-[#2b2956] transition-colors"
      >
        <span className="flex items-center justify-between gap-1 text-xl tracking-wider">
          <Icon className="group-hover:text-[#e8e6ff] text-[#d6d4ff] transition-colors w-10 h-6" />
          {name}
        </span>
        <FaChevronLeft
          className="w-4 h-4 group-hover:opacity-100 opacity-0 transition-opacity"
          style={{ rotate: "180deg" }}
        />
      </a>
    </li>
  );
};
const Sidebar = () => {
  const [isOpen, setIsOpen] = useState(true);
  const { t } = useTranslation();

  return (
    <div className="relative grid z-50">
      <button
        type="button"
        className="group flex items-center gap-2 absolute top-[5px] transition-all"
        style={{
          left: isOpen ? '85%' : '104%',
        }}
        onClick={() => {
          setIsOpen(!isOpen);
        }}
      >
        <FaChevronLeft
          className="text-[#d6d4ff] group-hover:!text-[#e8e6ff] text-2xl transition-all"
          style={{
            rotate: isOpen ? "0deg" : "-180deg",
            color: isOpen ? '#6b6d93' : '#d6d4ff'
          }}
        />
        <span
          className="group-hover:!opacity-100 group-hover:!visible text-white text-xs transition-all relative top-[1px]"
          style={{
            ...(isOpen && {
              opacity: 0,
              visibility: 'hidden'
            })
          }}
        >
          {isOpen ? 'CLOSE' : 'MENU'}
        </span>
      </button>
      <aside
        className="grid grid-rows-[0fr_1fr_0fr] py-2 bg-[#222338] text-white whitespace-nowrap transition-all"
        style={{
          width: isOpen ? "190px" : "0px",
          overflow: isOpen ? 'visible' : 'hidden'
        }}
      >
        <img
          src={appLogo}
          className="max-w-[100px] w-full relative left-[20px] filter-["
          alt="CFN Tracker logo"
        />
        <nav className="mt-5 w-full">
          <ul className="text-2xl font-semibold">
            <Link Icon={GiMagnifyingGlass} link="tracking" name={t('tracking')} />
            <Link Icon={FaScroll} link="history" name={t('history')} />
          </ul>
        </nav>
        <footer className="grid px-5 w-full text-xl">
          <LanguageSelector />
          <a
            target='_blank'
            href="https://twitter.com/greensoap_"
            className="font-extralight flex items-center justify-between gap-1 mt-1 text-[#d6d4ff] hover:text-white transition-colors"
          >
            <span className='flex items-center justify-between gap-2'>
              <FaTwitter className="text-[#f54952] w-4 h-4" />
              greensoap_
            </span>
          </a>
          <small className='text-sm mt-4 font-extralight'>CFN Tracker v2.1.0</small>
        </footer>
      </aside>
    </div>
  );
};

export default Sidebar;
