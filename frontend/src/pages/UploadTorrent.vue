<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Category, CategoryMetaTemplate } from '@/lib/proto/category/v1/category.pb'
import type { ResourceMeta } from '@/lib/proto/resource/v1/resource.pb'
import { CategoryService, ResourceService } from '@/services/grpc'

interface CategoryOption {
  id: string
  label: string
  templates: CategoryMetaTemplate[]
}

interface UploadFormState {
  title: string
  subtitle: string
  description: string
  categoryId: string
  tags: string
  screenshots: string
}

interface CustomMetaEntry {
  id: string
  key: string
  value: string
}

const router = useRouter()
const fileInputRef = ref<HTMLInputElement | null>(null)

const categories = ref<CategoryOption[]>([])
const loadingCategories = ref(false)
const submitting = ref(false)
const selectedFile = ref<File | null>(null)
const templateValues = reactive<Record<string, string>>({})
const customMetadata = ref<CustomMetaEntry[]>([])

const form = reactive<UploadFormState>({
  title: '',
  subtitle: '',
  description: '',
  categoryId: '',
  tags: '',
  screenshots: '',
})

const selectedCategory = computed(() => categories.value.find((item) => item.id === form.categoryId) || null)
const selectedCategoryTemplates = computed(() => {
  const templates = selectedCategory.value?.templates || []
  return templates
    .filter((item) => Boolean(item.key && item.key.trim()))
    .sort((left, right) => {
      const leftOrder = Number(left.sortOrder || 0)
      const rightOrder = Number(right.sortOrder || 0)
      if (leftOrder === rightOrder) {
        return String(left.key || '').localeCompare(String(right.key || ''))
      }
      return leftOrder - rightOrder
    })
})

const canSubmit = computed(
  () =>
    Boolean(
      form.title.trim() &&
        form.categoryId.trim() &&
        selectedFile.value &&
        !submitting.value,
    ),
)

const fileLabel = computed(() => {
  if (!selectedFile.value) {
    return '未选择文件'
  }
  return `${selectedFile.value.name} (${formatFileSize(selectedFile.value.size)})`
})

function triggerFilePicker() {
  fileInputRef.value?.click()
}

function flattenCategories(nodes: Category[] | undefined, depth = 0): CategoryOption[] {
  const output: CategoryOption[] = []
  for (const node of nodes || []) {
    const id = String(node.id || '').trim()
    if (!id || node.isEnabled === false) {
      continue
    }

    const name = String(node.name || node.slug || node.id || '未命名分类').trim()
    const label = `${'  '.repeat(depth)}${name}`
    output.push({
      id,
      label,
      templates: (node.metaTemplates || []).slice(),
    })

    if (node.children && node.children.length > 0) {
      output.push(...flattenCategories(node.children, depth + 1))
    }
  }
  return output
}

function createCustomMetaEntry(): CustomMetaEntry {
  return {
    id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
    key: '',
    value: '',
  }
}

function resetTemplateValues() {
  const current = { ...templateValues }
  for (const key of Object.keys(templateValues)) {
    delete templateValues[key]
  }

  for (const tpl of selectedCategoryTemplates.value) {
    const key = String(tpl.key || '').trim()
    if (!key) {
      continue
    }
    const previous = String(current[key] || '').trim()
    templateValues[key] = previous || String(tpl.defaultValue || '')
  }
}

function normalizeTemplateType(tpl: CategoryMetaTemplate): string {
  return String(tpl.type || '').trim().toLowerCase()
}

function parseTemplateOptions(tpl: CategoryMetaTemplate): string[] {
  const raw = String(tpl.options || '').trim()
  if (!raw) {
    return []
  }

  if (raw.startsWith('[') && raw.endsWith(']')) {
    try {
      const arr = JSON.parse(raw)
      if (Array.isArray(arr)) {
        return arr
          .map((item) => String(item || '').trim())
          .filter(Boolean)
      }
    } catch {
      // ignore invalid JSON and fallback to split parser
    }
  }

  const parts = raw
    .split(/[\n,，]/)
    .map((item) => item.trim())
    .filter(Boolean)

  return Array.from(new Set(parts))
}

function isTemplateSelect(tpl: CategoryMetaTemplate): boolean {
  const t = normalizeTemplateType(tpl)
  return ['select', 'enum', 'dropdown'].includes(t) && parseTemplateOptions(tpl).length > 0
}

function isTemplateTextarea(tpl: CategoryMetaTemplate): boolean {
  const t = normalizeTemplateType(tpl)
  return ['textarea', 'longtext', 'markdown'].includes(t)
}

function templatePlaceholder(tpl: CategoryMetaTemplate): string {
  const defaultValue = String(tpl.defaultValue || '').trim()
  if (defaultValue) {
    return `默认值：${defaultValue}`
  }
  return `请输入 ${tpl.label || tpl.key || '元数据'}`
}

