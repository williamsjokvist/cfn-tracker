import { useTranslation } from "react-i18next";
import { Icon } from "@iconify/react";

import * as Dialog from "@/ui/dialog";
import { Checkbox } from "@/ui/checkbox";
import { ActionButton } from "@/ui/action-button";
import type { StatOptions } from "./output-page";

type StatsDialogProps = {
  options: StatOptions
  onSelect: (option: string, checked: boolean) => void
}

export const StatsDialog: React.FC<StatsDialogProps> = ({ onSelect, options }) => {
  const { t } = useTranslation()
  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <ActionButton style={{ filter: "hue-rotate(-45deg)" }}>
          <Icon icon="bx:stats" className="mr-3 w-6 h-6" />
          {t("displayStats")}
        </ActionButton>
      </Dialog.Trigger>
      <Dialog.Content title="displayStats" description="statsWillBeDisplayed">
        <ul className="overflow-y-scroll h-72 mt-4">
          {Object.entries(options).map(([key, value]) => {
            if (key == "theme") return null;
            return (
              <li key={key}>
                <button
                  className="w-full cursor-pointer flex py-1 px-2 items-center text-lg hover:bg-[rgba(255,255,255,0.075)]"
                  onClick={() => onSelect(key, !value)}
                >
                  <Checkbox checked={options[key] == true} readOnly />
                  <span className="ml-2 text-center cursor-pointer capitalize">
                    {key}
                  </span>
                </button>
              </li>
            );
          })}
        </ul>
      </Dialog.Content>
    </Dialog.Root>
  )
}