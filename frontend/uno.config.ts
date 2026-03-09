import { defineConfig, presetUno, presetAttributify, presetIcons } from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons({
      scale: 1.2,
      warn: true,
    }),
  ],
  theme: {
    colors: {
      primary: '#3b82f6', // Tailwind blue-500
      secondary: '#10b981', // Tailwind emerald-500
      accent: '#8b5cf6', // Tailwind violet-500
    }
  }
})
