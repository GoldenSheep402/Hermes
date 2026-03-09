<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Category, CategoryMetaTemplate, MetaTemplatePreset } from '@/lib/proto/category/v1/category.pb'
import { CategoryService } from '@/services/grpc'

interface CategoryRow {
  id: string
  name: string
  slug: string
  description: string
  parentId: string
  icon: string
  sortOrder: number
  isEnabled: boolean
  depth: number
  metaTemplates: CategoryMetaTemplateRow[]
}

interface CategoryMetaTemplateRow {
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

interface CategoryPresetTemplate {
  key: string
  label: string
  type: string
  required?: boolean
  options?: string
  sortOrder?: number
  defaultValue?: string
}

interface CategoryPreset {
  value: string
  label: string
  name: string
  slug: string
  icon: string
  description: string
  templates: CategoryPresetTemplate[]
}

const categoryPresetList: CategoryPreset[] = [
  {
    value: 'movie',
    label: '影视资源',
    name: '电影',
    slug: 'movies',
    icon: '🎬',
    description: '电影、剧集、纪录片等视频资源',
    templates: [
      { key: 'language', label: '语言', type: 'text' },
      { key: 'region', label: '地区', type: 'text' },
      { key: 'year', label: '年份', type: 'number' },
      { key: 'source', label: '来源', type: 'text' },
      { key: 'resolution', label: '分辨率', type: 'select' },
      { key: 'video_codec', label: '视频编码', type: 'select' },
      { key: 'audio_codec', label: '音频编码', type: 'text' },
      { key: 'subtitle', label: '字幕', type: 'text' },
      { key: 'imdb_id', label: 'IMDb ID', type: 'text' },
      { key: 'douban_id', label: '豆瓣 ID', type: 'text' },
    ],
  },
  {
    value: 'music',
    label: '音乐资源',
    name: '音乐',
    slug: 'music',
    icon: '🎵',
    description: '专辑、单曲、演唱会音频等',
    templates: [
      { key: 'language', label: '语言', type: 'text' },
      { key: 'region', label: '地区', type: 'text' },
      { key: 'year', label: '年份', type: 'number' },
      { key: 'source', label: '来源', type: 'text' },
      { key: 'artist', label: '艺术家', type: 'text', required: true },
      { key: 'album', label: '专辑', type: 'text' },
      { key: 'format', label: '格式', type: 'select' },
      { key: 'bitrate', label: '码率', type: 'select' },
    ],
  },
  {
    value: 'book',
    label: '图书资源',
    name: '图书',
    slug: 'books',
    icon: '📚',
    description: '电子书、教程、杂志等文档资源',
    templates: [
      { key: 'language', label: '语言', type: 'text' },
      { key: 'region', label: '地区', type: 'text' },
      { key: 'year', label: '年份', type: 'number' },
      { key: 'source', label: '来源', type: 'text' },
      { key: 'author', label: '作者', type: 'text', required: true },
      { key: 'publisher', label: '出版社', type: 'text' },
      { key: 'isbn', label: 'ISBN', type: 'text' },
      { key: 'book_format', label: '书籍格式', type: 'select' },
    ],
  },
  {
    value: 'game',
    label: '游戏资源',
    name: '游戏',
    slug: 'games',
    icon: '🎮',
    description: 'PC/主机/移动平台游戏资源',
    templates: [
      { key: 'language', label: '语言', type: 'text' },
      { key: 'region', label: '地区', type: 'text' },
      { key: 'year', label: '年份', type: 'number' },
      { key: 'source', label: '来源', type: 'text' },
      { key: 'platform', label: '平台', type: 'select' },
      { key: 'game_version', label: '版本', type: 'text' },
      { key: 'multiplayer', label: '联机模式', type: 'select' },
      { key: 'crack_status', label: '破解状态', type: 'select' },
    ],
  },
  {
    value: 'software',
    label: '软件资源',
    name: '软件',
    slug: 'software',
    icon: '🧰',
    description: '应用程序、工具、插件等',
    templates: [
      { key: 'language', label: '语言', type: 'text' },
      { key: 'region', label: '地区', type: 'text' },
      { key: 'year', label: '年份', type: 'number' },
      { key: 'source', label: '来源', type: 'text' },
      { key: 'os', label: '操作系统', type: 'select' },
      { key: 'version', label: '版本', type: 'text', required: true },
      { key: 'license', label: '授权方式', type: 'select' },
      { key: 'official_website', label: '官网', type: 'url' },
    ],
  },
  {
    value: 'course',
    label: '课程资源',
    name: '课程',
    slug: 'courses',
    icon: '🧠',
    description: '教学视频、培训资料、学习内容',
    templates: [
      { key: 'language', label: '语言', type: 'text' },
      { key: 'region', label: '地区', type: 'text' },
      { key: 'year', label: '年份', type: 'number' },
      { key: 'source', label: '来源', type: 'text' },
      { key: 'instructor', label: '讲师', type: 'text' },
      { key: 'level', label: '难度', type: 'select' },
      { key: 'duration', label: '时长', type: 'text' },
    ],
  },
  {
    value: 'generic',
    label: '通用资源',
    name: '其他',
    slug: 'misc',
    icon: '📦',
    description: '不确定类型时使用通用预设',
    templates: [
      { key: 'language', label: '语言', type: 'text' },
      { key: 'region', label: '地区', type: 'text' },
      { key: 'year', label: '年份', type: 'number' },
      { key: 'source', label: '来源', type: 'text' },
      { key: 'author', label: '作者/发布者', type: 'text' },
      { key: 'version', label: '版本', type: 'text' },
      { key: 'format', label: '格式', type: 'text' },
    ],
  },
]

const categoryRows = ref<CategoryRow[]>([])
const categoryLoading = ref(false)
const categorySaving = ref(false)
const categoryEditorVisible = ref(false)
const categoryEditorMode = ref<'create' | 'edit'>('create')
const editingCategoryId = ref('')
const selectedCategoryPreset = ref('')
const selectedTemplateCategoryId = ref('')
const metaTemplatePresets = ref<MetaTemplatePreset[]>([])
const selectedMetaTemplatePreset = ref('')
const metaTemplateLoading = ref(false)
const metaTemplateSaving = ref(false)
const metaTemplatePresetApplying = ref(false)
const metaTemplateEditorVisible = ref(false)
const metaTemplateEditorMode = ref<'create' | 'edit'>('create')
const editingMetaTemplateId = ref('')

const categoryForm = reactive({
  name: '',
  slug: '',
  description: '',
  parentId: '',
  icon: '',
  sortOrder: 0,
  isEnabled: true,
})

const metaTemplateForm = reactive({
  categoryId: '',
  key: '',
  label: '',
  type: 'text',
  required: false,
  options: '',
  sortOrder: 0,
  defaultValue: '',
})

const categoryColumns = [
  { colKey: 'name', title: '类别名称', minWidth: 220 },
  { colKey: 'slug', title: 'Slug', width: 180 },
  { colKey: 'sortOrder', title: '排序', width: 90, align: 'right' },
  { colKey: 'isEnabled', title: '状态', width: 90 },
  { colKey: 'actions', title: '操作', width: 160 },
]

const categoryMetaTemplateColumns = [
  { colKey: 'key', title: 'Key', width: 180 },
  { colKey: 'label', title: '字段名', minWidth: 180 },
  { colKey: 'type', title: '类型', width: 110 },
  { colKey: 'required', title: '必填', width: 80 },
  { colKey: 'sortOrder', title: '排序', width: 80, align: 'right' },
  { colKey: 'defaultValue', title: '默认值', minWidth: 160 },
  { colKey: 'actions', title: '操作', width: 140 },
]

const metaTemplateTypeOptions = [
  { label: 'text', value: 'text' },
  { label: 'textarea', value: 'textarea' },
  { label: 'select', value: 'select' },
  { label: 'number', value: 'number' },
  { label: 'date', value: 'date' },
  { label: 'url', value: 'url' },
]

const parentCategoryOptions = computed(() =>
  categoryRows.value
    .filter((item) => item.id !== editingCategoryId.value)
    .map((item) => ({
      label: `${'  '.repeat(item.depth)}${item.name}`,
      value: item.id,
    })),
)

const categoryPresetOptions = computed(() =>
  categoryPresetList.map((item) => ({
    label: item.label,
    value: item.value,
  })),
)

const selectedCategoryPresetItem = computed(
  () => categoryPresetList.find((item) => item.value === selectedCategoryPreset.value) || null,
)

const metaTemplatePresetOptions = computed(() =>
  metaTemplatePresets.value
    .map((item) => ({
      label: item.label || item.value || '未命名预设',
      value: item.value || '',
    }))
    .filter((item) => Boolean(item.value)),
)

const selectedMetaTemplatePresetItem = computed(
  () => metaTemplatePresets.value.find((item) => item.value === selectedMetaTemplatePreset.value) || null,
)

const selectedTemplateCategory = computed(
  () => categoryRows.value.find((item) => item.id === selectedTemplateCategoryId.value) || null,
)

const selectedCategoryMetaTemplates = computed(() => {
  const list = selectedTemplateCategory.value?.metaTemplates || []
  return [...list].sort((left, right) => {
    if (left.sortOrder === right.sortOrder) {
      return left.key.localeCompare(right.key)
    }
    return left.sortOrder - right.sortOrder
  })
})

const filteredCategoryMetaTemplates = computed(() => {
  const preset = selectedMetaTemplatePresetItem.value
  if (!preset) {
    return selectedCategoryMetaTemplates.value
  }

  const templates = preset.templates || []
  const keySet = new Set(
    templates
      .map((item) => String(item.key || '').trim().toLowerCase())
      .filter(Boolean),
  )
  return selectedCategoryMetaTemplates.value.filter((item) => keySet.has(item.key.trim().toLowerCase()))
})

watch(selectedTemplateCategoryId, (categoryID) => {
  if (!metaTemplateEditorVisible.value || metaTemplateEditorMode.value !== 'create') {
    return
  }
  metaTemplateForm.categoryId = categoryID
})

function normalizeMetaTemplate(item: CategoryMetaTemplate | undefined): CategoryMetaTemplateRow | null {
  if (!item) {
    return null
  }
  const id = String(item.id || '').trim()
  const key = String(item.key || '').trim()
  if (!id || !key) {
    return null
  }

  return {
    id,
    categoryId: String(item.categoryId || '').trim(),
    key,
    label: String(item.label || key).trim(),
    type: String(item.type || 'text').trim().toLowerCase() || 'text',
    required: item.required === true,
    options: String(item.options || ''),
    sortOrder: Number(item.sortOrder || 0),
    defaultValue: String(item.defaultValue || ''),
  }
}

function flattenCategories(items: Category[] | undefined, depth = 0): CategoryRow[] {
  const result: CategoryRow[] = []
  for (const item of items || []) {
    const templates = (item.metaTemplates || [])
      .map((tpl) => normalizeMetaTemplate(tpl))
      .filter((tpl): tpl is CategoryMetaTemplateRow => Boolean(tpl))

    const row: CategoryRow = {
      id: item.id || '',
      name: item.name || '未命名',
      slug: item.slug || '',
      description: item.description || '',
      parentId: item.parentId || '',
      icon: item.icon || '',
      sortOrder: item.sortOrder || 0,
      isEnabled: item.isEnabled !== false,
      depth,
      metaTemplates: templates,
    }
    result.push(row)

    if (item.children && item.children.length > 0) {
      result.push(...flattenCategories(item.children, depth + 1))
    }
  }
  return result
}

async function loadCategories() {
  categoryLoading.value = true
  try {
    const response = await CategoryService.ListCategories({})
    const flattened = flattenCategories(response.categories)
    categoryRows.value = flattened.sort((left, right) => left.sortOrder - right.sortOrder)

    if (!selectedTemplateCategoryId.value) {
      selectedTemplateCategoryId.value = categoryRows.value[0]?.id || ''
    } else {
      const exists = categoryRows.value.some((item) => item.id === selectedTemplateCategoryId.value)
      if (!exists) {
        selectedTemplateCategoryId.value = categoryRows.value[0]?.id || ''
      }
    }
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '加载分类失败'
    MessagePlugin.error(message)
  } finally {
    categoryLoading.value = false
  }
}

async function loadMetaTemplatePresets() {
  try {
    const response = await CategoryService.ListMetaTemplatePresets({})
    metaTemplatePresets.value = (response.presets || []).filter((item) => Boolean(String(item.value || '').trim()))

    if (!selectedMetaTemplatePreset.value) {
      return
    }
    const exists = metaTemplatePresets.value.some((item) => item.value === selectedMetaTemplatePreset.value)
    if (!exists) {
      selectedMetaTemplatePreset.value = ''
    }
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '加载模板预设失败'
    MessagePlugin.error(message)
  }
}

function slugify(input: string): string {
  return input
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9\u4e00-\u9fa5\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
}

function resetCategoryForm() {
  selectedCategoryPreset.value = ''
  categoryForm.name = ''
  categoryForm.slug = ''
  categoryForm.description = ''
  categoryForm.parentId = ''
  categoryForm.icon = ''
  categoryForm.sortOrder = 0
  categoryForm.isEnabled = true
}

function suggestNextCategorySortOrder(): number {
  if (categoryRows.value.length === 0) {
    return 0
  }
  let maxSort = 0
  for (const item of categoryRows.value) {
    const value = Number(item.sortOrder || 0)
    if (value > maxSort) {
      maxSort = value
    }
  }
  return maxSort + 10
}

function applySelectedCategoryPreset() {
  const preset = selectedCategoryPresetItem.value
  if (!preset) {
    MessagePlugin.warning('请先选择一个预设')
    return
  }

  categoryForm.name = preset.name
  categoryForm.slug = preset.slug
  categoryForm.icon = preset.icon
  categoryForm.description = preset.description

  MessagePlugin.success(`已填充预设：${preset.label}，你可以继续修改后再创建`)
}

function normalizePresetTemplateType(type: string | undefined): string {
  const fieldType = String(type || 'text').trim().toLowerCase()
  if (['text', 'textarea', 'select', 'number', 'date', 'url'].includes(fieldType)) {
    return fieldType
  }
  return 'text'
}

async function applyMetaTemplatePreset() {
  const categoryId = selectedTemplateCategoryId.value.trim()
  if (!categoryId) {
    MessagePlugin.warning('请先选择一个分类')
    return
  }

  const preset = selectedMetaTemplatePresetItem.value
  if (!preset) {
    MessagePlugin.warning('请先选择一个模板预设')
    return
  }
  const presetTemplates = (preset.templates || []).filter((item) => Boolean(String(item.key || '').trim()))
  const presetName = preset.label || preset.value || '未命名预设'

  const existingKeys = new Set(
    selectedCategoryMetaTemplates.value
      .map((item) => item.key.trim().toLowerCase())
      .filter(Boolean),
  )

  const candidates = presetTemplates.filter((item) => {
    const key = String(item.key || '').trim().toLowerCase()
    return Boolean(key) && !existingKeys.has(key)
  })

  if (candidates.length === 0) {
    MessagePlugin.info('该预设字段已全部存在，无需添加')
    return
  }

  metaTemplatePresetApplying.value = true
  try {
    for (const [index, item] of candidates.entries()) {
      const key = String(item.key || '').trim().toLowerCase()
      const label = String(item.label || key).trim()
      const fieldType = normalizePresetTemplateType(item.type)
      const sortOrder = Number(item.sortOrder || 0) || (index + 1) * 10

      await CategoryService.CreateMetaTemplate({
        template: {
          categoryId,
          key,
          label,
          type: fieldType,
          required: item.required === true,
          options: String(item.options || '').trim(),
          sortOrder,
          defaultValue: String(item.defaultValue || '').trim(),
        },
      })
    }

    await loadCategories()
    selectedTemplateCategoryId.value = categoryId

    const skipped = presetTemplates.length - candidates.length
    MessagePlugin.success(`已应用预设「${presetName}」：新增 ${candidates.length}，跳过 ${Math.max(0, skipped)}`)
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '应用模板预设失败'
    MessagePlugin.error(message)
  } finally {
    metaTemplatePresetApplying.value = false
  }
}

function openCreateCategory() {
  categoryEditorMode.value = 'create'
  editingCategoryId.value = ''
  resetCategoryForm()
  categoryForm.sortOrder = suggestNextCategorySortOrder()
  categoryEditorVisible.value = true
}

function openEditCategory(row: CategoryRow) {
  categoryEditorMode.value = 'edit'
  editingCategoryId.value = row.id
  selectedCategoryPreset.value = ''
  categoryForm.name = row.name
  categoryForm.slug = row.slug
  categoryForm.description = row.description
  categoryForm.parentId = row.parentId
  categoryForm.icon = row.icon
  categoryForm.sortOrder = row.sortOrder
  categoryForm.isEnabled = row.isEnabled
  categoryEditorVisible.value = true
}

function focusCategoryTemplates(row: CategoryRow) {
  selectedTemplateCategoryId.value = row.id
  metaTemplateEditorVisible.value = false
}

function closeCategoryEditor() {
  categoryEditorVisible.value = false
  editingCategoryId.value = ''
  resetCategoryForm()
}

function resetMetaTemplateForm() {
  metaTemplateForm.categoryId = selectedTemplateCategoryId.value
  metaTemplateForm.key = ''
  metaTemplateForm.label = ''
  metaTemplateForm.type = 'text'
  metaTemplateForm.required = false
  metaTemplateForm.options = ''
  metaTemplateForm.sortOrder = 0
  metaTemplateForm.defaultValue = ''
}

function openCreateMetaTemplate() {
  if (!selectedTemplateCategoryId.value) {
    MessagePlugin.warning('请先在分类列表点击“模板”选择管理对象')
    return
  }

  metaTemplateEditorMode.value = 'create'
  editingMetaTemplateId.value = ''
  resetMetaTemplateForm()
  metaTemplateEditorVisible.value = true
}

function openEditMetaTemplate(row: CategoryMetaTemplateRow) {
  if (row.categoryId && row.categoryId !== selectedTemplateCategoryId.value) {
    selectedTemplateCategoryId.value = row.categoryId
  }
  metaTemplateEditorMode.value = 'edit'
  editingMetaTemplateId.value = row.id
  metaTemplateForm.categoryId = row.categoryId || selectedTemplateCategoryId.value
  metaTemplateForm.key = row.key
  metaTemplateForm.label = row.label
  metaTemplateForm.type = row.type || 'text'
  metaTemplateForm.required = row.required
  metaTemplateForm.options = row.options
  metaTemplateForm.sortOrder = row.sortOrder
  metaTemplateForm.defaultValue = row.defaultValue
  metaTemplateEditorVisible.value = true
}

function closeMetaTemplateEditor() {
  metaTemplateEditorVisible.value = false
  editingMetaTemplateId.value = ''
  resetMetaTemplateForm()
}

async function saveMetaTemplate() {
  const categoryId = (metaTemplateForm.categoryId || selectedTemplateCategoryId.value).trim()
  const key = metaTemplateForm.key.trim()
  const label = metaTemplateForm.label.trim()
  const fieldType = metaTemplateForm.type.trim().toLowerCase()

  if (!categoryId) {
    MessagePlugin.warning('请先在分类列表选择模板管理对象')
    return
  }
  if (!key) {
    MessagePlugin.warning('模板 Key 不能为空')
    return
  }
  if (!label) {
    MessagePlugin.warning('模板显示名不能为空')
    return
  }
  if (!fieldType) {
    MessagePlugin.warning('模板类型不能为空')
    return
  }
  if (fieldType === 'select' && !metaTemplateForm.options.trim()) {
    MessagePlugin.warning('select 类型需要填写可选项')
    return
  }

  metaTemplateSaving.value = true
  try {
    if (metaTemplateEditorMode.value === 'create') {
      await CategoryService.CreateMetaTemplate({
        template: {
          categoryId,
          key,
          label,
          type: fieldType,
          required: metaTemplateForm.required,
          options: metaTemplateForm.options.trim(),
          sortOrder: metaTemplateForm.sortOrder,
          defaultValue: metaTemplateForm.defaultValue.trim(),
        },
      })
      MessagePlugin.success('模板创建成功')
    } else {
      await CategoryService.UpdateMetaTemplate({
        template: {
          id: editingMetaTemplateId.value,
          categoryId,
          key,
          label,
          type: fieldType,
          required: metaTemplateForm.required,
          options: metaTemplateForm.options.trim(),
          sortOrder: metaTemplateForm.sortOrder,
          defaultValue: metaTemplateForm.defaultValue.trim(),
        },
      })
      MessagePlugin.success('模板更新成功')
    }

    closeMetaTemplateEditor()
    await loadCategories()
    selectedTemplateCategoryId.value = categoryId
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '保存模板失败'
    MessagePlugin.error(message)
  } finally {
    metaTemplateSaving.value = false
  }
}

async function removeMetaTemplate(row: CategoryMetaTemplateRow) {
  const confirmed = window.confirm(`确认删除模板「${row.label} (${row.key})」吗？`)
  if (!confirmed) {
    return
  }

  metaTemplateLoading.value = true
  try {
    await CategoryService.DeleteMetaTemplate({ id: row.id })
    MessagePlugin.success('模板已删除')
    await loadCategories()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '删除模板失败'
    MessagePlugin.error(message)
  } finally {
    metaTemplateLoading.value = false
  }
}

async function saveCategory() {
  const name = categoryForm.name.trim()
  const slug = (categoryForm.slug.trim() || slugify(name)).trim()

  if (!name) {
    MessagePlugin.warning('类别名称不能为空')
    return
  }

  if (!slug) {
    MessagePlugin.warning('Slug 不能为空')
    return
  }

  categorySaving.value = true
  try {
    if (categoryEditorMode.value === 'create') {
      const created = await CategoryService.CreateCategory({
        name,
        slug,
        description: categoryForm.description.trim(),
        parentId: categoryForm.parentId,
        icon: categoryForm.icon.trim(),
        sortOrder: categoryForm.sortOrder,
      })
      if (created.id) {
        selectedTemplateCategoryId.value = created.id
      }
      MessagePlugin.success('类别创建成功')
    } else {
      await CategoryService.UpdateCategory({
        category: {
          id: editingCategoryId.value,
          name,
          slug,
          description: categoryForm.description.trim(),
          parentId: categoryForm.parentId,
          icon: categoryForm.icon.trim(),
          sortOrder: categoryForm.sortOrder,
          isEnabled: categoryForm.isEnabled,
        },
      })
      MessagePlugin.success('类别更新成功')
    }

    closeCategoryEditor()
    await loadCategories()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '保存类别失败'
    MessagePlugin.error(message)
  } finally {
    categorySaving.value = false
  }
}

async function removeCategory(row: CategoryRow) {
  const confirmed = window.confirm(`确认删除类别「${row.name}」吗？`)
  if (!confirmed) {
    return
  }

  try {
    await CategoryService.DeleteCategory({ id: row.id })
    MessagePlugin.success('类别已删除')
    await loadCategories()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '删除类别失败'
    MessagePlugin.error(message)
  }
}

onMounted(async () => {
  await Promise.all([loadCategories(), loadMetaTemplatePresets()])
})
</script>

<template>
  <t-card title="种子类别管理" size="small" class="h-full">
    <template #actions>
      <div class="card-title-actions">
        <template v-if="!categoryEditorVisible">
          <t-button theme="primary" @click="openCreateCategory">新建类别</t-button>
          <t-button variant="outline" :loading="categoryLoading" @click="loadCategories">刷新列表</t-button>
        </template>
        <template v-else>
          <t-button variant="outline" @click="closeCategoryEditor">返回列表</t-button>
        </template>
      </div>
    </template>

