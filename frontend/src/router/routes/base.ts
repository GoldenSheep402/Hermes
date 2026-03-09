import type { RouteRecordRaw } from 'vue-router'

export const DEFAULT_LAYOUT = () => import('@/layouts/DefaultLayout.vue')

export const publicRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/Login.vue'),
    meta: {
      requiresAuth: false,
    },
  },
  {
    path: '/403',
    name: 'Forbidden',
    component: () => import('@/pages/Forbidden.vue'),
    meta: {
      requiresAuth: false,
    },
  },
]

export const NOT_FOUND_ROUTE: RouteRecordRaw = {
  path: '/:pathMatch(.*)*',
  name: 'NotFound',
  component: () => import('@/pages/NotFound.vue'),
  meta: {
    requiresAuth: false,
  },
}
