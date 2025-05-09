import type { AppRouteRecordRaw } from '../types'
import { DEFAULT_LAYOUT } from '../base'

const MY: AppRouteRecordRaw = {
  path: '/my',
  name: 'My',
  component: DEFAULT_LAYOUT,
  meta: {
    label: '我的',
    requiresAuth: true,
    icon: 'icon-user',
    order: 3,
  },
  children: [
    {
      path: 'info',
      name: 'MyInfo',
      component: () => import('@/views/my/index.vue'),
      meta: {
        label: '个人信息',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'magic-points',
      name: 'MagicPoints',
      component: () => import('@/views/my/magic-points.vue'),
      meta: {
        label: '魔力值明细',
        requiresAuth: true,
        roles: ['*'],
        hideInMenu: true,
      },
    },
    {
      path: 'magic-shop',
      name: 'MagicShop',
      component: () => import('@/views/my/magic-shop.vue'),
      meta: {
        label: '魔力值商店',
        requiresAuth: true,
        roles: ['*'],
        hideInMenu: true,
      },
    },
  ],
}

export default MY
 