    <div v-if="categoryEditorVisible" class="category-editor category-editor-focus">
      <p class="mb-3 text-sm font-600">
        {{ categoryEditorMode === 'create' ? '新建类别' : '编辑类别' }}
      </p>

      <div class="form-grid">
        <label v-if="categoryEditorMode === 'create'" class="lg:col-span-2">
          <span>快速预设</span>
          <div class="preset-picker-row">
            <t-select
              v-model="selectedCategoryPreset"
              :options="categoryPresetOptions"
              clearable
              placeholder="选择常见类别预设"
            />
            <t-button variant="outline" @click="applySelectedCategoryPreset">填充表单</t-button>
          </div>
          <p class="preset-hint">用于填充名称、Slug、图标和描述。创建后再到“元数据模板”里配置字段。</p>
        </label>

        <label>
          <span>名称 *</span>
          <t-input v-model="categoryForm.name" clearable placeholder="例如：Movies" />
        </label>

        <label>
          <span>Slug *</span>
          <t-input v-model="categoryForm.slug" clearable placeholder="例如：movies" />
        </label>

        <label>
          <span>父类别</span>
          <t-select v-model="categoryForm.parentId" clearable placeholder="无">
            <t-option
              v-for="item in parentCategoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </t-select>
        </label>

        <label>
          <span>图标</span>
          <t-input v-model="categoryForm.icon" clearable placeholder="例如：🎬" />
        </label>

