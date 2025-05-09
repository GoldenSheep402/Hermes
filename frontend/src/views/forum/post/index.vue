<script lang="ts" setup>
import Editor from '@/components/Editor/index.vue'
import { useAppStore } from '@/store'
import { Message } from '@arco-design/web-vue'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

interface Post {
  id: number
  title: string
  author: string
  createTime: string
  content: string
  replies: Reply[]
}

interface Reply {
  id: number
  author: string
  createTime: string
  content: string
}

const router = useRouter()
const appStore = useAppStore()
const theme = computed(() => appStore.theme as 'light' | 'dark')

// 模拟数据
const post = ref<Post>({
  id: 1,
  title: 'Stable Diffusion XL 3.0 发布了！',
  author: 'AI助手',
  createTime: '2024-03-20 10:00:00',
  content: `
    <h2>重大更新</h2>
    <p>Stable Diffusion XL 3.0 今天正式发布，带来了以下新特性：</p>
    <ul>
      <li>更高质量的图像生成</li>
      <li>更快的生成速度</li>
      <li>更好的文本理解能力</li>
    </ul>
    <p>欢迎大家下载体验！</p>
  `,
  replies: [
    {
      id: 1,
      author: '用户A',
      createTime: '2024-03-20 11:00:00',
      content: '太棒了！期待已久！',
    },
    {
      id: 2,
      author: '用户B',
      createTime: '2024-03-20 12:00:00',
      content: '生成速度确实快了很多，效果也很不错。',
    },
  ],
})

const replyContent = ref('')

function submitReply() {
  if (!replyContent.value.trim()) {
    Message.warning('请输入回复内容')
    return
  }

  // TODO: 实现回复逻辑
  post.value.replies.push({
    id: post.value.replies.length + 1,
    author: 'current_user',
    createTime: new Date().toLocaleString(),
    content: replyContent.value,
  })

  replyContent.value = ''
  Message.success('回复成功')
}

function goBack() {
  router.push({ name: 'ForumList' })
}
</script>

<template>
  <div class="forum-post">
    <a-button class="mb-4" @click="goBack">
      <template #icon>
        <icon-left />
      </template>
      返回论坛
    </a-button>

    <!-- 帖子内容 -->
    <a-card class="mb-4">
      <template #title>
        <div class="text-xl text-[var(--color-text-1)] font-bold">
          {{ post.title }}
        </div>
      </template>
      <template #extra>
        <div class="flex items-center gap-2">
          <span class="text-[var(--color-text-3)]">作者：{{ post.author }}</span>
          <span class="text-[var(--color-text-3)]">发布于：{{ post.createTime }}</span>
        </div>
      </template>
      <div class="prose max-w-none" v-html="post.content" />
    </a-card>

    <!-- 回复列表 -->
    <a-card class="mb-4">
      <template #title>
        <div class="text-lg text-[var(--color-text-1)] font-bold">
          回复 ({{ post.replies.length }})
        </div>
      </template>
      <div v-for="reply in post.replies" :key="reply.id" class="mb-4 border-b border-[var(--color-border)] pb-4 last:border-0">
        <div class="mb-2 flex items-center gap-2">
          <span class="text-[var(--color-text-1)] font-medium">{{ reply.author }}</span>
          <span class="text-sm text-[var(--color-text-3)]">{{ reply.createTime }}</span>
        </div>
        <div class="prose max-w-none" v-html="reply.content" />
      </div>
    </a-card>

    <!-- 回复框 -->
    <a-card>
      <template #title>
        <div class="text-lg text-[var(--color-text-1)] font-bold">
          发表回复
        </div>
      </template>
      <div class="mb-4">
        <Editor
          v-model="replyContent"
          placeholder="请输入回复内容..."
          min-height="200px"
          :theme="theme"
        />
      </div>
      <div class="flex justify-end">
        <a-button type="primary" @click="submitReply">
          发表回复
        </a-button>
      </div>
    </a-card>
  </div>
</template>

<style scoped>
.forum-post {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

:deep(.prose) {
  max-width: none;
  line-height: 1.6;
}

:deep(.prose h1),
:deep(.prose h2),
:deep(.prose h3),
:deep(.prose h4),
:deep(.prose h5),
:deep(.prose h6) {
  margin: 1em 0 0.5em;
  font-weight: bold;
  color: var(--color-text-1);
}

:deep(.prose p) {
  margin: 0.5em 0;
}

:deep(.prose ul),
:deep(.prose ol) {
  padding-left: 1.5em;
  margin: 0.5em 0;
}

:deep(.prose a) {
  color: var(--color-primary);
  text-decoration: none;
}

:deep(.prose a:hover) {
  color: var(--color-primary-hover);
}
</style>
