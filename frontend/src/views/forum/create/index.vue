<!-- eslint-disable vue/v-on-event-hyphenation -->
<script lang="ts" setup>
import Editor from '@/components/Editor/index.vue'
import { Message } from '@arco-design/web-vue'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const formData = ref({
  category: '',
  title: '',
  content: '',
})

const rules = {
  category: [{ required: true, message: '请选择版块' }],
  title: [{ required: true, message: '请输入标题' }],
  content: [{ required: true, message: '请输入内容' }],
}

const categories = [
  { value: '1', label: '站点公告' },
  { value: '2', label: '求助专区' },
  { value: '3', label: '技术交流' },
  { value: '4', label: '资源分享' },
]

function onSubmit() {
  // TODO: 实现发帖逻辑
  Message.success('发布成功')
  router.push({ name: 'ForumList' })
}

function onCancel() {
  router.back()
}
</script>

<template>
  <div class="create-post-container p-5">
    <h2 class="mb-4 text-2xl text-[var(--color-text-1)] font-bold">
      发布帖子
    </h2>

    <a-form
      :model="formData"
      :rules="rules"
      @submit="onSubmit"
    >
      <a-form-item field="category" label="版块">
        <a-select
          v-model="formData.category"
          :options="categories"
          placeholder="请选择版块"
        />
      </a-form-item>

      <a-form-item field="title" label="标题">
        <a-input
          v-model="formData.title"
          placeholder="请输入标题"
          allow-clear
        />
      </a-form-item>

      <a-form-item field="content" label="内容">
        <Editor
          v-model="formData.content"
          placeholder="请输入帖子内容..."
          min-height="400px"
        />
      </a-form-item>

      <a-form-item>
        <a-space>
          <a-button type="primary" html-type="submit">
            发布
          </a-button>
          <a-button @click="onCancel">
            取消
          </a-button>
        </a-space>
      </a-form-item>
    </a-form>
  </div>
</template>

<style scoped>
.create-post-container {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}
</style>