        <label>
          <span>排序</span>
          <t-input-number v-model="categoryForm.sortOrder" :min="0" />
        </label>

        <label v-if="categoryEditorMode === 'edit'">
          <span>启用</span>
          <t-switch
            :value="categoryForm.isEnabled"
            @change="(value: boolean) => (categoryForm.isEnabled = value)"
          />
        </label>

        <label class="lg:col-span-2">
          <span>描述</span>
          <t-textarea v-model="categoryForm.description" :autosize="{ minRows: 2, maxRows: 4 }" />
        </label>
      </div>

      <div class="mt-3 flex gap-2">
        <t-button theme="primary" :loading="categorySaving" @click="saveCategory">保存类别</t-button>
        <t-button variant="outline" @click="closeCategoryEditor">取消</t-button>
      </div>
    </div>

    <template v-else>
      <t-table
        row-key="id"
        size="small"
        bordered
        hover
        :columns="categoryColumns"
        :data="categoryRows"
        :loading="categoryLoading"
      >
        <template #name="{ row }">
          <div class="flex items-center gap-1" :style="{ paddingLeft: `${row.depth * 16}px` }">
            <span>{{ row.icon || '📁' }}</span>
            <span>{{ row.name }}</span>
          </div>
        </template>

        <template #isEnabled="{ row }">
          <t-tag size="small" :theme="row.isEnabled ? 'success' : 'danger'" variant="light">
            {{ row.isEnabled ? '启用' : '禁用' }}
          </t-tag>
        </template>

