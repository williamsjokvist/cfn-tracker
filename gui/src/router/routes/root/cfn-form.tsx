import React from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { CFNMachineContext } from "@/machine";
import { GetAvailableLogs, ResultsJSONExist } from "@@/go/core/CommandHandler";
import { ActionButton } from "@/ui/action-button";

type CfnFormValues = {
  cfn: string;
  restore?: boolean;
}
export const CFNForm: React.FC = () => {
  const { t } = useTranslation();
  const [_, send] = CFNMachineContext.useActor();

  const cfnInputRef = React.useRef<HTMLInputElement>(null)
  const restoreRef = React.useRef<HTMLInputElement>(null)


  const [oldCfns, setOldCfns] = React.useState<string[] | null>(null);
  const [lastJSONExist, setLastJSONExist] = React.useState<boolean>(false);
  
  const [restore, setRestore] = React.useState(false)

  React.useEffect(() => {
    if (oldCfns) return
    GetAvailableLogs().then(logs => setOldCfns(logs))
    ResultsJSONExist().then(exists => setLastJSONExist(exists))
  }, [oldCfns])

  const { register, handleSubmit, control } = useForm<CfnFormValues>()

  const onSubmit = (values: CfnFormValues) => {
    const cfn = cfnInputRef.current.value
    if (!cfn || cfn == '') return

    send({
      type: 'submit',
      cfn: cfn,
      restore
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
        disabled={restore}
        placeholder={t("cfnName")!}
        autoCapitalize="off"
        autoComplete="off"
        autoCorrect="off"
        autoSave="off"
      />
      {oldCfns && (
        <div className="mt-3 flex flex-wrap gap-2 content-center items-center text-center pr-3">
          {oldCfns.map(cfn => {
            return (
              <button
                disabled={false}
                onClick={_ => { 
                  if (!(cfnInputRef.current && restoreRef.current)) return

                  cfnInputRef.current.value = cfn
                  restoreRef.current.checked = false
                  setRestore(false)
                }}
                className="whitespace-nowrap bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-base backdrop-blur rounded-2xl transition-all items-center border-transparent border-opacity-5 border-[1px] px-3 py-1"
                type="button"
                key={cfn}
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
            ref={restoreRef}
            type="checkbox"
            className="w-7 h-7 rounded-md checked:border-2 checked:focus:border-[rgba(255,255,255,.25)] checked:hover:border-[rgba(255,255,255,.25)] checked:border-[rgba(255,255,255,.25)] border-2 border-[rgba(255,255,255,.25)] focus:border-2 cursor-pointer bg-transparent text-transparent focus:ring-offset-transparent focus:ring-transparent mr-4"
            onChange={e => {
              if (e.target.checked) {
                setRestore(e.target.checked)
                cfnInputRef.current.value = ''
              }
            }}
          />
          {t("restoreSession")}
        </div>
      )}
      <div className="flex justify-end">
        <ActionButton
          type="submit"
          className="mt-4"
          style={{ filter: "hue-rotate(156deg)" }}
        >
          {t("start")}
        </ActionButton>
      </div>
    </form>
  )
}
