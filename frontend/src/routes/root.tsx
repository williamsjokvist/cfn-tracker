import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { GiDeerTrack } from "react-icons/gi";
import { FaStop } from "react-icons/fa";
import { SetCFN, IsTracking, StopTracking  } from "../../wailsjs/go/main/App";

const Root = () => {
  const { t } = useTranslation();
  const [isLoading, setLoading] = useState(false);
  const [isCurrentlyTracking, setTracking] = useState(false);

  useEffect(() => {
    if (isCurrentlyTracking == true) return;

    const fetchIsTracking = async () => {
      const getIsTracking = await IsTracking();
      console.log(getIsTracking);
      setTracking(getIsTracking);
    };

    fetchIsTracking();
  }, []);

  return (
    <>
      <div className="w-fit max-w-md z-40 select-none">
        <h2 className="flex items-center gap-5 uppercase text-sm tracking-widest mb-4">
          {isCurrentlyTracking || isLoading ? "Tracking" : t("startTracking")}
          {(isCurrentlyTracking || isLoading) && (
            <div
              className="animate-spin inline-block w-8 h-8 border-[3px] border-current border-t-transparent text-pink-600 rounded-full"
              role="status"
              aria-label="loading"
            >
              <span className="sr-only">Loading...</span>
            </div>
          )}
        </h2>
        {isCurrentlyTracking && (
          <button
            onClick={(() => {
              StopTracking()
              setTracking(false)
              setLoading(false)
            }

            )}
            type="button"
            className="justify-self-end flex items-center justify-between gap-3 bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
          >
            <FaStop /> Stop
          </button>
        )}
        {!(isCurrentlyTracking || isLoading) && (
          <form
            className="grid gap-4 justify-items-start"
            onSubmit={(e) => {
              e.preventDefault();
              const cfn = (e.target as any).cfn.value;
              if (cfn == "") return;
              setLoading(true);

              const x = async () => {
                const res = await SetCFN(cfn);
                console.log(res);
                const getIsTracking = await IsTracking();
                console.log(getIsTracking);
                setTracking(getIsTracking);
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
              className="justify-self-end flex items-center justify-between gap-3 bg-[rgba(255,10,10,.1)] rounded-md px-5 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
            >
              <GiDeerTrack /> {t("start")}
            </button>
          </form>
        )}
      </div>
    </>
  );
};

export default Root;
