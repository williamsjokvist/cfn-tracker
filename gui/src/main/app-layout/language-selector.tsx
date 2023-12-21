import React from "react";
import { useTranslation } from "react-i18next";
import { Icon } from "@iconify/react";

import { HideableText } from "@/ui/hideable-text";
import { GetSupportedLanguages } from "@@/go/core/CommandHandler";

type LanguageSelectorProps = {
  isMinimized: boolean;
};
export const LanguageSelector: React.FC<LanguageSelectorProps> = ({
  isMinimized,
}) => {
  const { t, i18n } = useTranslation();

  const [langs, setLangs] = React.useState<string[]>([]);

  React.useEffect(() => {
    GetSupportedLanguages().then((langs) => setLangs(langs))
  }, [])
  
  const changeLanguage = (code: string) => {
    i18n.changeLanguage(code);
    window.localStorage.setItem("lng", code);
  };

  return (
    <div className="group flex">
      <button
        type="button"
        className="relative flex justify-center items-center h-[28px] text-[#d6d4ff] font-thin lowercase group-hover:text-white transition-colors"
      >
        <Icon
          icon="fa6-solid:globe"
          className="group-hover:text-white text-[#d6d4ff] mr-2 w-4 h-4 transition-all"
        />
        <HideableText text={t("language")} hide={isMinimized} />
      </button>
      <div className="absolute left-[98%] flex group-hover:opacity-100 group-hover:visible invisible opacity-0 transition-all">
        <Icon
          icon="fa6-solid:chevron-left"
          className="text-white w-3 h-3 relative right-4 top-2"
          style={{ transform: "rotate(180deg)" }}
        />
        <ul className="relative bottom-1 w-[195px] text-[16px] uppercase flex gap-2 justify-between group hover:bg-[rgba(0,0,0,.525)] text-[#bfbcff] transition-colors backdrop-blur leading-5 text-base py-2 px-3 rounded-lg bg-[rgba(0,0,0,.625)]">
          {langs.map((code) => (
            <li key={code}>
              <button
                onClick={() => changeLanguage(code)}
                type="button"
                className="cursor-pointer hover:!text-white transition-colors first-letter:uppercase"
                {...(i18n.resolvedLanguage === code && {
                  style: {
                    fontWeight: 600,
                  },
                })}
              >
                {new Intl.DisplayNames([code], { 
                  type: "language", 
                  languageDisplay: "standard", 
                  style: "narrow", 
                  fallback: "code" 
                  })
                  .of(code).split(/[ |（]/)[0]
                }
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};
