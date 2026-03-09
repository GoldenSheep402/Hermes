<script setup lang="ts">
</script>

<template>
  <t-card title="管理后台入口已迁移" size="small" class="h-full">
    <p class="text-sm text-[var(--muted-text)]">后台模块已拆分为独立页面，请从左上导航进入“管理后台 Staff”。</p>
    <div class="mt-3">
      <router-link to="/admin" class="jump-link">进入管理后台</router-link>
    </div>
  </t-card>
</template>

<style scoped>
.jump-link {
  display: inline-flex;
  border: 1px solid var(--app-border);
  background: var(--soft-bg);
  color: var(--app-text);
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 13px;
}

.jump-link:hover {
  border-color: #0ea5e9;
  color: #0369a1;
}
</style>
