import { defineStore } from 'pinia'

function applyThemeMode(isDarkMode: boolean) {
  if (typeof document === 'undefined') {
    return
  }
  document.documentElement.setAttribute('theme-mode', isDarkMode ? 'dark' : 'light')
}

export const useAppStore = defineStore('hermes/app', {
  state: () => ({
    isDarkMode: false,
  }),

  actions: {
    initializeTheme() {
      applyThemeMode(this.isDarkMode)
    },

    setDarkMode(enabled: boolean) {
      this.isDarkMode = enabled
      applyThemeMode(enabled)
    },

    toggleDarkMode() {
      this.setDarkMode(!this.isDarkMode)
    },
  },

  persist: {
    key: 'hermes-app',
    pick: ['isDarkMode'],
  },
})
