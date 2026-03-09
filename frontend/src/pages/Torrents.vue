<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { SearchIcon } from 'tdesign-icons-vue-next'
import type { Category } from '@/lib/proto/category/v1/category.pb'
import type { Resource, ResourceMeta } from '@/lib/proto/resource/v1/resource.pb'
import type { TorrentInfo } from '@/lib/proto/torrent/v1/torrent.pb'
import { CategoryService, ResourceService, TorrentService } from '@/services/grpc'
import { formatBytes, formatDateTime } from '@/utils/format'

type PromotionKey = 'normal' | 'free' | '2xfree' | '50down'

interface CategoryOption {
  id: string
  name: string
  icon: string
}

interface TorrentRecord {
  id: string
  categoryId: string
  categoryLabel: string
  categoryIcon: string
  title: string
  subtitle: string
  tags: string[]
  comments: number
  files: number
  addedAt: string
  size: number
  seeders: number
  leechers: number
  completed: number
  uploader: string
  anonymous: boolean
  imdbId: string
  imdbRating: number
  poster: string
  source: string
  resolution: string
  promotion: PromotionKey
  status: 'alive' | 'dead'
  downloaded: boolean
}

interface FilterState {
  keyword: string
  categories: string[]
  source: string
  resolution: string
  promotion: '' | PromotionKey
  status: 'all' | 'alive' | 'dead'
  undownloadedOnly: boolean
}

const defaultCategories: CategoryOption[] = [
  { id: 'movie', name: '电影', icon: '🎬' },
  { id: 'series', name: '剧集', icon: '📺' },
  { id: 'documentary', name: '纪录片', icon: '🎞️' },
  { id: 'anime', name: '动画', icon: '🧩' },
  { id: 'music', name: '音乐', icon: '🎵' },
]

const loading = ref(false)
const isMobile = ref(false)
const categoryOptions = ref<CategoryOption[]>([])
const torrents = ref<TorrentRecord[]>([])
const loadError = ref('')

const filters = reactive<FilterState>({
  keyword: '',
  categories: [],
  source: '',
  resolution: '',
  promotion: '',
  status: 'all',
  undownloadedOnly: false,
})

const sourceOptions = computed(() => {
  const set = new Set(['Blu-ray', 'WEB-DL', 'Remux', 'Encode'])
  for (const item of torrents.value) {
    if (item.source) {
      set.add(item.source)
    }
  }
  return Array.from(set)
})

const resolutionOptions = computed(() => {
  const set = new Set(['4K', '1080p', '720p'])
  for (const item of torrents.value) {
    if (item.resolution) {
      set.add(item.resolution)
    }
  }
  return Array.from(set)
})

function handleResize() {
  isMobile.value = window.innerWidth < 900
}

function toggleCategory(categoryId: string) {
  const index = filters.categories.indexOf(categoryId)
  if (index >= 0) {
    filters.categories.splice(index, 1)
  } else {
    filters.categories.push(categoryId)
  }
}

function resetFilters() {
  filters.keyword = ''
  filters.categories = []
  filters.source = ''
  filters.resolution = ''
  filters.promotion = ''
  filters.status = 'all'
  filters.undownloadedOnly = false
}

const desktopColumns = [
  { colKey: 'category', title: '类型', width: 70 },
  { colKey: 'title', title: '标题 / Name', minWidth: 420 },
  { colKey: 'meta', title: '互动 / 时间', width: 180 },
  { colKey: 'size', title: '体积', width: 110, align: 'right' },
  { colKey: 'seeders', title: 'S', width: 70, align: 'right' },
  { colKey: 'leechers', title: 'L', width: 70, align: 'right' },
  { colKey: 'completed', title: 'C', width: 90, align: 'right' },
  { colKey: 'uploader', title: '发布者', width: 120 },
]

const mobileColumns = [
  { colKey: 'title', title: '标题', minWidth: 270 },
  { colKey: 'size', title: '体积', width: 110, align: 'right' },
  { colKey: 'seeders', title: 'S', width: 70, align: 'right' },
]

const tableColumns = computed(() => (isMobile.value ? mobileColumns : desktopColumns))

