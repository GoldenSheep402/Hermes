import axios, { AxiosError } from 'axios'
import type { InternalAxiosRequestConfig } from 'axios'
import { MessagePlugin } from 'tdesign-vue-next'
import router from '@/router'
import pinia, { useAuthStore } from '@/store'

const http = axios.create({
  baseURL: '/gapi',
  timeout: 15000,
})

let unauthorizedHandled = false

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const authStore = useAuthStore(pinia)
  if (authStore.accessToken) {
    config.headers.Authorization = `Bearer ${authStore.accessToken}`
  }
  return config
})

http.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<{ message?: string }>) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore(pinia)
      authStore.handleUnauthorized()

      if (!unauthorizedHandled) {
        unauthorizedHandled = true
        await MessagePlugin.warning('登录状态已失效，请重新登录')
        const currentPath = router.currentRoute.value.fullPath
        await router.replace({
          name: 'Login',
          query: {
            redirect: currentPath,
          },
        })
        setTimeout(() => {
          unauthorizedHandled = false
        }, 300)
      }
    }

    const message = error.response?.data?.message || error.message || '请求失败'
    return Promise.reject(new Error(message))
  },
)

export default http
