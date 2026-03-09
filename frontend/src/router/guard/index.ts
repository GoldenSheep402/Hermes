import type { Router } from 'vue-router'
import setupAuthGuard from './auth'

export default function createRouteGuard(router: Router) {
  setupAuthGuard(router)
}
