import React from "react";
import { useTranslation } from "react-i18next";
import { Icon } from "@iconify/react";

import * as Dialog from "@/ui/dialog";
import { Checkbox } from "@/ui/checkbox";
import { ActionButton } from "@/ui/action-button";
import { GetCustomThemeList } from "@@/go/core/CommandHandler";

import defaultTheme from './themes/default.png'
import bladesTheme from './themes/blades.png'
import jaegerTheme from './themes/jaeger.png'
import nordTheme from './themes/nord.png'
import pillsTheme from './themes/pills.png'

type ThemeDialogProps = {
  selectedTheme: string
  onSelect: (theme: string) => void
}

// TODO: To preview themes, inject theme CSS from backend instead of images.
export const ThemeDialog: React.FC<ThemeDialogProps> = ({ selectedTheme, onSelect }) => {
  const [customThemes, setCustomThemes] = React.useState<string[]>([]);
  React.useEffect(() => {
    GetCustomThemeList().then(setCustomThemes).catch(console.error)
  }, [])

  const { t } = useTranslation()

  const ThemeItem = ({ name, img }) => (
    <li>
      <button onClick={() => onSelect(name)}>
        <img src={img} className='object-cover object-top overflow-hidden rounded-md h-[76px] w-[390px]'
          style={{ border: selectedTheme == name ? '2px solid lime' : '2px solid transparent' }}
        />
      </button>
    </li>
  )

  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <ActionButton style={{ filter: "hue-rotate(-180deg)", justifyContent: "center" }}>
          <Icon icon="ph:paint-bucket-fill" className="mr-3 w-6 h-6" />
          {t("selectTheme")}
        </ActionButton>
      </Dialog.Trigger>
      <Dialog.Content title="selectTheme">
        <ul className='mt-2 h-80 w-full pr-2 overflow-y-scroll'>
          {customThemes && customThemes.length > 0 &&
            customThemes.map((name) => (
              <li key={name}>
                <button
                  className="w-full cursor-pointer flex py-1 px-2 items-center text-lg hover:bg-[rgba(255,255,255,0.075)]"
                  onClick={() => onSelect(name)}
                >
                  <Checkbox checked={name === selectedTheme} readOnly />
                  <span className="ml-2 text-center cursor-pointer capitalize">
                    {name}
                  </span>
                </button>
              </li>
            ))
          }
          <ThemeItem name="default" img={defaultTheme} />
          <ThemeItem name="blades" img={bladesTheme} />
          <ThemeItem name="pills" img={pillsTheme} />
          <ThemeItem name="jaeger" img={jaegerTheme} />
          <ThemeItem name="nord" img={nordTheme} />
        </ul>
      </Dialog.Content>
    </Dialog.Root>
  )
}
