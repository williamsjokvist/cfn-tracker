export const signed = (value: number) => (value > 0 ? `+${value}` : value)

export const winRate = (wins: number, losses: number) =>
  wins + losses === 0 ? 0 : Math.floor((wins / (wins + losses)) * 100)
