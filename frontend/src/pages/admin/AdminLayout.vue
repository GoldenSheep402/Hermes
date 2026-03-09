<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/store'
import { ADMIN_SECTIONS } from './sections'

const route = useRoute()
const authStore = useAuthStore()

const visibleSections = computed(() =>
  ADMIN_SECTIONS.filter((item) => authStore.hasPermission(item.requiredPermission)),
)

function isSectionActive(routeName: string): boolean {
  return route.name === routeName
}
</script>

<template>
  <div class="grid h-full min-h-full w-full flex-1 items-stretch gap-4 overflow-hidden lg:grid-cols-[220px_minmax(0,1fr)]">
    <aside class="admin-aside h-full min-h-0 self-stretch overflow-y-auto rounded-lg border border-[var(--app-border)] bg-[var(--panel-bg)] p-2">
      <router-link
        to="/admin"
        class="admin-menu-btn"
        :class="route.name === 'AdminOverview' ? 'active' : ''"
      >
        管理后台首页
      </router-link>

      <router-link
        v-for="item in visibleSections"
        :key="item.key"
        :to="{ name: item.routeName }"
        class="admin-menu-btn"
        :class="isSectionActive(item.routeName) ? 'active' : ''"
      >
        {{ item.label }}
      </router-link>

      <p v-if="visibleSections.length === 0" class="px-2 py-3 text-xs text-rose-500">
        当前账号没有后台分区权限
      </p>
    </aside>

    <section class="admin-content flex h-full min-h-0 max-h-full self-stretch flex-col overflow-hidden">
      <div class="admin-view-host h-full min-h-0 min-w-0 flex-1">
        <router-view />
      </div>
    </section>
  </div>
</template>

<style scoped>
.admin-menu-btn {
  width: 100%;
  text-align: left;
  border: 1px solid transparent;
  background: transparent;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 13px;
  margin-bottom: 4px;
  color: var(--muted-text);
  transition: all 0.2s ease;
  display: block;
}

.admin-menu-btn:hover {
  background: var(--soft-bg);
  color: var(--app-text);
}

.admin-menu-btn.active {
  background: #0ea5e9;
  color: #fff;
}

.admin-aside {
  display: flex;
  flex-direction: column;
}

.admin-content :deep(.t-card) {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.admin-content :deep(.t-card__body) {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
}

.admin-view-host {
  display: flex;
  height: 100%;
  min-height: 0;
  width: 100%;
  flex: 1 1 auto;
}
</style>
