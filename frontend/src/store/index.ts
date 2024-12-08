import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import useAppStore from './modules/app'
import useTabBarStore from './modules/tab-bar'
import useUserStore from './modules/user'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

export { useAppStore, useTabBarStore, useUserStore }
export default pinia
