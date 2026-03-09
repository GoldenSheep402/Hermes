import { createRouter, createWebHistory } from 'vue-router'
import { routes } from './routes'
import createRouteGuard from './guard'

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

createRouteGuard(router)

export default router
