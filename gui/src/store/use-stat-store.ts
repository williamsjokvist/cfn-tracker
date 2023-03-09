import { create } from "zustand";
import { core } from "../../wailsjs/go/models";

const initialState: core.MatchHistory = {
  cfn: "",
  wins: 0,
  losses: 0,
  winRate: 0,
  lpGain: 0,
  lp: 0,
  opponent: "",
  opponentCharacter: "",
  opponentLeague: "",
  opponentLP: 0,
  totalLosses: 0,
  totalMatches: 0,
  totalWins: 0,
  result: false,
  winStreak: 0,
  timestamp: "",
  date: "",
};

type State = {
  matchHistory: core.MatchHistory | null;
  isTracking: boolean;
  isLoading: boolean;
  isPaused: boolean;
  isInitialized: boolean;
};

type Actions = {
  setMatchHistory: (mh: core.MatchHistory) => void;
  setTracking: (isTracking: boolean) => void;
  setLoading: (isLoading: boolean) => void;
  setPaused: (isPaused: boolean) => void;
  setInitialized: (isInitialized: boolean) => void;
  resetMatchHistory: () => void;
};

export const useStatStore = create<State & Actions>((set) => ({
  matchHistory: initialState,
  isTracking: false,
  isLoading: false,
  isInitialized: false,
  isPaused: false,
  setInitialized: (isInitialized) =>
    set((state) => ({ isInitialized: isInitialized })),
  setLoading: (isLoading) => set((state) => ({ isLoading: isLoading })),
  setMatchHistory: (mh) => set((state) => ({ matchHistory: mh })),
  setTracking: (isTracking) => set((state) => ({ isTracking: isTracking })),
  setPaused: (isPaused) => set((state) => ({ isPaused: isPaused })),
  resetMatchHistory: () => set((state) => ({ matchHistory: initialState })),
}));
