import React from "react";
import { createRoot } from "react-dom/client";
import { Icon } from "@iconify/react";

import { useErrorMessage } from "@/main/app-layout/error-message";

import * as Dialog from "@/ui/dialog";
import { Checkbox } from "@/ui/checkbox";
import { ActionButton } from "@/ui/action-button";

import { GetThemes } from "@@/go/core/CommandHandler";
import type { model } from "@@/go/models";

type ThemeDialogProps = {
  selectedTheme: string
  onSelect: (theme: string) => void
}

export const ThemeDialog: React.FC<ThemeDialogProps> = ({ selectedTheme, onSelect }) => {
  const containerRef = React.useRef<HTMLDivElement>(null)
  const [themes, setThemes] = React.useState<model.Theme[]>([])
  const [isOpen, setOpen] = React.useState(false)
  const setError = useErrorMessage()

  React.useEffect(() => {
    GetThemes().then(setThemes).catch(setError)
  }, [])

  React.useEffect(() => {
    if (!isOpen || !containerRef.current) return

    for (const theme of themes) {
      const container = document.createElement("div");
      containerRef.current.appendChild(container);
      const shadowRoot = container.attachShadow({ mode: "open" });

      createRoot(container).render(
        <div className="mb-4">
          <style>{theme.css.match(/@import.*;/g)}</style>
          <Checkbox id={`${theme.name}-checkbox`} checked={theme.name === selectedTheme} onChange={(e) => onSelect(theme.name)} />
          <label
            htmlFor={`${theme.name}-checkbox`}
            style={{ fontFamily: "League Spartan" }}
            className="capitalize text-lg font-bold cursor-pointer"
          >
            {theme.name}
          </label>
        </div>
      )

      createRoot(shadowRoot).render(
        <>
          <slot />
          <div className="stat-list">
            <style>{theme.css}</style>
            <div className="stat-item">
              <span className="stat-title">MR</span>
              <span className="stat-value">444</span>
            </div>
          </div>
        </>
      )
    }

    return () => {
      if (containerRef.current)
        containerRef.current.innerHTML = ""
    }
  }, [isOpen, themes, selectedTheme])

  return (
    <Dialog.Root onOpenChange={setOpen}>
      <Dialog.Trigger asChild>
        <ActionButton className="capitalize" style={{ filter: "hue-rotate(-180deg)", justifyContent: "center" }}>
          <Icon icon="ph:paint-bucket-fill" className="mr-3 w-6 h-6" />
          {selectedTheme}
        </ActionButton>
      </Dialog.Trigger>
      <Dialog.Content title="selectTheme">
        <div ref={containerRef} className='grid gap-4 mt-2 h-80 w-full pr-2 overflow-y-scroll'></div>
      </Dialog.Content>
    </Dialog.Root>
  )
}
