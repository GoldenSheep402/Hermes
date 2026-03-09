<template>
  <div class="min-h-screen grid place-items-center bg-[var(--app-bg)] px-4">
    <t-card title="404 - Not Found" size="small" class="w-full max-w-lg">
      <p class="text-sm text-[var(--muted-text)]">页面不存在，可能已被移动或删除。</p>
      <div class="mt-4">
        <router-link to="/home">
          <t-button theme="primary">返回首页</t-button>
        </router-link>
      </div>
    </t-card>
  </div>
</template>
