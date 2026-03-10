<script setup lang="ts">
import { useAuthStore } from '@/store'
import { formatBytes, formatNumber, formatRatio } from '@/utils/format'

const authStore = useAuthStore()
</script>

<template>
  <section class="home-dashboard grid w-full items-start gap-4 xl:grid-cols-[minmax(0,1fr)_340px]">
    <div class="grid gap-4">
      <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
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

        <t-card title="站内状态" size="small">
          <p class="line-item"><span>未读站内信</span><strong>{{ formatNumber(authStore.profile.inboxUnread) }}</strong></p>
          <p class="line-item"><span>邀请码</span><strong>{{ formatNumber(authStore.profile.inviteCount) }}</strong></p>
          <p class="line-item"><span>可用魔力</span><strong>{{ formatNumber(authStore.profile.bonusPoints) }}</strong></p>
        </t-card>
      </div>

      <t-card title="新手提示" size="small">
        <ul class="guide-list">
          <li>发布前先选择正确分类，系统会自动加载该分类元数据字段。</li>
          <li>下载后保持做种，建议优先保证分享率不低于 1.0。</li>
          <li>资源描述尽量包含来源、规格和完整信息，便于检索。</li>
        </ul>
      </t-card>
    </div>

    <t-card title="快速入口" size="small" class="quick-panel">
      <div class="grid gap-2">
        <router-link class="quick-entry" to="/torrents">进入种子列表</router-link>
        <router-link class="quick-entry" to="/upload">发布资源</router-link>
        <router-link class="quick-entry" to="/forums">进入论坛</router-link>
        <router-link class="quick-entry" to="/faq">查看 FAQ</router-link>
      </div>
    </t-card>
  </section>
</template>

<style scoped>
.home-dashboard {
  align-content: start;
}

.quick-panel {
  position: sticky;
  top: 0;
}

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

.guide-list {
  margin: 0;
  padding-left: 18px;
  display: grid;
  gap: 8px;
  font-size: 12px;
  color: var(--muted-text);
}

@media (max-width: 1279px) {
  .quick-panel {
    position: static;
  }
}
</style>
