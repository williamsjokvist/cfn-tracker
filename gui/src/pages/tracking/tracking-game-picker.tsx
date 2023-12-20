import React from "react";
import { useTranslation } from "react-i18next";
import { motion } from "framer-motion";

import sfv from "./games/sfv.png";
import sf6 from "./games/sf6.png";
import { ActionButton } from "@/ui/action-button";
import { GameButton } from "@/ui/game-button";
import { PageHeader } from "@/ui/page-header";

const GAMES = [
  {
    logo: sfv,
    code: "sfv",
    alt: "Street Fighter V",
  },
  {
    logo: sf6,
    code: "sf6",
    alt: "Street Fighter 6",
  },
];

type TrackingGamePickerProps = {
  onSubmit: (game: string) => void
}

export const TrackingGamePicker: React.FC<TrackingGamePickerProps> = ({ onSubmit }) => {
  const { t } = useTranslation();
  const [selectedGame, setSelectedGame] = React.useState<string | undefined>();

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
            selectedGame && onSubmit(selectedGame)
          }}
          disabled={!selectedGame}
        >
          {t("continueStep")}
        </ActionButton>
      </div>
    </>
  );
};
