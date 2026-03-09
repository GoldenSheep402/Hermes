<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/store'
import { ADMIN_SECTIONS } from './sections'

const authStore = useAuthStore()

const visibleSections = computed(() =>
  ADMIN_SECTIONS.filter((item) => authStore.hasPermission(item.requiredPermission)),
)
</script>

<template>
  <t-card title="管理后台" size="small" class="h-full">
    <p class="mb-4 text-sm text-[var(--muted-text)]">
      各管理模块已拆分为独立页面，选择一个分区进入。
    </p>

    <div class="grid gap-3 md:grid-cols-2 xl:grid-cols-3">
      <router-link
        v-for="item in visibleSections"
        :key="item.key"
        :to="{ name: item.routeName }"
        class="entry-card"
      >
        <p class="entry-title">{{ item.label }}</p>
        <p class="entry-desc">{{ item.description }}</p>
      </router-link>
    </div>

    <p v-if="visibleSections.length === 0" class="text-sm text-rose-500">
      当前账号没有可访问的后台分区。
    </p>
  </t-card>
</template>

<style scoped>
.entry-card {
  border: 1px solid var(--app-border);
  background: var(--soft-bg);
  border-radius: 10px;
  padding: 12px;
  display: grid;
  gap: 6px;
  transition: all 0.2s ease;
}

.entry-card:hover {
  border-color: #0ea5e9;
  transform: translateY(-1px);
}

.entry-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--app-text);
}

.entry-desc {
  margin: 0;
  font-size: 12px;
  color: var(--muted-text);
}
</style>
