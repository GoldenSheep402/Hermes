<script setup lang="ts">
import type { TorrentMetaData } from '@/lib/proto/torrent/v1/torrent.pb'
import RichTextViewer from '@/components/RichTextViewer/index.vue'
import { TorrentService } from '@/services/grpc.ts'
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const torrentMetadata = ref<TorrentMetaData[]>()

const sortedMetadata = computed(() => {
  if (!torrentMetadata.value)
    return []
  return [...torrentMetadata.value].sort((a, b) => {
    const orderA = a.order || 0
    const orderB = b.order || 0
    return orderA - orderB
  })
})

async function getTorrentDetail() {
  const torrentId = route.query.id as string
  try {
    const res = await TorrentService.GetTorrentV1({ id: torrentId })
    torrentMetadata.value = res.metadata
  }
  catch (error) {
    console.error('获取种子详情失败:', error)
  }
}

onMounted(() => {
  getTorrentDetail()
})
</script>

<template>
  <div class="p-5">
    <div class="rounded-sm bg-[--color-bg-2] p-5 shadow-sm">
      <div class="mb-8 flex items-center p-0.5 text-24px text-[--color-text-1] font-600 leading-[1.4]">
        <icon-info-circle class="mr-2 text-[var(--color-primary-6)]" />
        种子详情
      </div>

      <a-card class="mb-6" :bordered="true" :header-style="{ padding: '20px' }">
        <template #title>
          <div class="flex items-center text-[var(--color-text-1)]">
            <icon-list class="mr-2 text-[var(--color-primary-6)]" />
            基本信息
          </div>
        </template>
        <div class="p-4">
          <div v-for="(item, index) in sortedMetadata" :key="index" class="detail-item">
            <div class="detail-header">
              <div class="detail-label">
                <icon-tag class="mr-2 text-[var(--color-primary-6)]" />
                {{ item.key }}
              </div>
              <div class="detail-type">
                {{ item.type }}
              </div>
            </div>
            <div class="detail-content">
              <RichTextViewer
                v-if="item.key === '简介'"
                :content="item.value || ''"
                class="value-rich-text"
              />
              <div v-else class="value-text">
                {{ item.value || '-' }}
              </div>
            </div>
            <div v-if="item.description" class="detail-description">
              <icon-info-circle class="mr-1 text-[var(--color-text-3)]" />
              {{ item.description }}
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<style scoped lang="less">
.detail-item {
  margin-bottom: 16px;
  padding: 12px;
  border-bottom: 1px solid var(--color-border-2);

  &:last-child {
    margin-bottom: 0;
    border-bottom: none;
  }
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.detail-label {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 500;
  color: var(--color-text-1);
}

.detail-type {
  font-size: 12px;
  padding: 4px 8px;
  background-color: var(--color-fill-2);
  color: var(--color-text-3);
}

.detail-content {
  margin-bottom: 8px;
  padding: 8px;
}

.detail-description {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: var(--color-text-3);
  padding: 4px 8px;
}

.value-input {
  width: 100%;
  background-color: var(--color-bg-1);
}

.value-switch {
  margin: 8px 0;
}

.value-rich-text {
  margin: 12px 0;
}

.value-text {
  color: var(--color-text-1);
  word-break: break-all;
  line-height: 1.5;
}
</style>
