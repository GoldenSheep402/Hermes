<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { PermissionKeys } from '@/constants/permissions'
import { useAppStore, useAuthStore } from '@/store'
import type { NavItem } from '@/types/layout'
import { formatBytes, formatNumber, formatRatio } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

const navItems: NavItem[] = [
  { path: '/home', label: '首页 Home' },
  { path: '/torrents', label: '种子 Torrents' },
  { path: '/upload', label: '发布 Upload' },
  { path: '/forums', label: '论坛 Forums' },
  { path: '/top10', label: '排行榜 Top 10' },
  { path: '/rules', label: '规则 Rules' },
  { path: '/faq', label: 'FAQ' },
  { path: '/admin', label: '管理后台 Staff', requiredPermission: PermissionKeys.AdminPanelAccess },
]

const visibleNavItems = computed(() =>
  navItems.filter((item) => !item.requiredPermission || authStore.hasPermission(item.requiredPermission)),
)
const isAdminRoute = computed(() => route.path === '/admin' || route.path.startsWith('/admin/'))
const userInitial = computed(() => authStore.profile.username.slice(0, 1).toUpperCase() || 'U')
const ratioDanger = computed(() => authStore.ratio < 1)

function isActive(path: string): boolean {
  return route.path === path || route.path.startsWith(`${path}/`)
}

function handleThemeChange(value: boolean) {
  appStore.setDarkMode(value)
}

async function handleLogout() {
  authStore.logout()
  await router.replace('/login')
  MessagePlugin.success('已退出登录')
}
</script>

<template>
  <div
    class="flex flex-col bg-[var(--app-bg)] text-[var(--app-text)]"
    :class="isAdminRoute ? 'h-screen overflow-hidden' : 'min-h-screen'"
  >
    <header class="sticky top-0 z-40 border-b border-[var(--app-border)] bg-[var(--panel-bg)]/95 backdrop-blur">
      <section class="w-full px-3 py-3">
        <div class="grid gap-3 lg:grid-cols-[240px_1fr_280px] lg:items-center">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-lg bg-gradient-to-br from-sky-500 to-cyan-400"></div>
            <div>
              <p class="text-base font-700 leading-5">Hermes PT Station</p>
              <p class="text-xs text-[var(--muted-text)]">Nexus-style Tracker Interface</p>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2 md:grid-cols-4">
            <div class="monitor-card">
              <p class="monitor-label">⬆ Upload</p>
              <p class="monitor-value text-emerald-600">{{ formatBytes(authStore.profile.uploadBytes) }}</p>
            </div>
            <div class="monitor-card">
              <p class="monitor-label">⬇ Download</p>
              <p class="monitor-value text-sky-600">{{ formatBytes(authStore.profile.downloadBytes) }}</p>
            </div>
            <div class="monitor-card">
              <p class="monitor-label">📊 Ratio</p>
              <p class="monitor-value" :class="ratioDanger ? 'text-rose-600' : 'text-emerald-600'">
                {{ formatRatio(authStore.ratio) }}
              </p>
            </div>
            <div class="monitor-card">
              <p class="monitor-label">💎 Bonus</p>
              <p class="monitor-value text-amber-600">{{ formatNumber(authStore.profile.bonusPoints) }}</p>
            </div>
          </div>

          <div class="flex flex-wrap items-center justify-end gap-2">
            <div class="inline-flex items-center rounded-md bg-[var(--soft-bg)] px-2 py-1 text-xs">
              <t-avatar size="small">{{ userInitial }}</t-avatar>
              <span class="ml-2 font-600">{{ authStore.profile.username }}</span>
              <span class="ml-2 text-[var(--muted-text)]">{{ authStore.profile.className }}</span>
            </div>
            <t-badge :count="authStore.profile.inboxUnread" :max-count="99" show-zero>
              <t-button variant="outline" size="small">站内信</t-button>
            </t-badge>
            <t-button variant="outline" size="small">邀请码 {{ authStore.profile.inviteCount }}</t-button>
            <t-switch :value="appStore.isDarkMode" size="small" @change="handleThemeChange" />
            <t-button theme="danger" variant="outline" size="small" @click="handleLogout">登出</t-button>
          </div>
        </div>
      </section>

      <nav class="border-t border-[var(--app-border)] bg-[var(--panel-bg)]/90">
        <div class="flex w-full gap-1 overflow-x-auto px-3 py-2">
          <router-link
            v-for="item in visibleNavItems"
            :key="item.path"
            :to="item.path"
            class="rounded-md px-3 py-1.5 text-sm transition-colors"
            :class="isActive(item.path) ? 'bg-sky-500 text-white' : 'text-[var(--muted-text)] hover:bg-[var(--soft-bg)] hover:text-[var(--app-text)]'"
          >
            {{ item.label }}
          </router-link>
        </div>
      </nav>
    </header>

    <main
      class="flex min-h-0 w-full flex-1 px-3 py-4"
      :class="isAdminRoute ? 'h-0 overflow-hidden' : ''"
    >
      <div
        class="route-page-host flex min-h-0 w-full flex-1"
        :class="isAdminRoute ? 'h-full min-h-0 overflow-hidden' : ''"
      >
        <router-view />
      </div>
    </main>
  </div>
</template>

<style scoped>
.monitor-card {
  border: 1px solid var(--app-border);
  border-radius: 8px;
  background: var(--soft-bg);
  padding: 6px 8px;
}

.monitor-label {
  color: var(--muted-text);
  font-size: 11px;
  line-height: 1.1;
}

.monitor-value {
  font-size: 13px;
  font-weight: 700;
  line-height: 1.2;
  margin-top: 2px;
}
</style>
