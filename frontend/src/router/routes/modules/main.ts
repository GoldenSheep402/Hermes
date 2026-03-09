import type { RouteRecordRaw } from 'vue-router'
import { PermissionKeys } from '@/constants/permissions'
import { DEFAULT_LAYOUT } from '../base'

const mainRoute: RouteRecordRaw = {
  path: '/',
  component: DEFAULT_LAYOUT,
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
      path: 'torrents/:id',
      name: 'TorrentDetail',
      component: () => import('@/pages/TorrentDetail.vue'),
      meta: {
        requiresAuth: true,
      },
    },
    {
      path: 'upload',
      name: 'UploadTorrent',
      component: () => import('@/pages/UploadTorrent.vue'),
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
      component: () => import('@/pages/admin/AdminLayout.vue'),
      meta: {
        requiresAuth: true,
        requiredPermissions: [PermissionKeys.AdminPanelAccess],
      },
      children: [
        {
          path: '',
          name: 'AdminOverview',
          component: () => import('@/pages/admin/AdminOverview.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess],
          },
        },
        {
          path: 'site',
          name: 'AdminSiteSettings',
          component: () => import('@/pages/admin/AdminSiteSettings.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess, PermissionKeys.SiteSettingsManage],
          },
        },
        {
          path: 'invite',
          name: 'AdminInviteSettings',
          component: () => import('@/pages/admin/AdminInviteSettings.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess, PermissionKeys.InviteManage],
          },
        },
        {
          path: 'tracker',
          name: 'AdminTrackerSettings',
          component: () => import('@/pages/admin/AdminTrackerSettings.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess, PermissionKeys.TrackerManage],
          },
        },
        {
          path: 'torrent',
          name: 'AdminTorrentModeration',
          component: () => import('@/pages/admin/AdminTorrentModeration.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess, PermissionKeys.TorrentModerate],
          },
        },
        {
          path: 'category',
          name: 'AdminCategoryManage',
          component: () => import('@/pages/admin/AdminCategoryManage.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess, PermissionKeys.CategoryManage],
          },
        },
        {
          path: 'category/:id',
          name: 'AdminCategoryDetail',
          component: () => import('@/pages/admin/AdminCategoryDetail.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess, PermissionKeys.CategoryManage],
          },
        },
        {
          path: 'meta-template',
          name: 'AdminMetaTemplateManage',
          component: () => import('@/pages/admin/AdminMetaTemplateManage.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess, PermissionKeys.CategoryManage],
          },
        },
        {
          path: 'user',
          name: 'AdminUserManage',
          component: () => import('@/pages/admin/AdminUserManage.vue'),
          meta: {
            requiresAuth: true,
            requiredPermissions: [PermissionKeys.AdminPanelAccess, PermissionKeys.UserManage],
          },
        },
      ],
    },
  ],
}

export default mainRoute
