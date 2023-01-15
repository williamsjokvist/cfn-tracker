import { useTranslation } from "react-i18next";
const Match = () => {
  return (
    <li className="w-full h-9 bg-gray-200 rounded-md dark:bg-gray-700"></li>
  )
}
const History = () => {
  const { t } = useTranslation();

  return (
    <>
      <div className="w-full max-w-lg z-40 grid gap-4">
        <h2 className="uppercase text-sm tracking-widest">
          {t('history')}
        </h2>
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
    </>
  );
};

export default History;
