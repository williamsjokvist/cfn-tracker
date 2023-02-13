import { create } from 'zustand'

type State = {
  newVersionAvailable: boolean
}

type Actions = {
  setNewVersionAvailable: (newVersionAvailable: boolean) => void
}

export const useAppStore = create<State & Actions>(
  (set) => ({
      newVersionAvailable: false,
      setNewVersionAvailable: (newVersionAvailable: boolean) => set({ newVersionAvailable }),
  })
)