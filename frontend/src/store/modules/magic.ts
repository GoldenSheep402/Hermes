import { defineStore } from 'pinia'

export const useMagicStore = defineStore('magic', {
  state: () => ({
    magicPoints: Number(localStorage.getItem('magicPoints')) || 0,
    currentMagicPoints: Number(localStorage.getItem('currentMagicPoints')) || 0,
  }),
  actions: {
    setMagicPoints(points: number) {
      this.magicPoints = points
      this.currentMagicPoints = points
      localStorage.setItem('magicPoints', points.toString())
      localStorage.setItem('currentMagicPoints', points.toString())
    },
    addMagicPoints(points: number) {
      this.magicPoints += points
      this.currentMagicPoints += points
      localStorage.setItem('magicPoints', this.magicPoints.toString())
      localStorage.setItem('currentMagicPoints', this.currentMagicPoints.toString())
    },
    subtractMagicPoints(points: number) {
      const newPoints = this.currentMagicPoints - points
      this.magicPoints = newPoints
      this.currentMagicPoints = newPoints
      localStorage.setItem('magicPoints', newPoints.toString())
      localStorage.setItem('currentMagicPoints', newPoints.toString())
    },
  },
}) 