        <template #actions="{ row }">
          <div class="flex items-center gap-2">
            <t-link
              hover="color"
              :theme="row.id === selectedTemplateCategoryId ? 'success' : 'default'"
              @click="focusCategoryTemplates(row)"
            >
              模板
            </t-link>
            <t-link theme="primary" hover="color" @click="openEditCategory(row)">编辑</t-link>
            <t-link theme="danger" hover="color" @click="removeCategory(row)">删除</t-link>
          </div>
        </template>
      </t-table>

      <div class="meta-template-panel">
        <div class="meta-template-header">
          <div class="meta-template-title">
            <p class="text-sm font-600">元数据模板</p>
            <p class="preset-hint">
              当前分类：{{ selectedTemplateCategory?.name || '未选择' }}（在上方类别操作里点击“模板”切换）
            </p>
          </div>
        </div>

        <div class="meta-preset-row">
          <t-button theme="primary" :disabled="!selectedTemplateCategoryId" @click="openCreateMetaTemplate">
            新增模板
          </t-button>
          <t-select
            v-model="selectedMetaTemplatePreset"
            :options="metaTemplatePresetOptions"
            clearable
            class="w-64"
            placeholder="筛选预设（电影/动漫/电视剧/纪录片/软件/电子书）"
          />
          <t-button
            variant="outline"
            :loading="metaTemplatePresetApplying"
            :disabled="!selectedTemplateCategoryId || !selectedMetaTemplatePreset"
            @click="applyMetaTemplatePreset"
          >
            应用预设
          </t-button>
        </div>

