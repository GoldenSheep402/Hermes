import type { CategoryMetaTemplate } from '@/lib/proto/category/v1/category.pb'

export interface AdminCategoryMetaTemplateRow {
  id: string
  categoryId: string
  key: string
  label: string
  type: string
  required: boolean
  options: string
  sortOrder: number
  defaultValue: string
}

export interface AdminCategoryRow {
  id: string
  name: string
  slug: string
  description: string
  parentId: string
  icon: string
  sortOrder: number
  isEnabled: boolean
  depth: number
  metaTemplates: AdminCategoryMetaTemplateRow[]
}

export interface CategoryPresetTemplate {
  key: string
  label: string
  type: string
  required?: boolean
  options?: string
  sortOrder?: number
  defaultValue?: string
}

export interface CategoryPreset {
  value: string
  label: string
  name: string
  slug: string
  icon: string
  description: string
  templates: CategoryPresetTemplate[]
}

export interface UploadCategoryOption {
  id: string
  label: string
  templates: CategoryMetaTemplate[]
}

export interface UploadFormState {
  title: string
  subtitle: string
  description: string
  categoryId: string
  tags: string
  screenshots: string
}

export interface UploadCustomMetaEntry {
  id: string
  key: string
  value: string
}
