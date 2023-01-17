import { create } from 'zustand'

interface MatchHistory {
    cfn: string;
    wins: number;
    losses: number;
    winRate: any;
    lpGain: any;
    lp: number;
}

const initialState: MatchHistory = {
    cfn: '',
    wins: 0,
    losses: 0,
    winRate: 0,
    lpGain: 0,
    lp: 0
}

type State = {
    matchHistory: MatchHistory | null
}

type Actions = {
    setMatchHistory: (mh: MatchHistory) => void
    resetMatchHistory: () => void
}

export const useStatStore = create<State & Actions>(
    (set) => ({
        matchHistory: {
            cfn: '',
            wins: 0,
            losses: 0,
            winRate: 0,
            lpGain: 0,
            lp: 0
        },
        setMatchHistory: (mh) => set((state) => ({ matchHistory: mh }) ),
        resetMatchHistory: () => set((state) => ({ matchHistory: initialState}))
    })
)