        <p v-if="selectedMetaTemplatePresetItem" class="preset-hint">
          当前筛选：{{ selectedMetaTemplatePresetItem.label || selectedMetaTemplatePresetItem.value }}，可一键导入该预设字段。
        </p>
        <div v-if="selectedMetaTemplatePresetItem" class="preset-template-preview">
          <t-tag
            v-for="item in (selectedMetaTemplatePresetItem.templates || [])"
            :key="item.key"
            size="small"
            variant="light"
          >
            {{ item.label }} ({{ item.key }})
            <span v-if="item.required"> *</span>
          </t-tag>
        </div>

        <t-table
          row-key="id"
          size="small"
          bordered
          hover
          :columns="categoryMetaTemplateColumns"
          :data="filteredCategoryMetaTemplates"
          :loading="categoryLoading || metaTemplateLoading"
        >
          <template #required="{ row }">
            <t-tag size="small" :theme="row.required ? 'success' : 'default'" variant="light">
              {{ row.required ? '是' : '否' }}
            </t-tag>
          </template>

          <template #actions="{ row }">
            <div class="flex items-center gap-2">
              <t-link theme="primary" hover="color" @click="openEditMetaTemplate(row)">编辑</t-link>
              <t-link theme="danger" hover="color" @click="removeMetaTemplate(row)">删除</t-link>
            </div>
          </template>
        </t-table>

