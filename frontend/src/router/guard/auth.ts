import type { Router } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { useAuthStore } from '@/store'
import { LOGIN_ROUTE_NAME, WHITE_LIST } from '@/router/constants'

export default function setupAuthGuard(router: Router) {
  router.beforeEach((to) => {
    const authStore = useAuthStore()
    const routeName = typeof to.name === 'string' ? to.name : ''
    const isWhiteRoute = WHITE_LIST.has(routeName)
    const requiresAuth = to.meta.requiresAuth !== false && !isWhiteRoute

    if (to.name === LOGIN_ROUTE_NAME && authStore.isLoggedIn) {
      return { path: '/home' }
    }

    if (requiresAuth && !authStore.isLoggedIn) {
      return {
        name: LOGIN_ROUTE_NAME,
        query: {
          redirect: to.fullPath,
        },
      }
    }

    if (to.meta.requiresStaff && !authStore.isStaff) {
      MessagePlugin.warning('当前账号无管理后台权限')
      return { path: '/403' }
    }

    return true
  })
}
