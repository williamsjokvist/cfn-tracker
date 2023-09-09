import { Icon } from "@iconify/react";
import clsx from "clsx";

type NewVersionPromptProps = {
  onDismiss: () => void;
  text: string;
};

export const UpdatePrompt: React.FC<NewVersionPromptProps> = (props) => (
  <a
    className={clsx(
      "group absolute z-50 left-0 bottom-2",
      "cursor-pointer leading-5 text-base",
      "bg-[rgba(0,0,0,.625)] hover:bg-[rgba(0,0,0,.525)] text-[#bfbcff] hover:text-white transition-colors backdrop-blur",
      "ml-2 py-2 px-3 rounded-lg"
    )}
    onClick={props.onDismiss}
  >
    <Icon
      icon="radix-icons:update"
      className="group-hover:text-white inline text-[#49b3f5] transition-colors w-4 h-4 mr-2"
    />
    {props.text}
  </a>
);
