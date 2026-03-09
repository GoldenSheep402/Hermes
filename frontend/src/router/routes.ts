import type { RouteRecordRaw } from 'vue-router'

const DefaultLayout = () => import('@/layouts/DefaultLayout.vue')

export const routes: RouteRecordRaw[] = [
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
  {
    path: '/',
    component: DefaultLayout,
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: '',
        redirect: '/home',
      },
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/pages/Home.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'torrents',
        name: 'Torrents',
        component: () => import('@/pages/Torrents.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'forums',
        name: 'Forums',
        component: () => import('@/pages/PlaceholderPage.vue'),
        props: {
          title: 'Forums',
          description: '论坛版块建设中，后续会接入版面、帖子列表和回帖流。',
        },
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'top10',
        name: 'Top10',
        component: () => import('@/pages/PlaceholderPage.vue'),
        props: {
          title: 'Top 10',
          description: '排行榜页面建设中，后续会展示人气种子、活跃用户和分享率榜单。',
        },
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'rules',
        name: 'Rules',
        component: () => import('@/pages/PlaceholderPage.vue'),
        props: {
          title: 'Rules',
          description: '站点规则页面建设中，后续将包含分享率、安全与账号规范。',
        },
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'faq',
        name: 'Faq',
        component: () => import('@/pages/PlaceholderPage.vue'),
        props: {
          title: 'FAQ',
          description: '常见问题页面建设中，后续会补充下载、做种与客户端配置指南。',
        },
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'admin',
        name: 'Admin',
        component: () => import('@/pages/Admin.vue'),
        meta: {
          requiresAuth: true,
          requiresStaff: true,
        },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/pages/NotFound.vue'),
    meta: {
      requiresAuth: false,
    },
  },
]