function addCustomMetadata() {
  customMetadata.value.push(createCustomMetaEntry())
}

function removeCustomMetadata(id: string) {
  customMetadata.value = customMetadata.value.filter((item) => item.id !== id)
}

async function loadCategories() {
  loadingCategories.value = true
  try {
    const response = await CategoryService.ListCategories({})
    categories.value = flattenCategories(response.categories)

    if (categories.value.length > 0) {
      form.categoryId = categories.value[0].id
    }
    resetTemplateValues()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '分类加载失败'
    MessagePlugin.error(message)
  } finally {
    loadingCategories.value = false
  }
}

function formatFileSize(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) {
    return '0 B'
  }

  const units = ['B', 'KB', 'MB', 'GB']
  const power = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1)
  const value = bytes / 1024 ** power
  return `${value.toFixed(power > 1 ? 2 : 1)} ${units[power]}`
}

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) {
    selectedFile.value = null
    return
  }

  if (!file.name.toLowerCase().endsWith('.torrent')) {
    MessagePlugin.warning('请上传 .torrent 文件')
    target.value = ''
    selectedFile.value = null
    return
  }

  selectedFile.value = file
}

function validateTemplateRequired(): boolean {
  for (const tpl of selectedCategoryTemplates.value) {
    if (!tpl.required) {
      continue
    }
    const key = String(tpl.key || '').trim()
    if (!key) {
      continue
    }
    const value = String(templateValues[key] || '').trim()
    if (!value) {
      MessagePlugin.warning(`请填写必填元数据：${tpl.label || key}`)
      return false
    }
  }
  return true
}

function buildMetadata(): ResourceMeta[] {
  const metadataMap = new Map<string, string>()

  for (const tpl of selectedCategoryTemplates.value) {
    const key = String(tpl.key || '').trim()
    if (!key) {
      continue
    }
    const value = String(templateValues[key] || '').trim()
    if (!value) {
      continue
    }
    metadataMap.set(key, value)
  }

  for (const item of customMetadata.value) {
    const key = item.key.trim()
    const value = item.value.trim()
    if (!key || !value) {
      continue
    }
    metadataMap.set(key, value)
  }

  return Array.from(metadataMap.entries()).map(([key, value]) => ({ key, value }))
}

function buildTagNames(): string[] {
  const names = form.tags
    .split(/[\n,，]/)
    .map((item) => item.trim())
    .filter(Boolean)

  return Array.from(new Set(names))
}

function buildScreenshotUrls(): string[] {
  const urls = form.screenshots
    .split(/[\n,，]/)
    .map((item) => item.trim())
    .filter(Boolean)

  return Array.from(new Set(urls))
}

async function submitUpload() {
  if (!selectedFile.value) {
    MessagePlugin.warning('请选择种子文件')
    return
  }
  if (!form.title.trim()) {
    MessagePlugin.warning('标题不能为空')
    return
  }
  if (!form.categoryId.trim()) {
    MessagePlugin.warning('请选择分类')
    return
  }
  if (!validateTemplateRequired()) {
    return
  }

  submitting.value = true
  try {
    const buffer = await selectedFile.value.arrayBuffer()
    const torrentData = new Uint8Array(buffer)

    const response = await ResourceService.CreateResource({
      title: form.title.trim(),
      subtitle: form.subtitle.trim(),
      description: form.description.trim(),
      categoryId: form.categoryId,
      torrentData,
      metadata: buildMetadata(),
      tagNames: buildTagNames(),
      screenshotUrls: buildScreenshotUrls(),
    })

    MessagePlugin.success(`发布成功，资源 ID: ${response.id || 'N/A'}`)
    await router.push('/torrents')
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '上传失败'
    MessagePlugin.error(message)
  } finally {
    submitting.value = false
  }
}

watch(
  () => form.categoryId,
  () => {
    resetTemplateValues()
  },
)

onMounted(async () => {
  await loadCategories()
})
</script>

