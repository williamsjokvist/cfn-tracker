import { FaGlobe, FaChevronLeft } from "react-icons/fa";
import { useTranslation } from "react-i18next";
import { APP_LANGUAGES } from "@/i18n";

type LanguageSelectorProps = {
  isMinimized: boolean
}
export const LanguageSelector: React.FC<LanguageSelectorProps> = ({ isMinimized }) => {
  const { t, i18n } = useTranslation();

  return (
    <div className='group flex'>
      <button
        type="button"
        className="h-[28px] font-extralight lowercase relative flex justify-center items-center text-[#d6d4ff] group-hover:text-white transition-colors"
      >
        <FaGlobe className="w-4 h-4 text-[#d6d4ff] group-hover:text-white text-2xl transition-all mr-2" />
        {!isMinimized && t("language")}
      </button>
      <div className="absolute left-[98%] flex group-hover:opacity-100 group-hover:visible invisible opacity-0 transition-all">
        <FaChevronLeft
          className="text-white w-3 h-3 relative right-4 top-2"
          style={{ transform: "rotate(180deg)" }}
        />
        <ul className="w-[195px] text-[16px] uppercase italic flex justify-between bg-black px-3 py-2 relative bottom-2 rounded-lg">
          {APP_LANGUAGES.map((lng, index) => {
            return (
              <li key={lng.code}>
                <button
                  type="button"
                  style={{
                    fontWeight: i18n.resolvedLanguage === lng.code ? "600" : "normal",
                  }}
                  onClick={() => i18n.changeLanguage(lng.code)}
                >
                  {lng.nativeName}
                </button>
                {index !== APP_LANGUAGES.length - 1 && <i className='mx-2'>|</i>}
              </li>
            );
          })}
        </ul>
      </div>
    </div>
  );
};
