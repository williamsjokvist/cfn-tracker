import { Icon } from "@iconify/react";
import { Quit, WindowMinimise } from "@@/runtime";

export const AppTitleBar: React.FC = () => (
  <header style={{ "--draggable": "drag" } as React.CSSProperties}>
    <div className="flex justify-start">
      <div className="group mt-2 group ml-2 flex mb-3">
        <button
          aria-label="close"
          className="mr-[8px] p-[2px] w-[14px] h-[14px] group-hover:bg-red-500 bg-slate-600 flex items-center justify-center rounded-full transition-all"
          onClick={Quit}
        >
          <Icon
            icon="ep:close-bold"
            className="text-red-800 group-hover:opacity-100 opacity-0 transition-all"
          />
        </button>
        <button
          aria-label="close"
          className="p-[2px] w-[14px] h-[14px] group-hover:bg-yellow-500 bg-slate-600 flex items-center justify-center rounded-full transition-all"
          onClick={WindowMinimise}
        >
          <Icon
            icon="mingcute:minimize-fill"
            className="text-yellow-800 group-hover:opacity-100 opacity-0 transition-all"
          />
        </button>
      </div>
    </div>
  </header>
);
