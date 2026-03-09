import type { PermissionKey } from '@/constants/permissions'

export interface NavItem {
  path: string
  label: string
  requiredPermission?: PermissionKey
}