<template>
  <section class="grid gap-4 lg:grid-cols-[1.2fr_320px]">
    <t-card title="发布种子" size="small">
      <div class="form-grid">
        <label>
          <span>主标题 *</span>
          <t-input v-model="form.title" clearable placeholder="例如：Dune Part Two 2160p Blu-ray" />
        </label>

        <label>
          <span>副标题</span>
          <t-input v-model="form.subtitle" clearable placeholder="例如：中字 / HDR10+ / Atmos" />
        </label>

        <label>
          <span>分类 *</span>
          <t-select v-model="form.categoryId" :loading="loadingCategories" placeholder="选择分类">
            <t-option
              v-for="item in categories"
              :key="item.id"
              :label="item.label"
              :value="item.id"
            />
          </t-select>
        </label>

        <label>
          <span>文件 *</span>
          <div class="file-picker">
            <input
              ref="fileInputRef"
              class="hidden-file-input"
              hidden
              type="file"
              accept=".torrent"
              @change="onFileChange"
            />
            <div class="file-picker-row">
              <t-button size="small" variant="outline" @click="triggerFilePicker">选择文件</t-button>
              <span class="file-name">{{ fileLabel }}</span>
            </div>
          </div>
        </label>

        <label class="lg:col-span-2">
          <span>简介</span>
          <t-textarea
            v-model="form.description"
            :autosize="{ minRows: 4, maxRows: 8 }"
            placeholder="填写资源介绍、制作说明、使用说明等"
          />
        </label>

        <label
          v-for="tpl in selectedCategoryTemplates"
          :key="tpl.id || tpl.key"
          :class="isTemplateTextarea(tpl) ? 'lg:col-span-2' : ''"
        >
          <span>
            {{ tpl.label || tpl.key }}
            <span v-if="tpl.required" class="required-mark">*</span>
          </span>

          <t-select
            v-if="isTemplateSelect(tpl)"
            v-model="templateValues[tpl.key || '']"
            clearable
            :placeholder="templatePlaceholder(tpl)"
          >
            <t-option
              v-for="option in parseTemplateOptions(tpl)"
              :key="option"
              :label="option"
              :value="option"
            />
          </t-select>

          <t-textarea
            v-else-if="isTemplateTextarea(tpl)"
            v-model="templateValues[tpl.key || '']"
            :autosize="{ minRows: 2, maxRows: 6 }"
            :placeholder="templatePlaceholder(tpl)"
          />

          <t-input
            v-else
            v-model="templateValues[tpl.key || '']"
            clearable
            :placeholder="templatePlaceholder(tpl)"
          />
        </label>

        <label class="lg:col-span-2">
          <span>自定义元数据（可选）</span>
          <div class="custom-meta-list">
            <div v-if="customMetadata.length === 0" class="custom-meta-empty">
              当前未添加自定义元数据
            </div>

            <div v-for="item in customMetadata" :key="item.id" class="custom-meta-row">
              <t-input v-model="item.key" clearable placeholder="键，例如：author / isbn / game_version" />
              <t-input v-model="item.value" clearable placeholder="值" />
              <t-button variant="text" theme="danger" @click="removeCustomMetadata(item.id)">
                删除
              </t-button>
            </div>

            <t-button size="small" variant="outline" @click="addCustomMetadata">添加元数据</t-button>
          </div>
        </label>

        <label class="lg:col-span-2">
          <span>标签</span>
          <t-input v-model="form.tags" clearable placeholder="逗号分隔，例如：中字, Free, HDR" />
        </label>

        <label class="lg:col-span-2">
          <span>截图 URL</span>
          <t-textarea
            v-model="form.screenshots"
            :autosize="{ minRows: 2, maxRows: 4 }"
            placeholder="逗号或换行分隔多个地址"
          />
        </label>
      </div>

      <template #footer>
        <div class="flex gap-2">
          <t-button theme="primary" :loading="submitting" :disabled="!canSubmit" @click="submitUpload">
            发布资源
          </t-button>
          <router-link to="/torrents">
            <t-button variant="outline">返回列表</t-button>
          </router-link>
        </div>
      </template>
    </t-card>

    <t-card title="发布说明" size="small">
      <ul class="guide-list">
        <li>文件必须是 `.torrent` 格式。</li>
        <li>主标题和分类为必填项。</li>
        <li>分类模板元数据会自动显示，必填项需填写。</li>
        <li>你也可以额外添加任意键值元数据（用于非影视资源）。</li>
      </ul>
    </t-card>
  </section>
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

.form-grid label :deep(.t-switch),
.form-grid label :deep(.t-checkbox),
.form-grid label :deep(.t-radio-group) {
  align-self: flex-start;
  width: auto;
}

.required-mark {
  color: #ef4444;
  margin-left: 2px;
}

.file-picker {
  border: 1px dashed var(--app-border);
  border-radius: 8px;
  padding: 10px;
  background: var(--soft-bg);
}

.hidden-file-input {
  display: none;
}

.file-picker-row {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 28px;
}

.file-name {
  font-size: 12px;
  color: var(--muted-text);
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.custom-meta-list {
  display: grid;
  gap: 8px;
  padding: 10px;
  border: 1px dashed var(--app-border);
  border-radius: 8px;
  background: var(--soft-bg);
}

.custom-meta-empty {
  font-size: 12px;
  color: var(--muted-text);
}

.custom-meta-row {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 8px;
  align-items: center;
}

@media (max-width: 1023px) {
  .custom-meta-row {
    grid-template-columns: 1fr;
  }
}

.guide-list {
  margin: 0;
  padding-left: 18px;
  display: grid;
  gap: 8px;
  font-size: 12px;
  color: var(--muted-text);
}
</style>
