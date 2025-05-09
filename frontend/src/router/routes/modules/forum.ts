import type { AppRouteRecordRaw } from '../types'
import { DEFAULT_LAYOUT } from '../base'

const FORUM: AppRouteRecordRaw = {
  path: '/forum',
  name: 'Forum',
  component: DEFAULT_LAYOUT,
  meta: {
    label: '论坛',
    requiresAuth: true,
    icon: 'icon-message',
    order: 5,
  },
  children: [
    {
      path: 'list',
      name: 'ForumList',
      component: () => import('@/views/forum/index.vue'),
      meta: {
        label: '论坛首页',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'category/:id',
      name: 'ForumCategory',
      component: () => import('@/views/forum/category/index.vue'),
      meta: {
        label: '版块详情',
        requiresAuth: true,
        roles: ['*'],
        hideInMenu: true,
      },
    },
    {
      path: 'post/:id',
      name: 'ForumPost',
      component: () => import('@/views/forum/post/index.vue'),
      meta: {
        label: '帖子详情',
        requiresAuth: true,
        roles: ['*'],
        hideInMenu: true,
      },
    },
    {
      path: 'create',
      name: 'ForumCreate',
      component: () => import('@/views/forum/create/index.vue'),
      meta: {
        label: '发布帖子',
        requiresAuth: true,
        roles: ['*'],
        hideInMenu: true,
      },
    },
  ],
}

export default FORUM
