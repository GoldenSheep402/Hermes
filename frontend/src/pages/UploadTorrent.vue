<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Category } from '@/lib/proto/category/v1/category.pb'
import type { ResourceMeta } from '@/lib/proto/resource/v1/resource.pb'
import { CategoryService, ResourceService } from '@/services/grpc'

interface CategoryOption {
  id: string
  label: string
}

interface UploadFormState {
  title: string
  subtitle: string
  description: string
  categoryId: string
  imdbId: string
  imdbRating: string
  source: string
  resolution: string
  poster: string
  tags: string
  screenshots: string
}

const router = useRouter()
const fileInputRef = ref<HTMLInputElement | null>(null)

const categories = ref<CategoryOption[]>([])
const loadingCategories = ref(false)
const submitting = ref(false)
const selectedFile = ref<File | null>(null)

const form = reactive<UploadFormState>({
  title: '',
  subtitle: '',
  description: '',
  categoryId: '',
  imdbId: '',
  imdbRating: '',
  source: '',
  resolution: '',
  poster: '',
  tags: '',
  screenshots: '',
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

function flattenCategories(nodes: Category[] | undefined): Category[] {
  const queue = [...(nodes || [])]
  const output: Category[] = []

  while (queue.length > 0) {
    const current = queue.shift()
    if (!current) {
      continue
    }
    output.push(current)
    if (current.children && current.children.length > 0) {
      queue.push(...current.children)
    }
  }

  return output
}

async function loadCategories() {
  loadingCategories.value = true
  try {
    const response = await CategoryService.ListCategories({})
    const flattened = flattenCategories(response.categories)
    categories.value = flattened
      .filter((item) => item.isEnabled !== false)
      .map((item) => ({
        id: item.id || '',
        label: item.name || item.slug || item.id || '未命名分类',
      }))
      .filter((item) => item.id)

    if (categories.value.length > 0) {
      form.categoryId = categories.value[0].id
    }
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

function buildMetadata(): ResourceMeta[] {
  const metadata: ResourceMeta[] = []
  const append = (key: string, value: string) => {
    const text = value.trim()
    if (text) {
      metadata.push({ key, value: text })
    }
  }

  append('imdb', form.imdbId)
  append('imdb_rating', form.imdbRating)
  append('source', form.source)
  append('resolution', form.resolution)
  append('poster', form.poster)

  return metadata
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
            placeholder="填写影片/资源介绍、压制信息、发布说明等"
          />
        </label>

        <label>
          <span>IMDB ID</span>
          <t-input v-model="form.imdbId" clearable placeholder="例如：tt15239678" />
        </label>

        <label>
          <span>IMDB 评分</span>
          <t-input v-model="form.imdbRating" clearable placeholder="例如：8.5" />
        </label>

        <label>
          <span>来源 Source</span>
          <t-input v-model="form.source" clearable placeholder="例如：Blu-ray / WEB-DL / Remux" />
        </label>

        <label>
          <span>分辨率</span>
          <t-input v-model="form.resolution" clearable placeholder="例如：4K / 1080p" />
        </label>

        <label class="lg:col-span-2">
          <span>海报 URL</span>
          <t-input v-model="form.poster" clearable placeholder="封面图地址" />
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
        <li>标签、截图、元数据会写入 `CreateResource` 请求。</li>
        <li>发布成功后会跳转回种子列表页。</li>
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

.guide-list {
  margin: 0;
  padding-left: 18px;
  display: grid;
  gap: 8px;
  font-size: 12px;
  color: var(--muted-text);
}
</style>
