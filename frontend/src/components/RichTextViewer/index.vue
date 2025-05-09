<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  content: string
}>()

const htmlContent = computed(() => {
  return props.content
})
</script>

<template>
  <div class="rich-text-viewer" v-html="htmlContent" />
</template>

<style scoped>
.rich-text-viewer {
  min-height: 100px;
  padding: 8px;
  border: 1px solid var(--color-border);
  border-radius: 4px;
  background-color: var(--color-bg-2);
}
</style>
