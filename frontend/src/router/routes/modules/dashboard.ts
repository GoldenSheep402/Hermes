import type { AppRouteRecordRaw } from '../types'
import { DEFAULT_LAYOUT } from '../base'

const DASHBOARD: AppRouteRecordRaw = {
  path: '/dashboard',
  name: 'dashboard',
  component: DEFAULT_LAYOUT,
  meta: {
    label: '仪表盘',
    requiresAuth: true,
    icon: 'icon-dashboard',
    order: 1,
  },
  children: [
    {
      path: 'workplace',
      name: 'Workplace',
      component: () => import('@/views/dashboard/workplace/index.vue'),
      meta: {
        label: '工作台',
        // locale: 'menu.dashboard.workplace',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'site-stats',
      name: 'SiteStats',
      component: () => import('@/views/dashboard/site-stats/index.vue'),
      meta: {
        label: '全站数据',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'tracker-stats',
      name: 'TrackerStats',
      component: () => import('@/views/dashboard/tracker-stats/index.vue'),
      meta: {
        label: 'Tracker负载',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
}

export default DASHBOARD