        <div v-if="metaTemplateEditorVisible" class="category-editor">
          <p class="mb-3 text-sm font-600">
            {{ metaTemplateEditorMode === 'create' ? '新建模板' : '编辑模板' }}
          </p>
          <p class="mb-3 text-xs text-[var(--muted-text)]">
            所属分类：{{ selectedTemplateCategory?.name || '未选择' }}
          </p>

          <div class="form-grid">
            <label>
              <span>字段 Key *</span>
              <t-input v-model="metaTemplateForm.key" clearable placeholder="例如：author / isbn / season" />
            </label>

            <label>
              <span>字段名称 *</span>
              <t-input v-model="metaTemplateForm.label" clearable placeholder="例如：作者 / ISBN / 季度" />
            </label>

            <label>
              <span>字段类型 *</span>
              <t-select v-model="metaTemplateForm.type" :options="metaTemplateTypeOptions" />
            </label>

            <label>
              <span>排序</span>
              <t-input-number v-model="metaTemplateForm.sortOrder" :min="0" />
            </label>

            <label>
              <span>必填</span>
              <t-switch
                :value="metaTemplateForm.required"
                @change="(value: boolean) => (metaTemplateForm.required = value)"
              />
            </label>

            <label>
              <span>默认值</span>
              <t-input v-model="metaTemplateForm.defaultValue" clearable />
            </label>

