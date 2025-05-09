<!-- eslint-disable vue/v-on-event-hyphenation -->
<script lang="ts" setup>
import { Message } from '@arco-design/web-vue'
import { ref } from 'vue'

const props = defineProps<{
  onSuccess: (url: string) => void
}>()

const uploadRef = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

function handleClick() {
  uploadRef.value?.click()
}

async function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file)
    return

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    Message.error('请选择图片文件')
    return
  }

  // 检查文件大小（限制为 5MB）
  if (file.size > 5 * 1024 * 1024) {
    Message.error('图片大小不能超过 5MB')
    return
  }

  uploading.value = true
  try {
    // TODO: 实现实际的图片上传逻辑
    // 这里暂时使用 FileReader 在本地预览
    const reader = new FileReader()
    reader.onload = () => {
      props.onSuccess(reader.result as string)
      uploading.value = false
    }
    reader.readAsDataURL(file)
  }
  catch (error) {
    console.error(error)
    Message.error('图片上传失败')
    uploading.value = false
  }

  // 清除 input 的值，允许重复上传相同的文件
  target.value = ''
}
</script>

<template>
  <div class="image-uploader">
    <input
      ref="uploadRef"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleFileChange"
    >
    <a-button
      size="small"
      :loading="uploading"
      @click="handleClick"
    >
      <template #icon>
        <icon-upload />
      </template>
      上传图片
    </a-button>
  </div>
</template>
