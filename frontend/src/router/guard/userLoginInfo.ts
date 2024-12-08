import type { LocationQueryRaw, Router } from 'vue-router'
import { useUserStore } from '@/store'

import NProgress from 'nprogress' // progress bar

export default function setupUserLoginInfoGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    NProgress.start()
    const userStore = useUserStore()
    if (userStore.isLogin) {
      if (userStore.role) {
        next()
      }
      else {
        try {
          await userStore.info()
          next()
        }
        catch {
          await userStore.logout()
          next({
            name: 'login',
            query: {
              redirect: to.name,
              ...to.query,
            } as LocationQueryRaw,
          })
        }
      }
    }
    else {
      if (to.name === 'login') {
        next()
        return
      }
      if (to.name === 'register') {
        next()
        return
      }
      next({
        name: 'login',
        query: {
          redirect: to.name,
          ...to.query,
        } as LocationQueryRaw,
      })
    }
  })
}
