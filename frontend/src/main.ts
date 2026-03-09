import { createApp } from 'vue'
import TDesign from 'tdesign-vue-next'

import '@unocss/reset/tailwind.css'
import 'virtual:uno.css'
import 'tdesign-vue-next/es/style/index.css'

import App from '@/App.vue'
import router from '@/router'
import pinia, { useAppStore } from '@/store'
import '@/style.css'

const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(TDesign)

useAppStore(pinia).initializeTheme()

app.mount('#app')
