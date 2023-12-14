import React from "react";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";

import sfvLogo from "@/img/logo_sfv.png";
import sf6Logo from "@/img/logo_sf6.png";
import { ActionButton } from "@/ui/action-button";
import { TRACKING_MACHINE } from "@/main/machine";
import { GameButton } from "@/ui/game-button";
import { PageHeader } from "@/ui/page-header";
import { useMachine } from "@xstate/react";

const GAMES = [
  {
    logo: sfvLogo,
    code: "sfv",
    alt: "Street Fighter V",
  },
  {
    logo: sf6Logo,
    code: "sf6",
    alt: "Street Fighter 6",
  },
];

export const TrackingGamePicker: React.FC = () => {
  const { t } = useTranslation();
  const [selectedGame, setSelectedGame] = React.useState<string | undefined>();
  const [_, send] = useMachine(TRACKING_MACHINE);

  return (
    <>
      <PageHeader text={t("pickGame")}/>
      <div className="flex items-center flex-col gap-10 justify-self-center justify-center">
        <motion.ul
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.125 }}
          className="flex justify-center w-full gap-8"
        >
          {GAMES.map((game) => (
            <li key={game.code}>
              <GameButton
                {...(game.code == selectedGame && {
                  style: {
                    outline: "1px solid lightblue",
                    background: "rgb(248 250 252 / 0.05)",
                  },
                })}
                onClick={() => setSelectedGame(game.code)}
                {...game}
              />
            </li>
          ))}
        </motion.ul>
        <ActionButton
          onClick={() => {
            selectedGame &&
              send({
                type: "submit",
                game: selectedGame,
              });
          }}
          disabled={!selectedGame}
        >
          {t("continueStep")}
        </ActionButton>
      </div>
    </>
  );
};
