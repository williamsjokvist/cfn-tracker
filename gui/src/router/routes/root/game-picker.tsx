import React from 'react';
import { useTranslation } from 'react-i18next';
import { useAnimate, stagger } from "framer-motion";

import sfvLogo from '@/img/logo_sfv.png'
import sf6Logo from '@/img/logo_sf6.png'
import { ActionButton } from '@/ui/action-button';
import { CFNMachineContext } from '@/machine';

type GameButtonProps = {
  logo: string;
  code: string;
  alt: string;
} & React.HTMLAttributes<HTMLButtonElement>
const GameButton: React.FC<GameButtonProps> = (props) => {
  return (
    <button 
      type='button' 
      {...props}
      className='w-52 p-3 rounded-lg hover:bg-slate-50 hover:bg-opacity-5 transition-colors'
    >
      <img src={props.logo} alt={props.alt} className='pointer-events-none select-none'/>
    </button>
  )
}

const GAMES = [{
  logo: sfvLogo,
  code: 'sfv',
  alt: 'Street Fighter V'
},
{
  logo: sf6Logo,
  code: 'sf6',
  alt: 'Street Fighter 6'
}]

export const GamePicker: React.FC = () => {
  const { t } = useTranslation()
  const [selectedGame, setSelectedGame] = React.useState('')
  const [_, send] = CFNMachineContext.useActor();

  const [scope, animate] = useAnimate()

  React.useEffect(() => {
    animate("li", 
      { opacity: [0, 1] }, 
      { delay: stagger(0.125, { ease: "linear" }) }
    )
  }, [])

  return (
    <div className='w-full flex items-center flex-col gap-10'>
      <ul ref={scope} className='flex justify-center w-full gap-8'>
        {GAMES.map((game) => {
          return (
            <li key={game.code}>
              <GameButton 
                {...(game.code == selectedGame && {
                  style: {
                    outline: '1px solid lightblue',
                    background: 'rgb(248 250 252 / 0.05)'
                  }
                })}
                onClick={() => setSelectedGame(game.code)}
                {...game}
              />
            </li>
          )
        })}
      </ul>
      <ActionButton
        onClick={() => {
          if (selectedGame == '') return
          send({
            type: 'submit',
            game: selectedGame
          })
        }}
        style={{
          ...((selectedGame == '') && {
            filter: 'saturate(0)',
          })
        }}
        >
        {t('continueStep')}
      </ActionButton>
    </div>
  )
}