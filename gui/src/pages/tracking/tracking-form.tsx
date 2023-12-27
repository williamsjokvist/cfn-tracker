import React from "react";
import { useLoaderData } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";
import { Icon } from "@iconify/react";

import { TrackingMachineContext } from "@/machines/tracking-machine";
import { ActionButton } from "@/ui/action-button";
import { Checkbox } from "@/ui/checkbox";
import { model } from "@@/go/models";
import { PageHeader } from "@/ui/page-header";

export const TrackingForm: React.FC = () => {
  const { t } = useTranslation();
  const users = (useLoaderData() ?? []) as model.User[]
  const trackingActor = TrackingMachineContext.useActorRef()
  
  const playerIdInputRef = React.useRef<HTMLInputElement>(null);
  const restoreRef = React.useRef<HTMLInputElement>(null);
  const [playerIdInput, setPlayerIdInput] = React.useState<string>("");

  const onSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    if (playerIdInput == "") return;
    trackingActor.send({
      type: "submit",
      user: {
        displayName: users.find((old) => old.code == playerIdInput) ?? playerIdInput,
        code: playerIdInput,
      },
      restore: restoreRef.current && restoreRef.current.checked,
    });
  };

  const playerChipClicked = (player: model.User) => {
    if (playerIdInputRef.current) {
      playerIdInputRef.current.value = player.code;
      setPlayerIdInput(player.code);
    }
  };

  const clearInput = () => {
    if (playerIdInputRef.current) {
      playerIdInputRef.current.value = "";
      setPlayerIdInput("");
    }
  }

  return (
    <>
      <PageHeader text={t("startTracking")} />
      <motion.form
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.125 }}
        className=" relative h-full overflow-y-scroll overflow-x-visible w-full justify-self-center flex flex-col pt-12 gap-5 px-56 pb-4"
        onSubmit={onSubmit}
      >
        <h3 className="text-lg">{t("enterCfnName")}</h3>
        <div className='relative'>
          <input
            ref={playerIdInputRef}
            onChange={(e) => setPlayerIdInput(e.target.value)}
            className="bg-transparent border-b-2 border-0 focus:ring-offset-transparent focus:ring-transparent border-b-[rgba(255,255,255,0.275)] focus:border-white hover:border-white outline-none focus:outline-none hover:text-white transition-colors pt-4 pb-3 px-4 pr-12 block w-full text-lg text-gray-300"
            type="text"
            placeholder={t("cfnName")!}
            autoCapitalize="off"
            autoComplete="off"
            autoCorrect="off"
            autoSave="off"
          />
          {playerIdInput.length > 0 && (
            <button 
              type="button" 
              onClick={clearInput}
              aria-label="Clear" 
              className='absolute top-0 right-0 mt-4 mr-4 text-[#bfbcff] hover:text-white hover:bg-[rgba(255,255,255,.11)] transition-colors rounded-md'
            >
              <Icon icon="ci:close-big" className='w-6 h-6' />
            </button>
          )}
        </div>
        {users.length > 0 && (
          <div className="flex flex-wrap gap-2 content-center items-center text-center">
            {users.map((player) => (
              <button
                key={player.displayName}
                type="button"
                onClick={() => playerChipClicked(player)}
                className="whitespace-nowrap bg-[rgb(255,255,255,0.075)] hover:bg-[rgb(255,255,255,0.125)] text-base rounded-2xl transition-all items-center px-3 pt-1"
              >
                {player.displayName}
              </button>
            ))}
          </div>
        )}
        <footer className="flex items-center w-full">
          {users.some((old) => old.code == playerIdInput) && (
            <div className="group flex items-center">
              <Checkbox ref={restoreRef} id="restore-session" />
              <label
                htmlFor="restore-session"
                className="text-lg cursor-pointer text-gray-300 group-hover:text-white transition-colors"
              >
                {t("restoreSession")}
              </label>
            </div>
          )}
          <ActionButton
            type="submit"
            style={{ filter: "hue-rotate(-65deg)" }}
            className="ml-auto"
          >
            {t("start")}
          </ActionButton>
        </footer>
      </motion.form>
    </>
  );
};