            <label class="lg:col-span-2">
              <span>可选项（select 类型使用，逗号或换行分隔）</span>
              <t-textarea
                v-model="metaTemplateForm.options"
                :autosize="{ minRows: 2, maxRows: 4 }"
                placeholder="例如：1080p, 2160p 或一行一个选项"
              />
            </label>
          </div>

          <div class="mt-3 flex gap-2">
            <t-button theme="primary" :loading="metaTemplateSaving" @click="saveMetaTemplate">保存模板</t-button>
            <t-button variant="outline" @click="closeMetaTemplateEditor">取消</t-button>
          </div>
        </div>
      </div>
    </template>
  </t-card>
</template>

<style scoped>
.form-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  gap: 12px;
}

@media (min-width: 1024px) {
  .form-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 1023px) {
  .preset-picker-row {
    grid-template-columns: 1fr;
  }
}

.form-grid label {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 6px;
  font-size: 12px;
  color: var(--muted-text);
}

.form-grid label :deep(.t-switch),
.form-grid label :deep(.t-checkbox),
.form-grid label :deep(.t-radio-group) {
  align-self: flex-start;
  width: auto;
}

.category-editor {
  margin-top: 14px;
  border-top: 1px dashed var(--app-border);
  padding-top: 14px;
}

.preset-picker-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 8px;
  align-items: center;
}

.preset-hint {
  margin: 0;
  font-size: 12px;
  color: var(--muted-text);
}

.preset-template-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.meta-preset-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.meta-template-panel {
  margin-top: 18px;
  border-top: 1px dashed var(--app-border);
  padding-top: 14px;
  display: grid;
  gap: 12px;
}

.meta-template-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.meta-template-title {
  display: grid;
  gap: 4px;
}

.card-title-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
}
</style>
