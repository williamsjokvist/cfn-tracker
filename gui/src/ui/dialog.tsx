import React from "react";
import { Icon } from "@iconify/react";
import { useTranslation } from "react-i18next";
import * as DialogPrimitive from "@radix-ui/react-dialog";

import type { LocalizationKey } from "@/main/i18n";

type DialogContentProps = {
  title: LocalizationKey
  description?: LocalizationKey
} & React.PropsWithChildren<DialogPrimitive.DialogContentProps>

export const DialogContent = React.forwardRef((
  props: DialogContentProps, 
  ref: React.MutableRefObject<HTMLDivElement>
) => {
  const { t } = useTranslation();
  return (
    <DialogPrimitive.Portal>
      <DialogPrimitive.Content {...props} ref={ref} className="fixed z-[9999] left-[50%] top-[50%] translate-x-[-50%] translate-y-[-50%] w-[450px] text-white bg-black backdrop-blur-xl bg-opacity-60 px-5 py-6 rounded-3xl text-lg h-[420px] shadow-[0_3px_16px_rgba(0,0,0,.5)]">
        <header className="flex justify-between">
          <div>
            <DialogPrimitive.Title className="text-2xl font-bold">
              {t(props.title)}
            </DialogPrimitive.Title>
            {props.description && 
              <DialogPrimitive.Description className="text-lg leading-none">
                {t(props.description)}
              </DialogPrimitive.Description>
            }
          </div>
          <DialogCloseButton />
        </header>
        {props.children}
      </DialogPrimitive.Content>
      <DialogPrimitive.Overlay className="bg-black bg-opacity-20 fixed top-0 left-0 right-0 bottom-0 grid place-items-center overflow-y-auto" />
    </DialogPrimitive.Portal>
  )
});

const DialogCloseButton = () => (
  <DialogPrimitive.Close aria-label="Close" className="bg-[#202020] hover:bg-[#2b2a33] w-11 h-11 rounded-full">
    <Icon icon="ci:close-big" width={28} className='mx-auto' />
  </DialogPrimitive.Close>
)

export default {
  Root: DialogPrimitive.Root,
  Trigger: DialogPrimitive.Trigger,
  Content: DialogContent
}
