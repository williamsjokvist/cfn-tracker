import React from "react";
import { useTranslation } from "react-i18next";
import { useAnimate } from "framer-motion";

import sfvLogo from "@/img/logo_sfv.png";
import sf6Logo from "@/img/logo_sf6.png";
import { ActionButton } from "@/ui/action-button";
import { CFNMachineContext } from "@/main/machine";
import { GameButton } from "@/ui/game-button";

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
  const [_, send] = CFNMachineContext.useActor();

  const [scope, animate] = useAnimate();

  React.useEffect(() => {
    animate("li", { opacity: [0, 1] }, { delay: 0.25 });
  }, []);

  return (
    <div className="w-full flex items-center flex-col gap-10">
      <ul ref={scope} className="flex justify-center w-full gap-8">
        {GAMES.map((game) => {
          return (
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
          );
        })}
      </ul>
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
  );
};