const filteredTorrents = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase()

  return torrents.value.filter((item) => {
    if (keyword) {
      const matchedKeyword =
        item.title.toLowerCase().includes(keyword) ||
        item.subtitle.toLowerCase().includes(keyword) ||
        item.imdbId.toLowerCase().includes(keyword)
      if (!matchedKeyword) {
        return false
      }
    }

    if (filters.categories.length > 0 && !filters.categories.includes(item.categoryId)) {
      return false
    }

    if (filters.source && item.source !== filters.source) {
      return false
    }

    if (filters.resolution && item.resolution !== filters.resolution) {
      return false
    }

    if (filters.promotion && item.promotion !== filters.promotion) {
      return false
    }

    if (filters.status !== 'all' && item.status !== filters.status) {
      return false
    }

    if (filters.undownloadedOnly && item.downloaded) {
      return false
    }

    return true
  })
})

function promotionLabel(promotion: PromotionKey): string {
  if (promotion === 'free') {
    return 'Free'
  }
  if (promotion === '2xfree') {
    return '2xFree'
  }
  if (promotion === '50down') {
    return '50% Down'
  }
  return 'Normal'
}

function tagTheme(tag: string): 'primary' | 'warning' | 'danger' | 'success' {
  if (tag.includes('Free')) {
    return 'success'
  }
  if (tag.includes('中字')) {
    return 'primary'
  }
  if (tag.includes('置顶')) {
    return 'warning'
  }
  return 'danger'
}

function parseIntLike(value: string | number | undefined): number {
  if (typeof value === 'number') {
    return Number.isFinite(value) ? value : 0
  }
  if (typeof value === 'string' && value.trim() !== '') {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : 0
  }
  return 0
}

function normalizeText(value: string): string {
  return value.trim().toLowerCase()
}

function guessCategoryIcon(name: string): string {
  const text = normalizeText(name)
  if (text.includes('movie') || text.includes('电影')) {
    return '🎬'
  }
  if (text.includes('series') || text.includes('剧') || text.includes('tv')) {
    return '📺'
  }
  if (text.includes('doc') || text.includes('纪录')) {
    return '🎞️'
  }
  if (text.includes('anime') || text.includes('动画')) {
    return '🧩'
  }
  if (text.includes('music') || text.includes('音乐')) {
    return '🎵'
  }
  return '📦'
}

function flattenCategories(categories: Category[] | undefined): Category[] {
  const result: Category[] = []
  const stack = [...(categories || [])]
  while (stack.length > 0) {
    const current = stack.shift()
    if (!current) {
      continue
    }
    result.push(current)
    if (current.children && current.children.length > 0) {
      stack.push(...current.children)
    }
  }
  return result
}

async function loadCategoryOptions() {
  try {
    const response = await CategoryService.ListCategories({})
    const categories = flattenCategories(response.categories)
    if (categories.length === 0) {
      categoryOptions.value = defaultCategories
      return
    }
    categoryOptions.value = categories.map((item) => ({
      id: item.id || item.name || Math.random().toString(36).slice(2),
      name: item.name || item.slug || '未分类',
      icon: item.icon || guessCategoryIcon(item.name || item.slug || ''),
    }))
  } catch {
    categoryOptions.value = defaultCategories
  }
}

function getMetaValue(metadata: ResourceMeta[] | undefined, keys: string[]): string {
  if (!metadata || metadata.length === 0) {
    return ''
  }

  const keySet = new Set(keys.map((key) => key.toLowerCase()))
  for (const item of metadata) {
    const key = (item.key || '').toLowerCase()
    if (keySet.has(key)) {
      return item.value || ''
    }
  }
  return ''
}

function parsePromotion(resource: Resource): PromotionKey {
  if (resource.isFree && resource.doubleUpload) {
    return '2xfree'
  }
  if (resource.isFree) {
    return 'free'
  }
  return 'normal'
}

function getCategoryInfo(categoryId: string | undefined): CategoryOption {
  if (categoryId) {
    const matched = categoryOptions.value.find((item) => item.id === categoryId)
    if (matched) {
      return matched
    }
  }

  const fallbackName = categoryId || '未分类'
  return {
    id: categoryId || 'unknown',
    name: fallbackName,
    icon: guessCategoryIcon(fallbackName),
  }
}

