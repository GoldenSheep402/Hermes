<script setup lang="ts">
import type { Category } from '@/lib/proto/category/v1/category.pb.ts'
import { CategoryService } from '@/services/grpc.ts'
import { Notification } from '@arco-design/web-vue'
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const category = ref<Category>({
  id: '',
  name: '',
  description: '',
  metaData: [],
})

function fetchCategoryDetail() {
  CategoryService.GetCategory({ id: route.params.id as string }).then((res) => {
    if (res.category) {
      category.value = res.category
    }
  }).catch((err) => {
    console.error(err)
    handleNotification('error', '错误', '获取类别详情失败')
  })
}

function handleNotification(type: string, title: string, content: string) {
  switch (type) {
    case 'success':
      Notification.success({
        title,
        content,
      })
      break
    case 'error':
      Notification.error({
        title,
        content,
      })
      break
    case 'warning':
      Notification.warning({
        title,
        content,
      })
      break
    default:
      Notification.info({
        title,
        content,
      })
  }
}

function goBack() {
  router.push({ name: 'CategoryList' })
}

onMounted(() => {
  fetchCategoryDetail()
})
</script>

<template>
  <div class="p-5">
    <div class="bg-[--color-bg-2] p-5">
      <div class="mb-4">
        <a-space>
          <a-button class="cursor-pointer" @click="goBack">
            <template #icon>
              <icon-left />
            </template>
            返回列表
          </a-button>
        </a-space>
      </div>

      <div class="mb-4">
        <h2 class="mb-2 text-2xl text-[var(--color-text-1)] font-bold">
          {{ category.name }}
        </h2>
        <p class="text-[var(--color-text-3)]">
          {{ category.description }}
        </p>
      </div>

      <a-card title="元数据配置">
        <a-table :data="category.metaData">
          <template #columns>
            <a-table-column title="序号" data-index="order" :width="80" />
            <a-table-column title="键" data-index="key" :width="150" />
            <a-table-column title="类型" data-index="type" :width="100" />
            <a-table-column title="描述" data-index="description" :width="200" />
            <a-table-column title="默认值">
              <template #cell="{ record }">
                <div v-if="record.type === 'richText'" class="rich-text-content">
                  <div class="prose max-w-none" v-html="record.defaultValue" />
                </div>
                <span v-else>{{ record.defaultValue }}</span>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<style scoped lang="less">
.rich-text-content {
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid var(--color-border);
  padding: 16px;
  border-radius: 4px;
  background-color: var(--color-bg-2);
}

:deep(.prose) {
  h1, h2, h3, h4, h5, h6 {
    color: var(--color-text-1);
    margin-top: 1.5em;
    margin-bottom: 0.5em;
  }

  p {
    color: var(--color-text-2);
    margin-bottom: 1em;
  }

  ul, ol {
    color: var(--color-text-2);
    margin-bottom: 1em;
    padding-left: 1.5em;
  }

  li {
    margin-bottom: 0.5em;
  }

  code {
    background-color: var(--color-fill-2);
    padding: 0.2em 0.4em;
    border-radius: 3px;
    font-family: monospace;
  }

  pre {
    background-color: var(--color-fill-2);
    padding: 1em;
    border-radius: 4px;
    overflow-x: auto;
    margin-bottom: 1em;
  }

  blockquote {
    border-left: 4px solid var(--color-border);
    padding-left: 1em;
    margin-left: 0;
    color: var(--color-text-3);
  }

  img {
    max-width: 100%;
    height: auto;
    border-radius: 4px;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 1em;
  }

  th, td {
    border: 1px solid var(--color-border);
    padding: 0.5em;
  }

  th {
    background-color: var(--color-fill-2);
  }
}
</style>
