import React from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { CFNMachineContext } from "@/machine";
import { GetAvailableLogs } from "@@/go/core/CommandHandler";

type CfnFormValues = {
  cfn: string;
  restore?: boolean;
}
export const CFNForm: React.FC = () => {
  const { t } = useTranslation();
  const [_, send] = CFNMachineContext.useActor();
  const [oldCfns, setOldCfns] = React.useState<string[] | null>(null);
  const [lastJSONExist, setLastJSONExist] = React.useState<boolean>(false);
  const cfnInputRef = React.useRef<HTMLInputElement>(null)

  React.useCallback(() => {
    GetAvailableLogs().then(logs => setOldCfns(logs))
    console.log(oldCfns)
  }, [oldCfns])

  const { register, handleSubmit, control } = useForm<CfnFormValues>()
/*
  const { cfn, restore } = useWatch({
    control
  })
*/
  const onSubmit = (values: CfnFormValues) => {
    const cfn = cfnInputRef.current.value
    if (!cfn || cfn == '') return
    send({
      type: 'submit',
      cfn,
      restore: false
    })
  }

  return (
    <form
      className="max-w-[450px] mx-auto"
      onSubmit={handleSubmit(onSubmit)}
    >
      <h3 className="mb-2 text-lg">{t("enterCfnName")}:</h3>
      <input
        {...register('cfn')}
        ref={cfnInputRef}
        className="bg-transparent border-b-2 border-0 focus:ring-offset-transparent focus:ring-transparent border-b-[rgba(255,255,255,0.275)] focus:border-white hover:border-white outline-none focus:outline-none hover:text-white transition-colors py-3 px-4 block w-full text-lg text-gray-300"
        type="text"
        placeholder={t("cfnName")!}
        autoCapitalize="off"
        autoComplete="off"
        autoCorrect="off"
        autoSave="off"
      />
      {oldCfns && (
        <div className="mt-3 flex flex-wrap gap-2 content-center items-center text-center pr-3">
          {oldCfns.map((cfn, index) => {
            return (
              <button
                disabled={false}
                onClick={_ => {
                  cfnInputRef.current.value = cfn
                }}
                className="whitespace-nowrap bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-base backdrop-blur rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
                type="button"
                key={index}
              >
                {cfn}
              </button>
            );
          })}
        </div>
      )}
      {lastJSONExist && (
        <div className={`text-lg flex items-center mt-4`}>
          <input
            {...register('restore')}
            type="checkbox"
            className="w-7 h-7 rounded-md checked:border-2 checked:focus:border-[rgba(255,255,255,.25)] checked:hover:border-[rgba(255,255,255,.25)] checked:border-[rgba(255,255,255,.25)] border-2 border-[rgba(255,255,255,.25)] focus:border-2 cursor-pointer bg-transparent text-transparent focus:ring-offset-transparent focus:ring-transparent mr-4"
            onChange={e => {
              if (e.target.checked)
                cfnInputRef.current.value = ''
            }}
          />
          {t("restoreSession")}
        </div>
      )}
      <div className="flex justify-end">
        <button
          type="submit"
          className="mt-4 select-none text-center bg-[rgba(255,10,10,.1)] rounded-md px-7 py-3 border-[#FF3D51] hover:bg-[#FF3D51] border-[1px] transition-colors font-semibold text-md"
          style={{ filter: "hue-rotate(156deg)" }}
        >
          {t("start")}
        </button>
      </div>
    </form>
  )
}