function mapResourceToTorrent(resource: Resource, torrent: TorrentInfo | undefined): TorrentRecord {
  const category = getCategoryInfo(resource.categoryId)
  const metadata = resource.metadata || []
  const tags = (resource.tags || []).map((item) => item.name || '').filter(Boolean)

  if (resource.isFree) {
    tags.push('Free')
  }
  if (resource.doubleUpload) {
    tags.push('2xUP')
  }
  if (resource.isSticky) {
    tags.push('置顶')
  }

  const torrentSize = parseIntLike(torrent?.size ?? resource.torrentSize)
  const imdbRating = parseIntLike(getMetaValue(metadata, ['imdb_rating', 'imdbRate', 'rating']))
  const promotion = parsePromotion(resource)
  const status = torrent?.isActive === false || resource.status === 2 ? 'dead' : 'alive'

  return {
    id: resource.id || torrent?.id || '',
    categoryId: category.id,
    categoryLabel: category.name,
    categoryIcon: category.icon,
    title: resource.title || torrent?.name || 'Untitled',
    subtitle: resource.subtitle || '',
    tags,
    comments: resource.commentCount || 0,
    files: torrent?.fileCount || 0,
    addedAt: resource.createdAt || '',
    size: torrentSize,
    seeders: torrent?.seedCount || resource.seedCount || 0,
    leechers: torrent?.leechCount || resource.leechCount || 0,
    completed: torrent?.snatchCount || resource.snatchCount || 0,
    uploader: resource.uploaderName || resource.uploaderId || 'Anonymous',
    anonymous: false,
    imdbId: getMetaValue(metadata, ['imdb', 'imdb_id', 'imdbId']),
    imdbRating,
    poster: getMetaValue(metadata, ['poster', 'cover', 'image']),
    source: getMetaValue(metadata, ['source', 'media', 'format']),
    resolution: getMetaValue(metadata, ['resolution', 'res']),
    promotion,
    status,
    downloaded: false,
  }
}

async function loadTorrentDetailMap(resources: Resource[]): Promise<Map<string, TorrentInfo>> {
  const torrentIds = Array.from(new Set(resources.map((item) => item.torrentId).filter(Boolean))) as string[]
  const detailMap = new Map<string, TorrentInfo>()

  for (let i = 0; i < torrentIds.length; i += 10) {
    const chunk = torrentIds.slice(i, i + 10)
    const settled = await Promise.allSettled(
      chunk.map(async (torrentId) => {
        const response = await TorrentService.GetTorrent({ id: torrentId })
        if (response.torrent) {
          detailMap.set(torrentId, response.torrent)
        }
      }),
    )

    for (const item of settled) {
      if (item.status === 'rejected') {
        continue
      }
    }
  }

  return detailMap
}

