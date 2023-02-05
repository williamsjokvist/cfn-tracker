import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { FaChevronLeft } from "react-icons/fa";
import { MdOutlineDelete } from 'react-icons/md'
import { IMatchHistory } from '../types/match-history'
import {
  GetMatchLog,
  DeleteMatchLog
} from "../../wailsjs/go/main/App";

const History = () => {
  const { t } = useTranslation();
  const [matchLog, setMatchLog] = useState<IMatchHistory[]>();
  const [isSpecified, setSpecified] = useState(false)
  const fetchLog = async () => {
    const log = await GetMatchLog()
    setMatchLog([])
    setTimeout(() => {
      setMatchLog(log)
      setSpecified(false)
    }, 50)
  }

  useEffect(() => {
    !matchLog && fetchLog()
  }, [matchLog])

  return (
    <main className="grid grid-rows-[0fr_1fr] min-h-screen max-h-screen z-40 flex-1 text-white mx-auto">
      <header className='border-b border-slate-50 border-opacity-10 backdrop-blur select-none' style={{
        '--wails-draggable': 'drag'
      } as React.CSSProperties}>
        <h2 className="pt-4 px-8 flex items-center justify-between gap-5 uppercase text-sm tracking-widest mb-4">
          {t('history')}
        </h2>
      </header>
      <div className="relative w-full pt-2 z-40 pb-4">
        <div className='px-8 mb-2 h-10 border-b border-slate-50 border-opacity-10 '>
          {isSpecified && (
            <>
              <button onClick={() => {
                fetchLog()
              }} className='backdrop-blur leading-4 h-8 inline-block mr-3 rounded-xl transition-all items-center border-transparent hover:border-white border-opacity-5 border-[1px] px-3 py-1'>
                <FaChevronLeft className='w-4 h-4 inline mr-2' />
                {t('goBack')}
              </button>
            </>
          )}
          <button onClick={() => {
            DeleteMatchLog()
            setTimeout(() => {
              setMatchLog([])
            }, 50)
          }} className='backdrop-blur leading-4 h-8 inline-block float-right rounded-xl transition-all items-center border-transparent hover:border-white border-opacity-5 border-[1px] px-3 py-1'>
            <MdOutlineDelete className='w-4 h-4 inline mr-2' />
            {t('delete')}
          </button>
        </div>
        <div className='overflow-y-scroll max-h-[320px] h-full px-8'>
          <table className="w-full border-spacing-y-1 border-separate min-w-[525px]">
            <thead>
              <tr>
                <th className='text-left px-3 whitespace-nowrap'>{t('time')}</th>
                <th className='text-left px-3 whitespace-nowrap'>CFN</th>
                <th className='text-left px-3 whitespace-nowrap'>{t('opponent')}</th>
                <th className='text-left px-3 whitespace-nowrap'>{t('character')}</th>
                <th className='text-left px-3 whitespace-nowrap'>{t('result')}</th>
                <th className='text-left px-3 whitespace-nowrap'>{t('lpGain')}</th>
              </tr>
            </thead>
            <tbody>
              {matchLog && matchLog.map((log, index) => {
                return (
                  <tr key={index} className="w-full backdrop-blur">
                    <td className='whitespace-nowrap text-center rounded-l-xl rounded-r-none bg-slate-50 bg-opacity-5 px-3 py-2'>
                      {log.timestamp}
                    </td>
                    <td
                      onClick={() => {
                        setMatchLog([])
                        setTimeout(() => {
                          setMatchLog(matchLog.filter(ml => ml.cfn == log.cfn))
                          setSpecified(true)
                        }, 50)
                      }}
                      className='whitespace-nowrap rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 hover:underline cursor-pointer'>
                      {log.cfn}
                    </td>
                    <td
                      onClick={() => {
                        setMatchLog([])
                        setTimeout(() => {
                          setMatchLog(matchLog.filter(ml => ml.opponent == log.opponent))
                          setSpecified(true)
                        }, 50)
                      }}
                      className='whitespace-nowrap w-full rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 hover:underline cursor-pointer'>
                      {log.opponent}
                    </td>
                    <td
                      onClick={() => {
                        setMatchLog([])
                        setTimeout(() => {
                          setMatchLog(matchLog.filter(ml => ml.opponentCharacter == log.opponentCharacter))
                          setSpecified(true)
                        }, 50)
                      }}
                      className='rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 text-center hover:underline cursor-pointer'>
                      {log.opponentCharacter}
                    </td>
                    <td className='rounded-none bg-slate-50 bg-opacity-5 px-3 py-2 text-center' style={{
                      color: log.result == true ? 'lime' : 'red'
                    }}>
                      {log.result == true ? 'W' : 'L'}
                    </td>
                    <td className='rounded-r-xl rounded-l-none bg-slate-50 bg-opacity-5 px-3 py-2 text-center'>
                      {(log.lpGain > 0) && '+'}
                      {log.lpGain}
                    </td>
                  </tr>
                )
              })
              }
            </tbody>
          </table>
        </div>
      </div>
    </main>
  );
};

export default History;
