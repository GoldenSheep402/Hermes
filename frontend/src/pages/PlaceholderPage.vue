<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>

<template>
  <section class="w-full">
    <t-card :title="title" size="small" class="w-full">
      <p class="text-sm text-[var(--muted-text)]">{{ description }}</p>
    </t-card>
  </section>
</template>
