import { create } from 'zustand'
import { IMatchHistory } from '../types/match-history'


const initialState: IMatchHistory = {
    cfn: '',
    wins: 0,
    losses: 0,
    winRate: 0,
    lpGain: 0,
    lp: 0,
    opponent: '',
    opponentCharacter: '',
    opponentLP: '',
    totalLosses: 0,
    totalMatches: 0,
    totalWins: 0,
    result: false,
    winStreak: 0,
    timestamp: ''
}

type State = {
    matchHistory: IMatchHistory | null,
    isTracking: boolean,
    isLoading: boolean,
    isPaused: boolean,
    isInitialized: boolean
}

type Actions = {
    setMatchHistory: (mh: IMatchHistory) => void
    setTracking: (isTracking: boolean) => void 
    setLoading: (isLoading: boolean) => void 
    setPaused: (isPaused: boolean) => void 
    setInitialized: (isInitialized: boolean) => void 
    resetMatchHistory: () => void
}

export const useStatStore = create<State & Actions>(
    (set) => ({
        matchHistory: initialState,
        isTracking: false,
        isLoading: false,
        isInitialized: false,
        isPaused: false,
        setInitialized: (isInitialized) => set((state) => ({ isInitialized: isInitialized }) ),
        setLoading: (isLoading) => set((state) => ({ isLoading: isLoading }) ),
        setMatchHistory: (mh) => set((state) => ({ matchHistory: mh }) ),
        setTracking: (isTracking) => set((state) => ({ isTracking: isTracking }) ),
        setPaused: (isPaused) => set((state) => ({ isPaused: isPaused }) ),
        resetMatchHistory: () => set((state) => ({ matchHistory: initialState}))
    })
)