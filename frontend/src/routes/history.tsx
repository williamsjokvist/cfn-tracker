import { useTranslation } from "react-i18next";
const Match = () => {
  return (
    <li className="w-full h-9 bg-gray-200 rounded-md dark:bg-gray-700"></li>
  )
}
const History = () => {
  const { t } = useTranslation();

  return (
    <main className="grid grid-rows-[0fr_1fr] min-h-screen">
      <header className='border-b border-slate-50 border-opacity-10 --wails-draggable select-none'>
        <h2 className="pt-4 px-8 pl-12 flex items-center justify-between gap-5 uppercase text-sm tracking-widest mb-4">
          {t('history')}
        </h2>
      </header>
      <div className="w-full z-40 grid justify-items-center justify-center pt-10">
        <ul className="space-y-3 animate-pulse min-w-[525px] max-h-[300px] overflow-y-scroll">
          <Match />
          <Match />
          <Match />
          <Match />
          <Match />
          <Match />
          <Match />
          <Match />
          <Match />
          <Match />
        </ul>
      </div>
    </main>
  );
};

export default History;
