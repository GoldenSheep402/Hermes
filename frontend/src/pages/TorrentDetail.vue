<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Resource } from '@/lib/proto/resource/v1/resource.pb'
import type { TorrentFile, TorrentInfo } from '@/lib/proto/torrent/v1/torrent.pb'
import { ResourceService, TorrentService } from '@/services/grpc'
import type { TorrentFileRow as FileRow } from '@/types/torrent'
import { formatBytes, formatDateTime } from '@/utils/format'
import { downloadTorrentFile } from '@/utils/torrent'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const downloading = ref(false)
const loadError = ref('')

const resource = ref<Resource | null>(null)
const torrent = ref<TorrentInfo | null>(null)
const fileRows = ref<FileRow[]>([])

const fileColumns = [
  { colKey: 'path', title: '文件路径', minWidth: 520 },
  { colKey: 'size', title: '大小', width: 120, align: 'right' },
]

const resourceId = computed(() => String(route.params.id || '').trim())

const pageTitle = computed(() => resource.value?.title || torrent.value?.name || '种子详情')

function normalizeText(value: string): string {
  return value.trim().toLowerCase()
}

function parseBooleanLike(value: string): boolean {
  const text = normalizeText(value)
  return text === '1' || text === 'true' || text === 'yes' || text === 'y' || text === 'on'
}

function isResourceAnonymous(item: Resource | null): boolean {
  if (!item) {
    return false
  }

  const uploaderText = normalizeText(item.uploaderName || '')
  if (uploaderText === 'anonymous' || uploaderText === '匿名') {
    return true
  }

  const metadata = item.metadata || []
  const anonymousKeys = new Set(['anonymous', 'is_anonymous', 'anon'])
  for (const entry of metadata) {
    const key = normalizeText(entry.key || '')
    if (!anonymousKeys.has(key)) {
      continue
    }
    if (parseBooleanLike(entry.value || '')) {
      return true
    }
  }

  const tags = item.tags || []
  for (const tag of tags) {
    const name = normalizeText(tag.name || '')
    if (name === 'anonymous' || name === '匿名') {
      return true
    }
  }

  return false
}

const uploaderDisplayName = computed(() =>
  isResourceAnonymous(resource.value) ? '匿名' : resource.value?.uploaderName || '未知用户',
)

const promotionTags = computed(() => {
  const tags: string[] = []
  if (resource.value?.isFree) {
    tags.push('Free')
  }
  if (resource.value?.doubleUpload) {
    tags.push('2xUP')
  }
  if (resource.value?.isSticky) {
    tags.push('置顶')
  }
  return tags
})

const canDownload = computed(() => Boolean(resource.value?.torrentId || torrent.value?.id))

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

function formatDateText(value: unknown): string {
  if (value instanceof Date) {
    return formatDateTime(value.toISOString())
  }
  if (typeof value === 'string') {
    return formatDateTime(value)
  }
  return '-'
}

function mapTorrentFiles(files: TorrentFile[] | undefined): FileRow[] {
  return (files || []).map((item, index) => ({
    id: item.id || `${index}`,
    path: item.path || '未知路径',
    size: parseIntLike(item.size),
  }))
}

interface StructureNode {
  name: string
  kind: 'dir' | 'file'
  size: number
  children: Map<string, StructureNode>
}

function ensureDirNode(container: Map<string, StructureNode>, name: string): StructureNode {
  const existed = container.get(name)
  if (existed) {
    return existed
  }

  const node: StructureNode = {
    name,
    kind: 'dir',
    size: 0,
    children: new Map<string, StructureNode>(),
  }
  container.set(name, node)
  return node
}

function upsertFileNode(container: Map<string, StructureNode>, name: string, size: number) {
  const existed = container.get(name)
  if (existed) {
    existed.kind = 'file'
    existed.size = size
    existed.children.clear()
    return
  }

  container.set(name, {
    name,
    kind: 'file',
    size,
    children: new Map<string, StructureNode>(),
  })
}

function sortedNodes(nodes: Map<string, StructureNode>): StructureNode[] {
  return Array.from(nodes.values()).sort((left, right) => {
    if (left.kind !== right.kind) {
      return left.kind === 'dir' ? -1 : 1
    }
    return left.name.localeCompare(right.name, 'zh-CN', { sensitivity: 'base' })
  })
}

function buildStructureLines(files: FileRow[], rootName: string): string[] {
  const root = new Map<string, StructureNode>()

  for (const file of files) {
    const segments = file.path
      .split(/[\\/]+/)
      .map((item) => item.trim())
      .filter(Boolean)

    if (segments.length === 0) {
      continue
    }

    let current = root
    for (let i = 0; i < segments.length; i += 1) {
      const segment = segments[i]
      const isLast = i === segments.length - 1

      if (isLast) {
        upsertFileNode(current, segment, file.size)
      } else {
        const dirNode = ensureDirNode(current, segment)
        current = dirNode.children
      }
    }
  }

  const lines: string[] = [`📦 ${rootName || '种子结构'}`]

  function walk(nodes: Map<string, StructureNode>, prefix: string) {
    const items = sortedNodes(nodes)
    for (let i = 0; i < items.length; i += 1) {
      const item = items[i]
      const isLast = i === items.length - 1
      const connector = isLast ? '└─ ' : '├─ '
      if (item.kind === 'dir') {
        lines.push(`${prefix}${connector}📁 ${item.name}/`)
        walk(item.children, `${prefix}${isLast ? '   ' : '│  '}`)
      } else {
        lines.push(`${prefix}${connector}📄 ${item.name} (${formatBytes(item.size)})`)
      }
    }
  }

  walk(root, '')
  return lines
}

