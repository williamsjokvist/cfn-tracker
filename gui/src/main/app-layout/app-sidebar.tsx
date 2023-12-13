import React from "react";
import { useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useAnimate } from "framer-motion";
import { Icon } from "@iconify/react";

import { GetAppVersion } from "@@/go/core/CommandHandler";
import { LanguageSelector } from "@/i18n/language-selector";
import { BrowserOpenURL } from "@@/runtime";

import { NavigationLink } from "@/ui/nav-link";
import { HideableText } from "@/ui/hideable-text";

import { AppTitleBar } from "./app-titlebar";

import "@/styles/sidebar.sass";

const NavigationItems = [
  {
    icon: "ri:search-line",
    selectedIcon: "ri:search-fill",
    href: "tracking",
    main: true,
  },
  {
    icon: "ion:document-text-outline",
    selectedIcon: "ion:document-text",
    href: "sessions",
    main: false,
  },
  {
    icon: "clarity:sign-out-line",
    selectedIcon: "clarity:sign-out-solid",
    href: "output",
    main: false,
  },
];

export const AppSidebar: React.FC = () => {
  const { t } = useTranslation();
  const location = useLocation();
  const [appVersion, setAppVersion] = React.useState("");
  const [isMinimized, setMinimized] = React.useState(
    window.localStorage.getItem("sidebar-minimized") == "true"
  );

  const [scope, animate] = useAnimate();

  React.useEffect(() => {
    animate("a, button", { opacity: [0, 1] }, { delay: 0.25 });
  }, []);

  React.useEffect(() => {
    !appVersion && GetAppVersion().then((v) => setAppVersion(v));
  }, [appVersion]);

  React.useEffect(() => {
    if (isMinimized) {
      localStorage.setItem("sidebar-minimized", "true");
    } else if (!isMinimized) {
      localStorage.removeItem("sidebar-minimized");
    }
  }, [isMinimized]);

  return (
    <aside
      ref={scope}
      className="sidebar"
      style={{
        width: isMinimized ? "76px" : "175px",
      }}
    >
      <AppTitleBar />

      <nav className="mt-5 w-full">
        {NavigationItems.map((navItem) => (
          <NavigationLink
            key={navItem.href}
            href={navItem.main ? "" : navItem.href}
            icon={navItem.icon}
            name={t(navItem.href)}
            selectedIcon={navItem.selectedIcon}
            isSelected={
              navItem.main
                ? location.pathname == "/"
                : location.pathname.includes(navItem.href)
            }
            isMinimized={isMinimized}
          />
        ))}
      </nav>

      <footer className="grid w-full text-xl px-2">
        <LanguageSelector isMinimized={isMinimized} />

        {/* Twitter */}
        <a
          target="#"
          onClick={() => BrowserOpenURL("https://twitter.com/greensoap_")}
          className={`h-[28px] cursor-pointer w-full group font-extralight flex justify-between items-center mt-1 text-[#d6d4ff] hover:text-white transition-colors`}
        >
          <span className={`flex items-center justify-between lowercase`}>
            <Icon
              icon="fa6-brands:twitter"
              className="text-[#49b3f5] w-4 h-4 mr-2 transition-colors group-hover:text-white"
            />
            <HideableText text={t("follow")} hide={isMinimized} />
          </span>
          <Icon
            icon="fa6-solid:arrow-up"
            className="relative right-[-8px] w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity"
            style={{ transform: "rotate(45deg)" }}
          />
        </a>

        {/* Minimize */}
        <button
          type="button"
          className={`h-[28px] cursor-pointer w-full group font-extralight flex items-center mt-1 text-[#d6d4ff] hover:text-white transition-colors`}
          onClick={() => setMinimized(!isMinimized)}
        >
          <Icon
            icon="fa6-solid:chevron-left"
            className="group-hover:text-white text-[#d6d4ff] w-4 h-4 transition-all"
            style={{ transform: isMinimized ? "rotate(-180deg)" : "none" }}
          />
          <HideableText
            className="ml-2"
            text={t("minimize")}
            hide={isMinimized}
          />
        </button>

        {/* Version */}
        <a
          target="#"
          className="text-sm mt-4 font-extralight cursor-pointer hover:underline"
          onClick={() =>
            BrowserOpenURL("https://github.com/GreenSoap/cfn-tracker/releases")
          }
        >
          {isMinimized ? `v${appVersion}` : `CFN Tracker v${appVersion}`}
        </a>
      </footer>
    </aside>
  );
};
