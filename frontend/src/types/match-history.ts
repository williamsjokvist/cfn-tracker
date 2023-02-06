export interface IMatchHistory {
  cfn: string,
  losses: number,
  lp: number,
  lpGain: number,
  opponent: string,
  opponentCharacter: string,
  opponentLP: string,
  totalLosses: number,
  totalMatches: number,
  totalWins: number,
  winRate: number,
  wins: number,
  result: boolean,
  winStreak: number,
  timestamp: string
}