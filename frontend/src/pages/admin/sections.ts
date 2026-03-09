import type { PermissionKey } from '@/constants/permissions'
import { PermissionKeys } from '@/constants/permissions'

export interface AdminSectionItem {
  key: 'site' | 'invite' | 'tracker' | 'torrent' | 'category' | 'metaTemplate' | 'user'
  label: string
  description: string
  routeName: string
  requiredPermission: PermissionKey
}

export const ADMIN_SECTIONS: AdminSectionItem[] = [
  {
    key: 'site',
    label: '站点基础设置',
    description: '维护站点名称、公告和维护模式',
    routeName: 'AdminSiteSettings',
    requiredPermission: PermissionKeys.SiteSettingsManage,
  },
  {
    key: 'invite',
    label: '注册与邀请控制',
    description: '管理开放注册、邀请码和邮箱验证策略',
    routeName: 'AdminInviteSettings',
    requiredPermission: PermissionKeys.InviteManage,
  },
  {
    key: 'tracker',
    label: 'Tracker 参数',
    description: '调整 announce、flush、freeleech 等参数',
    routeName: 'AdminTrackerSettings',
    requiredPermission: PermissionKeys.TrackerManage,
  },
  {
    key: 'torrent',
    label: '种子管理与审核',
    description: '执行死种清理和举报处理等运营动作',
    routeName: 'AdminTorrentModeration',
    requiredPermission: PermissionKeys.TorrentModerate,
  },
  {
    key: 'category',
    label: '种子类别管理',
    description: '维护类别层级和元数据模板',
    routeName: 'AdminCategoryManage',
    requiredPermission: PermissionKeys.CategoryManage,
  },
  {
    key: 'metaTemplate',
    label: '元数据模板',
    description: '维护各类别的元数据字段',
    routeName: 'AdminMetaTemplateManage',
    requiredPermission: PermissionKeys.CategoryManage,
  },
  {
    key: 'user',
    label: '用户管理',
    description: '封禁账号、重置流量、等级提升',
    routeName: 'AdminUserManage',
    requiredPermission: PermissionKeys.UserManage,
  },
]
