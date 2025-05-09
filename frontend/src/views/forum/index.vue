<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

interface Category {
  id: number
  name: string
  description: string
  topics: number
  posts: number
  postCount: number
  lastPost?: {
    id: number
    title: string
    author: string
    time: string
  }
}

const categories = ref<Category[]>([
  {
    id: 1,
    name: '站点公告',
    description: '重要通知和站点更新信息',
    topics: 10,
    posts: 56,
    postCount: 56,
    lastPost: {
      id: 1,
      title: '关于最新的种子规则更新',
      author: 'admin',
      time: '2025-04-09 12:30',
    },
  },
  {
    id: 2,
    name: '求助专区',
    description: '遇到问题？在这里寻求帮助',
    topics: 234,
    posts: 1205,
    postCount: 1205,
    lastPost: {
      id: 2,
      title: '下载速度很慢怎么办？',
      author: 'user123',
      time: '2025-04-09 11:45',
    },
  },
  {
    id: 3,
    name: '技术交流',
    description: '分享和讨论各种技术话题',
    topics: 456,
    posts: 3102,
    postCount: 3102,
    lastPost: {
      id: 3,
      title: 'Stable Diffusion XL 3.0 发布了！',
      author: 'tech_guru',
      time: '2025-04-09 10:15',
    },
  },
  {
    id: 4,
    name: '资源分享',
    description: '分享各种有趣的资源',
    topics: 789,
    posts: 5431,
    postCount: 5431,
    lastPost: {
      id: 4,
      title: '整理了一批AI模型资源',
      author: 'share_master',
      time: '2025-04-09 09:00',
    },
  },
])

function createPost() {
  router.push({ name: 'ForumCreate' })
}

function viewPost(postId: string) {
  router.push({
    name: 'ForumPost',
    params: { id: postId },
  })
}

function viewCategory(categoryId: number) {
  router.push({
    name: 'ForumCategory',
    params: { id: categoryId },
  })
}
</script>

<template>
  <div class="p-5">
    <div class="forum-container">
      <a-card class="mb-4">
        <template #title>
          <div class="text-xl text-[var(--color-text-1)] font-bold">
            论坛
          </div>
        </template>
        <template #extra>
          <a-button type="primary" @click="createPost">
            <template #icon>
              <icon-plus />
            </template>
            发布帖子
          </a-button>
        </template>
        <div class="forum-categories">
          <div
            v-for="category in categories"
            :key="category.id"
            class="forum-category cursor-pointer border border-[var(--color-border)] rounded-sm p-4 transition-colors hover:border-[var(--color-primary)]"
            @click="viewCategory(category.id)"
          >
            <div class="mb-2 flex items-center justify-between">
              <div class="flex items-center gap-2">
                <icon-folder class="text-[var(--color-primary)]" />
                <span class="text-lg text-[var(--color-text-1)] font-medium">{{ category.name }}</span>
              </div>
              <span class="text-[var(--color-text-3)]">{{ category.postCount }} 个帖子</span>
            </div>
            <p class="mb-2 text-[var(--color-text-3)]">
              {{ category.description }}
            </p>
            <div v-if="category.lastPost" class="text-sm">
              <span class="text-[var(--color-text-3)]">最后回复：</span>
              <a-link
                class="text-blue-500 font-medium hover:text-blue-600"
                @click.stop="viewPost(category.lastPost.id.toString())"
              >
                {{ category.lastPost.title }}
              </a-link>
              <span class="ml-2 text-[var(--color-text-3)]">
                由 {{ category.lastPost.author }} 于 {{ category.lastPost.time }}
              </span>
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<style scoped>
.forum-container {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 1rem;
}

.forum-categories {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
}

.forum-category {
  transition: all 0.3s ease;
}

.forum-category:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
