<script setup lang="ts">
import type { TorrentMetaData } from '@/lib/proto/torrent/v1/torrent.pb'
import { TorrentService } from '@/services/grpc.ts'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const torrentMetadata = ref<TorrentMetaData[]>()

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
    <div class="mt-5">
      <a-typography-title :heading="4">
        种子详情
      </a-typography-title>
      <a-table :data="torrentMetadata">
        <template #columns>
          <a-table-column title="键" data-index="key" :width="150" />
          <a-table-column title="描述" data-index="description" :width="300" />
          <a-table-column title="类型" data-index="type" :width="100" />
          <a-table-column title="值" :width="200">
            <template #cell="{ record }">
              <a-input v-if="record.type === 'string'" v-model="record.value" />
              <a-input-number v-else-if="record.type === 'number'" v-model="record.value" />
              <a-switch v-else-if="record.type === 'switch'" v-model="record.value" />
              <a-input v-if="record.type === 'textarea'" v-model="record.value" />
              <div v-else>{{ record.value }}</div>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>
  </div>
</template>
