<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>

<template>
  <t-card :title="title" size="small">
    <p class="text-sm text-[var(--muted-text)]">{{ description }}</p>
  </t-card>
</template>
