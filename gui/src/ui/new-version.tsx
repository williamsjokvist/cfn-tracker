import { Icon } from "@iconify/react";
import { useTranslation } from "react-i18next";
import { BrowserOpenURL } from "@@/runtime/runtime";

type NewVersionPromptProps = {
  hasNewVersion: boolean;
  setNewVersion: (hasNewVersion: boolean) => void
}

export const NewVersionPrompt: React.FC<NewVersionPromptProps> = ( { hasNewVersion, setNewVersion }) => {
  const { t } = useTranslation()
  if (!hasNewVersion) return null
  return (
    <a
      className="z-50 absolute left-full ml-2 bottom-2 group hover:bg-[rgba(0,0,0,.525)] text-[#bfbcff] hover:text-white transition-colors backdrop-blur cursor-pointer leading-5 text-base py-2 px-3 rounded-lg bg-[rgba(0,0,0,.625)]"
      onClick={() => {
        BrowserOpenURL("https://williamsjokvist.github.io/cfn-tracker/");
        setNewVersion(false);
      }}
    >
      <Icon icon='radix-icons:update' className="group-hover:text-white inline text-[#49b3f5] transition-colors w-4 h-4 mr-2" />
      {t(`newVersionAvailable`)}
    </a>
  )
}