async function fetchTorrents() {
  loading.value = true
  loadError.value = ''

  try {
    const response = await ResourceService.ListResources({
      keyword: filters.keyword.trim(),
      status: 0,
      page: 1,
      pageSize: 200,
      sortBy: 'created_at',
      sortOrder: 'desc',
    })

    const resources = response.resources || []

    if (resources.length === 0) {
      torrents.value = []
      return
    }

    const detailMap = await loadTorrentDetailMap(resources)
    torrents.value = resources.map((resource) => mapResourceToTorrent(resource, detailMap.get(resource.torrentId || '')))

    // 如果后端分类为空，至少保证筛选按钮可用
    if (categoryOptions.value.length === 0) {
      const fallbackFromData = Array.from(
        new Set(
          torrents.value
            .map((item) => item.categoryId)
            .filter((categoryId) => categoryId && categoryId !== 'unknown'),
        ),
      ).map((id) => ({ id, name: id, icon: guessCategoryIcon(id) }))

      categoryOptions.value = fallbackFromData.length > 0 ? fallbackFromData : defaultCategories
    }
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '加载失败'
    loadError.value = message
    MessagePlugin.error(`种子列表加载失败: ${message}`)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  handleResize()
  window.addEventListener('resize', handleResize)
  await loadCategoryOptions()
  await fetchTorrents()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<template>
  <div class="space-y-3">
    <t-card title="种子筛选" size="small">
      <div class="grid gap-3 xl:grid-cols-[1.4fr_1fr_1fr_1fr_auto_auto]">
        <t-input v-model="filters.keyword" clearable placeholder="标题搜索 / IMDB ID" @enter="fetchTorrents">
          <template #prefix-icon>
            <search-icon />
          </template>
        </t-input>

        <t-select v-model="filters.source" clearable placeholder="格式/媒介">
          <t-option v-for="option in sourceOptions" :key="option" :label="option" :value="option" />
        </t-select>

        <t-select v-model="filters.resolution" clearable placeholder="分辨率">
          <t-option v-for="option in resolutionOptions" :key="option" :label="option" :value="option" />
        </t-select>

        <t-select v-model="filters.promotion" clearable placeholder="促销状态">
          <t-option label="Free" value="free" />
          <t-option label="2xFree" value="2xfree" />
          <t-option label="50% Down" value="50down" />
        </t-select>

        <t-button theme="primary" :loading="loading" @click="fetchTorrents">查询</t-button>
        <t-button variant="outline" @click="resetFilters">重置</t-button>
      </div>

      <div class="mt-3 flex flex-wrap items-center gap-2">
        <button
          v-for="item in categoryOptions"
          :key="item.id"
          class="category-btn"
          :class="filters.categories.includes(item.id) ? 'active' : ''"
          @click="toggleCategory(item.id)"
        >
          <span>{{ item.icon }}</span>
          <span>{{ item.name }}</span>
        </button>

        <div class="ml-auto flex items-center gap-4 text-xs text-[var(--muted-text)]">
          <t-checkbox v-model="filters.undownloadedOnly">仅显示未下载</t-checkbox>
          <t-radio-group v-model="filters.status" variant="default-filled">
            <t-radio-button value="all">全部</t-radio-button>
            <t-radio-button value="alive">活种</t-radio-button>
            <t-radio-button value="dead">死种</t-radio-button>
          </t-radio-group>
        </div>
      </div>

      <p v-if="loadError" class="mt-3 text-xs text-rose-600">{{ loadError }}</p>
    </t-card>

    <t-card size="small">
      <template #title>
        <div class="flex items-center gap-2">
          <span>种子列表</span>
          <t-tag size="small" theme="primary" variant="light">{{ filteredTorrents.length }}</t-tag>
        </div>
      </template>

      <t-table
        row-key="id"
        size="small"
        table-layout="auto"
        :data="filteredTorrents"
        :columns="tableColumns"
        :loading="loading"
        :max-height="680"
        bordered
        hover
      >
        <template #category="{ row }">
          <span class="text-lg" :title="row.categoryLabel">{{ row.categoryIcon }}</span>
        </template>

        <template #title="{ row }">
          <t-popup trigger="hover" placement="right" :show-arrow="true">
            <template #content>
              <div class="w-58 p-2">
                <img
                  v-if="row.poster"
                  class="h-42 w-full rounded-md object-cover"
                  :src="row.poster"
                  :alt="row.title"
                />
                <div v-else class="grid h-42 place-items-center rounded-md bg-[var(--soft-bg)] text-xs text-[var(--muted-text)]">
                  No Poster
                </div>
                <p class="mt-2 text-sm font-600">{{ row.title }}</p>
                <p class="text-xs text-[var(--muted-text)]">IMDB: {{ row.imdbRating || '-' }} / 10</p>
                <p class="text-xs text-[var(--muted-text)]">{{ row.imdbId || 'N/A' }}</p>
              </div>
            </template>

            <div class="cursor-pointer space-y-1">
              <div class="flex flex-wrap items-center gap-1">
                <t-tag
                  v-for="tag in [...row.tags, promotionLabel(row.promotion)]"
                  :key="tag"
                  size="small"
                  variant="light"
                  :theme="tagTheme(tag)"
                >
                  {{ tag }}
                </t-tag>
              </div>
              <p class="text-sm font-600 leading-5 text-[var(--app-text)]">{{ row.title }}</p>
              <p class="text-xs leading-4 text-[var(--muted-text)]">{{ row.subtitle }}</p>
            </div>
          </t-popup>
        </template>

        <template #meta="{ row }">
          <div class="text-xs leading-5 text-[var(--muted-text)]">
            <p>💬 {{ row.comments }} · 📁 {{ row.files }}</p>
            <p>{{ formatDateTime(row.addedAt) }}</p>
          </div>
        </template>

        <template #size="{ row }">
          <span class="font-600">{{ formatBytes(row.size) }}</span>
        </template>

        <template #seeders="{ row }">
          <span class="font-700 text-emerald-600">{{ row.seeders }}</span>
        </template>

        <template #leechers="{ row }">
          <span class="font-700 text-rose-600">{{ row.leechers }}</span>
        </template>

        <template #completed="{ row }">
          <span class="font-600">{{ row.completed }}</span>
        </template>

        <template #uploader="{ row }">
          <span>{{ row.anonymous ? 'Anonymous' : row.uploader }}</span>
        </template>
      </t-table>
    </t-card>
  </div>
</template>

<style scoped>
.category-btn {
  border: 1px solid var(--app-border);
  border-radius: 7px;
  padding: 4px 9px;
  font-size: 12px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--panel-bg);
  color: var(--muted-text);
  transition: all 0.2s ease;
}

.category-btn:hover {
  border-color: #0ea5e9;
  color: #0369a1;
}

.category-btn.active {
  border-color: #0284c7;
  background: #e0f2fe;
  color: #075985;
}

:global(html[theme-mode='dark']) .category-btn.active {
  background: rgba(14, 165, 233, 0.2);
  color: #bae6fd;
}
</style>