const structureLines = computed(() => {
  if (fileRows.value.length > 0) {
    return buildStructureLines(fileRows.value, pageTitle.value)
  }

  if (torrent.value?.name) {
    const size = parseIntLike(torrent.value.size)
    const sizeText = size > 0 ? ` (${formatBytes(size)})` : ''
    return [`📦 ${pageTitle.value}`, `└─ 📄 ${torrent.value.name}${sizeText}`]
  }

  return ['📦 种子结构', '└─ 暂无文件结构数据']
})

async function loadDetail() {
  if (!resourceId.value) {
    loadError.value = '缺少资源 ID'
    resource.value = null
    torrent.value = null
    fileRows.value = []
    return
  }

  loading.value = true
  loadError.value = ''
  try {
    const resourceResp = await ResourceService.GetResource({ id: resourceId.value })
    if (!resourceResp.resource) {
      throw new Error('资源不存在或已删除')
    }

    resource.value = resourceResp.resource
    const torrentId = resourceResp.resource.torrentId || ''
    if (!torrentId) {
      torrent.value = null
      fileRows.value = []
      return
    }

    const [torrentResp, filesResp] = await Promise.allSettled([
      TorrentService.GetTorrent({ id: torrentId }),
      TorrentService.ListTorrentFiles({ torrentId }),
    ])

    if (torrentResp.status === 'fulfilled') {
      torrent.value = torrentResp.value.torrent || null
    } else {
      torrent.value = null
    }

    if (filesResp.status === 'fulfilled') {
      fileRows.value = mapTorrentFiles(filesResp.value.files)
    } else {
      fileRows.value = []
    }
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '加载失败'
    loadError.value = message
    MessagePlugin.error(`种子详情加载失败: ${message}`)
  } finally {
    loading.value = false
  }
}

async function handleDownload() {
  const torrentId = resource.value?.torrentId || torrent.value?.id || ''
  if (!torrentId) {
    MessagePlugin.warning('当前资源缺少可下载的 torrentId')
    return
  }

  downloading.value = true
  try {
    await downloadTorrentFile(torrentId, resource.value?.title || torrent.value?.name || torrentId)
    MessagePlugin.success('种子下载已开始')
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '下载失败'
    MessagePlugin.error(message)
  } finally {
    downloading.value = false
  }
}

watch(
  () => route.params.id,
  () => {
    void loadDetail()
  },
  { immediate: true },
)
</script>

<template>
  <div class="space-y-3">
    <t-card :loading="loading" size="small">
      <template #title>
        <div class="flex items-center gap-2">
          <span>{{ pageTitle }}</span>
          <t-tag v-for="tag in promotionTags" :key="tag" size="small" theme="success" variant="light">
            {{ tag }}
          </t-tag>
        </div>
      </template>

      <div class="mb-3 flex flex-wrap gap-2">
        <t-button theme="primary" variant="outline" :loading="downloading" :disabled="!canDownload" @click="handleDownload">
          下载种子
        </t-button>
        <t-button variant="outline" @click="router.push('/torrents')">返回列表</t-button>
      </div>

      <p v-if="loadError" class="m-0 text-xs text-rose-600">{{ loadError }}</p>

      <div v-else class="grid grid-cols-1 gap-2 md:grid-cols-2">
        <p class="m-0 flex items-center justify-between gap-3 rounded-lg border border-dashed border-[var(--app-border)] px-2.5 py-2 text-xs">
          <span>发布者</span>
          <strong class="max-w-[60%] text-right text-[var(--app-text)]">{{ uploaderDisplayName }}</strong>
        </p>
        <p class="m-0 flex items-center justify-between gap-3 rounded-lg border border-dashed border-[var(--app-border)] px-2.5 py-2 text-xs">
          <span>发布时间</span>
          <strong class="max-w-[60%] text-right text-[var(--app-text)]">{{ formatDateText(resource?.createdAt) }}</strong>
        </p>
        <p class="m-0 flex items-center justify-between gap-3 rounded-lg border border-dashed border-[var(--app-border)] px-2.5 py-2 text-xs">
          <span>体积</span>
          <strong class="max-w-[60%] text-right text-[var(--app-text)]">{{ formatBytes(parseIntLike(torrent?.size || resource?.torrentSize)) }}</strong>
        </p>
        <p class="m-0 flex items-center justify-between gap-3 rounded-lg border border-dashed border-[var(--app-border)] px-2.5 py-2 text-xs">
          <span>做种/下载/完成</span>
          <strong class="max-w-[60%] text-right text-[var(--app-text)]">
            {{ torrent?.seedCount || 0 }} / {{ torrent?.leechCount || 0 }} / {{ torrent?.snatchCount || 0 }}
          </strong>
        </p>
      </div>

      <div class="mt-3">
        <p class="mb-1.5 text-xs text-[var(--muted-text)]">简介</p>
        <p class="m-0 whitespace-pre-wrap rounded-lg border border-dashed border-[var(--app-border)] bg-[var(--soft-bg)] p-2.5 text-[13px] leading-6">
          {{ resource?.description || '暂无简介' }}
        </p>
      </div>

      <div class="mt-3">
        <p class="mb-1.5 text-xs text-[var(--muted-text)]">种子结构</p>
        <pre class="m-0 overflow-x-auto whitespace-pre rounded-lg border border-dashed border-[var(--app-border)] bg-[var(--soft-bg)] p-2.5 text-xs leading-5 text-[var(--app-text)]">{{ structureLines.join('\n') }}</pre>
      </div>
    </t-card>

    <t-card title="文件列表" size="small">
      <t-table row-key="id" size="small" :data="fileRows" :columns="fileColumns" bordered :loading="loading">
        <template #size="{ row }">
          <span class="font-600">{{ formatBytes(row.size) }}</span>
        </template>
      </t-table>
    </t-card>
  </div>
</template>
