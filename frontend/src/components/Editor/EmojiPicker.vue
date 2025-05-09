<!-- eslint-disable vue/v-on-event-hyphenation -->
<script lang="ts" setup>
import { ref } from 'vue'

const props = defineProps<{
  onSelect: (emoji: string) => void
}>()

const visible = ref(false)
const buttonRef = ref<HTMLElement | null>(null)

// 常用表情列表
const emojis = [
  '😀',
  '😃',
  '😄',
  '😁',
  '😆',
  '😅',
  '😂',
  '🤣',
  '😊',
  '😇',
  '🙂',
  '🙃',
  '😉',
  '😌',
  '😍',
  '🥰',
  '😘',
  '😗',
  '😙',
  '😚',
  '😋',
  '😛',
  '😝',
  '😜',
  '🤪',
  '🤨',
  '🧐',
  '🤓',
  '😎',
  '🤩',
  '🥳',
  '😏',
  '😒',
  '😞',
  '😔',
  '😟',
  '😕',
  '🙁',
  '☹️',
  '😣',
  '😖',
  '😫',
  '😩',
  '🥺',
  '😢',
  '😭',
  '😤',
  '😠',
  '😡',
  '🤬',
  '🤯',
  '😳',
  '🥵',
  '🥶',
  '😱',
  '😨',
  '👍',
  '👎',
  '👊',
  '✌️',
  '🤞',
  '🤝',
  '🙏',
  '💪',
]

function handleSelect(emoji: string) {
  props.onSelect(emoji)
  visible.value = false
}
</script>

<template>
  <div class="emoji-picker">
    <a-popover
      v-model:popup-visible="visible"
      trigger="click"
      position="bl"
    >
      <a-button
        ref="buttonRef"
        size="small"
        @click="visible = !visible"
      >
        <template #icon>
          <icon-emoji />
        </template>
      </a-button>
      <template #content>
        <div class="emoji-grid">
          <button
            v-for="emoji in emojis"
            :key="emoji"
            class="emoji-button"
            @click="handleSelect(emoji)"
          >
            {{ emoji }}
          </button>
        </div>
      </template>
    </a-popover>
  </div>
</template>

<style scoped>
.emoji-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 4px;
  padding: 8px;
  max-width: 320px;
}

.emoji-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  font-size: 20px;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.emoji-button:hover {
  background-color: var(--color-fill-2);
}
</style>
