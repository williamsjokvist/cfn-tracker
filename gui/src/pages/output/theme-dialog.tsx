import { useTranslation } from "react-i18next";
import { useLoaderData } from "react-router-dom";
import { Icon } from "@iconify/react";

import * as Dialog from "@/ui/dialog";
import { ActionButton } from "@/ui/action-button";

import defaultTheme from './themes/default.png'
import bladesTheme from './themes/blades.png'
import jaegerTheme from './themes/jaeger.png'
import nordTheme from './themes/nord.png'
import pillsTheme from './themes/pills.png'

type ThemeDialogProps = {
  selectedTheme: string
  onSelect: (theme: string) => void
}

export const ThemeDialog: React.FC<ThemeDialogProps> = ({ selectedTheme, onSelect }) => {
  const themes = useLoaderData() as string[];

  const { t } = useTranslation()

  const ThemeItem = ({ name, img }) => (
    <li>
      <button onClick={() => onSelect(name)}>
        <img src={img} className='object-cover object-top overflow-hidden w-full rounded-md h-[76px]' 
          style={{ border: selectedTheme == name ? '2px solid lime' : '2px solid transparent'}}
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
          <ThemeItem name="default" img={defaultTheme} />
          <ThemeItem name="blades" img={bladesTheme} />
          <ThemeItem name="pills" img={pillsTheme}  />
          <ThemeItem name="jaeger" img={jaegerTheme}  />
          <ThemeItem name="nord" img={nordTheme} />
        </ul>
      </Dialog.Content>
  </Dialog.Root>
  )
}
