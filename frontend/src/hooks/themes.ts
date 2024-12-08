import { useAppStore } from '@/store'
import { computed } from 'vue'

export default function useThemes() {
  const appStore = useAppStore()
  const isDark = computed(() => {
    return appStore.theme === 'dark'
  })
  return {
    isDark,
  }
}
