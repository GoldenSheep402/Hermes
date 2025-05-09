<script lang="ts" setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

interface Post {
  id: string
  title: string
  author: string
  createTime: string
  replies: number
  views: number
  lastReply?: {
    author: string
    time: string
  }
}

// 模拟数据
const categories: Record<string, string> = {
  1: '技术讨论',
  2: '资源分享',
  3: '问题求助',
  4: '其他',
}

const descriptions: Record<string, string> = {
  1: '技术相关讨论',
  2: '资源分享区',
  3: '问题求助区',
  4: '其他话题',
}

function getCategoryName(id: string) {
  return categories[id] || '未知版块'
}

function getCategoryDescription(id: string) {
  return descriptions[id] || ''
}

const categoryInfo = ref({
  id: route.params.id,
  name: getCategoryName(route.params.id as string),
  description: getCategoryDescription(route.params.id as string),
})

const posts = ref<Post[]>([
  {
    id: '1',
    title: 'Stable Diffusion XL 3.0 发布了！',
    author: 'tech_guru',
    createTime: '2025-04-09 10:15',
    replies: 23,
    views: 156,
    lastReply: {
      author: 'ai_lover',
      time: '2025-04-09 15:30',
    },
  },
  {
    id: '2',
    title: '如何提高生成图片的质量？',
    author: 'newbie123',
    createTime: '2025-04-09 09:20',
    replies: 15,
    views: 89,
    lastReply: {
      author: 'expert_user',
      time: '2025-04-09 14:45',
    },
  },
  {
    id: '3',
    title: '分享一些实用的 Prompt',
    author: 'prompt_master',
    createTime: '2025-04-09 08:30',
    replies: 45,
    views: 328,
    lastReply: {
      author: 'learner',
      time: '2025-04-09 15:20',
    },
  },
])

function createPost() {
  router.push({
    name: 'ForumCreate',
    query: { category: categoryInfo.value.id },
  })
}

function viewPost(postId: string) {
  router.push({
    name: 'ForumPost',
    params: { id: postId },
  })
}

function goBack() {
  router.push({ name: 'ForumList' })
}
</script>

<template>
  <div class="category-container p-5">
    <div class="mb-4">
      <a-space>
        <a-button class="cursor-pointer" @click="goBack">
          <template #icon>
            <icon-left />
          </template>
          返回论坛
        </a-button>
      </a-space>
    </div>

    <div class="mb-4 flex items-center justify-between">
      <div>
        <h2 class="mb-2 text-2xl text-[var(--color-text-1)] font-bold">
          {{ categoryInfo.name }}
        </h2>
        <p class="text-[var(--color-text-3)]">
          {{ categoryInfo.description }}
        </p>
      </div>
      <a-button type="primary" class="cursor-pointer" @click="createPost">
        发布帖子
      </a-button>
    </div>

    <a-card>
      <a-table :data="posts" :bordered="false" :pagination="{ pageSize: 20 }">
        <template #columns>
          <a-table-column title="主题" data-index="title">
            <template #cell="{ record }">
              <a
                class="cursor-pointer text-[var(--color-primary)] font-medium hover:text-[var(--color-primary-hover)]"
                @click="viewPost(record.id)"
              >
                {{ record.title }}
              </a>
            </template>
          </a-table-column>
          <a-table-column title="作者" data-index="author" />
          <a-table-column title="回复/查看" align="center">
            <template #cell="{ record }">
              {{ record.replies }}/{{ record.views }}
            </template>
          </a-table-column>
          <a-table-column title="最后回复">
            <template #cell="{ record }">
              <div v-if="record.lastReply">
                <div class="text-sm text-[var(--color-text-3)]">
                  by {{ record.lastReply.author }}
                </div>
                <div class="text-sm text-[var(--color-text-4)]">
                  {{ record.lastReply.time }}
                </div>
              </div>
              <span v-else>暂无</span>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<style scoped>
.category-container {
  width: 100%;
}
</style>
