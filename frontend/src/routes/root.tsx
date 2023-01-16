import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { GiDeerTrack } from "react-icons/gi";
import { FaStop } from "react-icons/fa";
import { Track, IsTracking, StopTracking  } from "../../wailsjs/go/main/App";

const Root = () => {
  const { t } = useTranslation();
  const [isLoading, setLoading] = useState(false);
  const [isCurrentlyTracking, setTracking] = useState(false);

  useEffect(() => {
    if (isCurrentlyTracking == true) return;

    const fetchIsTracking = async () => {
      const getIsTracking = await IsTracking();
      setTracking(getIsTracking);
    };

    fetchIsTracking();
  }, []);

  return (
    <div className='grid grid-rows-[0fr_1fr]'>
      <header className='border-b border-slate-50 border-opacity-10 --wails-draggable'>
        <h2 className="pt-4 px-8 pl-12 flex items-center justify-between gap-5 uppercase text-sm tracking-widest mb-4">
            {isCurrentlyTracking || isLoading ? "Tracking" : t("startTracking")}
            {(isCurrentlyTracking || isLoading) && (
              <div
                className="animate-spin inline-block w-5 h-5 border-[3px] border-current border-t-transparent text-pink-600 rounded-full"
                role="status"
                aria-label="loading"
              >
              </div>
            )}
          </h2>
      </header>
      <div className="z-40 h-full select-none grid place-items-center">
        {isCurrentlyTracking && (
          <button
            onClick={(() => {
              StopTracking()
              setTracking(false)
              setLoading(false)
            }
            )}
            type="button"
            className="mt-10 flex items-center justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
          >
            <FaStop className='mr-3'/> Stop
          </button>
        )}
        {!(isCurrentlyTracking || isLoading) && (
          <form
            className="grid justify-items-start mt-10"
            onSubmit={(e) => {
              e.preventDefault();
              const cfn = (e.target as any).cfn.value;
              if (cfn == "") return;
              setLoading(true);

              const x = async () => {
                const isTracking = await Track(cfn);
                setTracking(isTracking);
                if (isTracking == false) {
                  alert("Failed to track CFN")
                } else {
                  console.log('is Tracking')
                  setTracking(true)
                }
                setLoading(false);
              };
              x();
            }}
          >
            <input
              disabled={isLoading}
              type="text"
              name="cfn"
              className="py-3 px-4 block w-full border-gray-200 rounded-md text-lg focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
              placeholder={t("cfnName")!}
            ></input>
            <button
              disabled={isLoading}
              type="submit"
              className="mt-5 justify-self-end flex items-center justify-between bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
            >
              <GiDeerTrack className='mr-3'/> {t("start")}
            </button>
          </form>
        )}
      </div>
    </div>
  );
};

export default Root;
