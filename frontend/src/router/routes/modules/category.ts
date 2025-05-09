import type { AppRouteRecordRaw } from '../types'
import { DEFAULT_LAYOUT } from '../base'

const CATEGORY: AppRouteRecordRaw = {
  path: '/category',
  name: 'Category',
  component: DEFAULT_LAYOUT,
  meta: {
    label: '类别',
    requiresAuth: true,
    icon: 'icon-apps',
    order: 2,
  },
  children: [
    {
      path: 'list',
      name: 'CategoryList',
      component: () => import('@/views/category/index.vue'),
      meta: {
        label: '列表',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'detail/:id',
      name: 'CategoryDetail',
      component: () => import('@/views/category/detail.vue'),
      meta: {
        label: '类别详情',
        requiresAuth: true,
        roles: ['*'],
        hideInMenu: true,
      },
    },
  ],
}

export default CATEGORY
