import { Icon } from "@iconify/react";
import clsx from "clsx";
import { HideableText } from "./hideable-text";

type NavigationLinkProps = {
  icon: string;
  href: string;
  name: string;
  isMinimized: boolean;
  isSelected?: boolean;
  selectedIcon?: string;
};
export const NavigationLink: React.FC<NavigationLinkProps> = (props) => (
  <a
    href={`#/${props.href}`}
    className={clsx(
      "flex items-center justify-between",
      "group flex items-center justify-between",
      "hover:!text-white hover:bg-slate-50 active:bg-[rgba(255,255,255,.075)] hover:bg-opacity-5 transition-colors",
      "text-lg text-[#bfbcff] text-opacity-80 rounded py-2 px-1"
    )}
    style={{
      fontWeight: props.isSelected ? "600" : "200",
      color: props.isSelected ? "#d6d4ff" : "#bfbcff",
    }}
  >
    <span className="flex items-center justify-between">
      <Icon
        icon={props.isSelected ? props.selectedIcon : props.icon}
        className="text-[#f85961] transition-colors w-10 h-7 mr-1"
      />
      <HideableText text={props.name} hide={props.isMinimized} />
    </span>
    <Icon
      icon="fa6-solid:chevron-left"
      className="w-3 h-3 group-hover:opacity-100 opacity-0 transition-opacity"
      style={{
        transform: "rotate(180deg)",
        display: props.isMinimized ? "none" : "block",
      }}
    />
  </a>
);
