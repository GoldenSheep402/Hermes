<template>
  <div class="min-h-screen grid place-items-center bg-[var(--app-bg)] px-4">
    <t-card title="403 - Forbidden" size="small" class="w-full max-w-lg">
      <p class="text-sm text-[var(--muted-text)]">当前账号没有访问该页面的权限，请联系管理员。</p>
      <div class="mt-4 flex gap-2">
        <router-link to="/home">
          <t-button theme="primary">返回首页</t-button>
        </router-link>
        <router-link to="/login">
          <t-button variant="outline">重新登录</t-button>
        </router-link>
      </div>
    </t-card>
  </div>
</template>
