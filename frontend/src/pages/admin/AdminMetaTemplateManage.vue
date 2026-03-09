<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Category, CategoryMetaTemplate, MetaTemplatePreset } from '@/lib/proto/category/v1/category.pb'
import { CategoryService } from '@/services/grpc'
import type {
  AdminCategoryMetaTemplateRow as CategoryMetaTemplateRow,
  AdminCategoryRow as CategoryRow,
} from '@/types/category'

const categoryLoading = ref(false)
const metaTemplateLoading = ref(false)
const metaTemplateSaving = ref(false)
const metaTemplatePresetApplying = ref(false)
const categoryRows = ref<CategoryRow[]>([])
const selectedCategoryId = ref('')
const metaTemplatePresets = ref<MetaTemplatePreset[]>([])
const selectedMetaTemplatePreset = ref('')
const metaTemplateEditorVisible = ref(false)
const metaTemplateEditorMode = ref<'create' | 'edit'>('create')
const editingMetaTemplateId = ref('')

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

const categoryOptions = computed(() =>
  categoryRows.value.map((item) => ({
    label: `${'— '.repeat(item.depth)}${item.icon || '📁'} ${item.name}`,
    value: item.id,
  })),
)

const selectedCategory = computed(
  () => categoryRows.value.find((item) => item.id === selectedCategoryId.value) || null,
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

const categoryMetaTemplates = computed(() => {
  const items = selectedCategory.value?.metaTemplates || []
  return [...items].sort((left, right) => {
    if (left.sortOrder === right.sortOrder) {
      return left.key.localeCompare(right.key)
    }
    return left.sortOrder - right.sortOrder
  })
})

const filteredCategoryMetaTemplates = computed(() => {
  const preset = selectedMetaTemplatePresetItem.value
  if (!preset) {
    return categoryMetaTemplates.value
  }

  const templates = preset.templates || []
  const keySet = new Set(
    templates
      .map((item) => String(item.key || '').trim().toLowerCase())
      .filter(Boolean),
  )
  return categoryMetaTemplates.value.filter((item) => keySet.has(item.key.trim().toLowerCase()))
})

watch(selectedCategoryId, (value) => {
  if (!metaTemplateEditorVisible.value || metaTemplateEditorMode.value !== 'create') {
    return
  }
  metaTemplateForm.categoryId = value
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

    if (!selectedCategoryId.value) {
      selectedCategoryId.value = categoryRows.value[0]?.id || ''
    } else {
      const exists = categoryRows.value.some((item) => item.id === selectedCategoryId.value)
      if (!exists) {
        selectedCategoryId.value = categoryRows.value[0]?.id || ''
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
    metaTemplatePresets.value = (response.presets || []).filter((item) =>
      Boolean(String(item.value || '').trim()),
    )

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

function normalizePresetTemplateType(type: string | undefined): string {
  const fieldType = String(type || 'text').trim().toLowerCase()
  if (['text', 'textarea', 'select', 'number', 'date', 'url'].includes(fieldType)) {
    return fieldType
  }
  return 'text'
}

function suggestNextMetaTemplateSortOrder(): number {
  if (categoryMetaTemplates.value.length === 0) {
    return 10
  }
  let maxSort = 0
  for (const item of categoryMetaTemplates.value) {
    const value = Number(item.sortOrder || 0)
    if (value > maxSort) {
      maxSort = value
    }
  }
  return maxSort + 10
}

function findDuplicateMetaTemplateKey(key: string, ignoreID = ''): boolean {
  const target = key.trim().toLowerCase()
  if (!target) {
    return false
  }
  return categoryMetaTemplates.value.some((item) => item.id !== ignoreID && item.key.trim().toLowerCase() === target)
}

function resetMetaTemplateForm() {
  metaTemplateForm.categoryId = selectedCategoryId.value
  metaTemplateForm.key = ''
  metaTemplateForm.label = ''
  metaTemplateForm.type = 'text'
  metaTemplateForm.required = false
  metaTemplateForm.options = ''
  metaTemplateForm.sortOrder = suggestNextMetaTemplateSortOrder()
  metaTemplateForm.defaultValue = ''
}

async function applyMetaTemplatePreset() {
  const categoryID = selectedCategoryId.value.trim()
  if (!categoryID) {
    MessagePlugin.warning('请先选择一个分类')
    return
  }

  const preset = selectedMetaTemplatePresetItem.value
  if (!preset) {
    MessagePlugin.warning('请先选择一个模板预设')
    return
  }
  const presetTemplates = (preset.templates || []).filter((item) =>
    Boolean(String(item.key || '').trim()),
  )
  const presetName = preset.label || preset.value || '未命名预设'

  const existingKeys = new Set(
    categoryMetaTemplates.value
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
          categoryId: categoryID,
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
    selectedCategoryId.value = categoryID
    const skipped = presetTemplates.length - candidates.length
    MessagePlugin.success(
      `已应用预设「${presetName}」：新增 ${candidates.length}，跳过 ${Math.max(0, skipped)}`,
    )
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '应用模板预设失败'
    MessagePlugin.error(message)
  } finally {
    metaTemplatePresetApplying.value = false
  }
}

function openCreateMetaTemplate() {
  if (!selectedCategoryId.value) {
    MessagePlugin.warning('请先选择一个分类')
    return
  }

  metaTemplateEditorMode.value = 'create'
  editingMetaTemplateId.value = ''
  resetMetaTemplateForm()
  metaTemplateEditorVisible.value = true
}

function openEditMetaTemplate(row: CategoryMetaTemplateRow) {
  if (row.categoryId && row.categoryId !== selectedCategoryId.value) {
    selectedCategoryId.value = row.categoryId
  }
  metaTemplateEditorMode.value = 'edit'
  editingMetaTemplateId.value = row.id
  metaTemplateForm.categoryId = row.categoryId || selectedCategoryId.value
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
  const categoryID = (metaTemplateForm.categoryId || selectedCategoryId.value).trim()
  const key = metaTemplateForm.key.trim()
  const label = metaTemplateForm.label.trim()
  const fieldType = metaTemplateForm.type.trim().toLowerCase()
  const ignoreID = metaTemplateEditorMode.value === 'edit' ? editingMetaTemplateId.value : ''
  const sortOrder = Number(metaTemplateForm.sortOrder || 0)

  if (!categoryID) {
    MessagePlugin.warning('请先选择一个分类')
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
  if (findDuplicateMetaTemplateKey(key, ignoreID)) {
    MessagePlugin.warning('模板 Key 已存在，请更换后再保存')
    return
  }

  metaTemplateSaving.value = true
  try {
    if (metaTemplateEditorMode.value === 'create') {
      await CategoryService.CreateMetaTemplate({
        template: {
          categoryId: categoryID,
          key,
          label,
          type: fieldType,
          required: metaTemplateForm.required,
          options: metaTemplateForm.options.trim(),
          sortOrder,
          defaultValue: metaTemplateForm.defaultValue.trim(),
        },
      })
      MessagePlugin.success('字段创建成功')
    } else {
      await CategoryService.UpdateMetaTemplate({
        template: {
          id: editingMetaTemplateId.value,
          categoryId: categoryID,
          key,
          label,
          type: fieldType,
          required: metaTemplateForm.required,
          options: metaTemplateForm.options.trim(),
          sortOrder,
          defaultValue: metaTemplateForm.defaultValue.trim(),
        },
      })
      MessagePlugin.success('字段更新成功')
    }

    closeMetaTemplateEditor()
    await loadCategories()
    selectedCategoryId.value = categoryID
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '保存字段失败'
    MessagePlugin.error(message)
  } finally {
    metaTemplateSaving.value = false
  }
}

async function removeMetaTemplate(row: CategoryMetaTemplateRow) {
  const confirmed = window.confirm(`确认删除字段「${row.label} (${row.key})」吗？`)
  if (!confirmed) {
    return
  }

  metaTemplateLoading.value = true
  try {
    await CategoryService.DeleteMetaTemplate({ id: row.id })
    MessagePlugin.success('字段已删除')
    await loadCategories()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '删除字段失败'
    MessagePlugin.error(message)
  } finally {
    metaTemplateLoading.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadCategories(), loadMetaTemplatePresets()])
})
</script>

<template>
  <t-card title="元数据模板管理" size="small" class="h-full">
    <template #actions>
      <div class="card-title-actions">
        <t-button theme="primary" :disabled="!selectedCategoryId" @click="openCreateMetaTemplate">新增元数据字段</t-button>
        <t-button variant="outline" :loading="categoryLoading" @click="loadCategories">刷新列表</t-button>
      </div>
    </template>

    <div class="meta-template-panel">
      <div class="meta-template-header">
        <div class="meta-template-title">
          <p class="text-sm font-600">当前分类</p>
          <p class="preset-hint">
            {{ selectedCategory?.name ? `${selectedCategory.icon || '📁'} ${selectedCategory.name}` : '未选择' }}
          </p>
        </div>
      </div>

      <div class="meta-preset-row">
        <t-select
          v-model="selectedCategoryId"
          :options="categoryOptions"
          clearable
          class="w-80"
          placeholder="选择分类"
        />
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
          :disabled="!selectedCategoryId || !selectedMetaTemplatePreset"
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
          {{ metaTemplateEditorMode === 'create' ? '新建字段' : '编辑字段' }}
        </p>
        <p class="mb-3 text-xs text-[var(--muted-text)]">
          所属分类：{{ selectedCategory?.name || '未选择' }}
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
          <t-button theme="primary" :loading="metaTemplateSaving" @click="saveMetaTemplate">保存字段</t-button>
          <t-button variant="outline" @click="closeMetaTemplateEditor">取消</t-button>
        </div>
      </div>
    </div>
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

.form-grid label {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 6px;
  font-size: 12px;
  color: var(--muted-text);
}

.form-grid label :deep(.t-switch) {
  align-self: flex-start;
  width: auto;
}

.category-editor {
  margin-top: 14px;
  border-top: 1px dashed var(--app-border);
  padding-top: 14px;
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
