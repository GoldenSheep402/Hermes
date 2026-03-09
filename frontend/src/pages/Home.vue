<script setup lang="ts">
import { useAuthStore } from '@/store'
import { formatBytes, formatNumber, formatRatio } from '@/utils/format'

const authStore = useAuthStore()
</script>

<template>
  <section class="grid gap-4 md:grid-cols-3">
    <t-card title="账号总览" size="small">
      <p class="line-item"><span>用户名</span><strong>{{ authStore.profile.username }}</strong></p>
      <p class="line-item"><span>用户等级</span><strong>{{ authStore.profile.className }}</strong></p>
      <p class="line-item"><span>信誉积分</span><strong>{{ formatNumber(authStore.profile.reputation) }}</strong></p>
    </t-card>

    <t-card title="流量状态" size="small">
      <p class="line-item"><span>上传总量</span><strong class="text-emerald-600">{{ formatBytes(authStore.profile.uploadBytes) }}</strong></p>
      <p class="line-item"><span>下载总量</span><strong class="text-sky-600">{{ formatBytes(authStore.profile.downloadBytes) }}</strong></p>
      <p class="line-item"><span>分享率</span><strong :class="authStore.ratio < 1 ? 'text-rose-600' : 'text-emerald-600'">{{ formatRatio(authStore.ratio) }}</strong></p>
    </t-card>

    <t-card title="快速入口" size="small">
      <div class="grid gap-2">
        <router-link class="quick-entry" to="/torrents">进入种子列表</router-link>
        <router-link class="quick-entry" to="/forums">进入论坛</router-link>
        <router-link class="quick-entry" to="/faq">查看 FAQ</router-link>
      </div>
    </t-card>
  </section>
</template>

<style scoped>
.line-item {
  margin: 0;
  padding: 8px 0;
  border-bottom: 1px dashed var(--app-border);
  display: flex;
  justify-content: space-between;
  gap: 12px;
  font-size: 13px;
}

.line-item:last-child {
  border-bottom: none;
}

.quick-entry {
  border: 1px solid var(--app-border);
  border-radius: 8px;
  padding: 8px 10px;
  font-size: 13px;
  color: var(--app-text);
  background: var(--soft-bg);
}

.quick-entry:hover {
  border-color: #0ea5e9;
  color: #0369a1;
}
</style>
