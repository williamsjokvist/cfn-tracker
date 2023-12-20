import { CheckForUpdate } from "@@/go/core/CommandHandler";
import { BrowserOpenURL } from "@@/runtime/runtime";
import { Icon } from "@iconify/react";
import clsx from "clsx";
import React from "react";
import { useTranslation } from "react-i18next";

export const UpdatePrompt: React.FC = () => {
  const { t } = useTranslation();
  const [hasUpdate, setHasUpdate] = React.useState(false);

  React.useEffect(() => {
    CheckForUpdate().then((hasUpdate: boolean) => setHasUpdate(hasUpdate))
  }, [])

  if (hasUpdate === false) {
    return null
  }

  return (
    <a
      className={clsx(
        "group absolute z-50 left-0 bottom-2",
        "cursor-pointer leading-5 text-base",
        "bg-[rgba(0,0,0,.625)] hover:bg-[rgba(0,0,0,.525)] text-[#bfbcff] hover:text-white transition-colors backdrop-blur",
        "ml-2 py-2 px-3 rounded-lg"
      )}
      onClick={() => {
        BrowserOpenURL("https://cfn.williamsjokvist.se/");
        setHasUpdate(false);
      }}
    >
      <Icon
        icon="radix-icons:update"
        className="group-hover:text-white inline text-[#49b3f5] transition-colors w-4 h-4 mr-2"
      />
      {t("newVersionAvailable")}
    </a>
  )
};