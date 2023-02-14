import { useTranslation } from "react-i18next";

const Statistics = () => {
  const { t } = useTranslation();

  return (
    <main className="grid grid-rows-[0fr_1fr] min-h-screen max-h-screen z-40 flex-1 text-white mx-auto">
      <header
        className="border-b border-slate-50 border-opacity-10 backdrop-blur select-none"
        style={
          {
            "--wails-draggable": "drag",
          } as React.CSSProperties
        }
      >
        <h2 className="pt-4 px-8 flex items-center justify-between gap-5 uppercase text-sm tracking-widest mb-4">
          {t("statistics")}
        </h2>
      </header>
      <div className="relative w-full pt-2 z-40 pb-4">
        <div className="flex flex-col items-center justify-center w-full h-full">
          <p>The goods are in progress</p>
        </div>
      </div>
    </main>
  );
};

export default Statistics;
