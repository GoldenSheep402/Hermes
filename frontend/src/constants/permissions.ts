export const PermissionKeys = {
  AdminPanelAccess: 'admin.panel.access',
  SiteSettingsManage: 'site.settings.manage',
  InviteManage: 'invite.manage',
  TrackerManage: 'tracker.manage',
  TorrentModerate: 'torrent.moderate',
  CategoryManage: 'torrent.category.manage',
  UserManage: 'user.manage',
} as const

export type PermissionKey = (typeof PermissionKeys)[keyof typeof PermissionKeys]

const RolePermissionMap: Record<string, PermissionKey[]> = {
  user: [],
  staff: [
    PermissionKeys.AdminPanelAccess,
    PermissionKeys.SiteSettingsManage,
    PermissionKeys.InviteManage,
    PermissionKeys.TrackerManage,
    PermissionKeys.TorrentModerate,
    PermissionKeys.CategoryManage,
    PermissionKeys.UserManage,
  ],
}

export function resolvePermissionsByRoles(roles: string[]): PermissionKey[] {
  const permissionSet = new Set<PermissionKey>()

  for (const role of roles) {
    const normalizedRole = role.trim().toLowerCase()
    const rolePermissions = RolePermissionMap[normalizedRole]
    if (!rolePermissions) {
      continue
    }
    for (const permission of rolePermissions) {
      permissionSet.add(permission)
    }
  }

  return Array.from(permissionSet)
}
