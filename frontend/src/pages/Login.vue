<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { useAuthStore } from '@/store'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const form = reactive({
  email: '',
  password: '',
})

async function handleLogin() {
  if (!form.email.trim() || !form.password.trim()) {
    MessagePlugin.warning('请输入邮箱和密码')
    return
  }

  loading.value = true
  try {
    await authStore.login({ email: form.email, password: form.password })
    MessagePlugin.success('登录成功')
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/home'
    await router.replace(redirect)
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '登录失败'
    MessagePlugin.error(message)
  } finally {
    loading.value = false
  }
}

</script>

<template>
  <div class="login-page">
    <div class="login-background"></div>
    <section class="login-panel">
      <div class="mb-6">
        <p class="text-xs text-sky-600 font-600 tracking-[0.14em]">PRIVATE TRACKER</p>
        <h1 class="mt-2 text-3xl font-700">Hermes PT Station</h1>
        <p class="mt-2 text-sm text-[var(--muted-text)]">登录后访问种子列表与用户流量面板。</p>
      </div>

      <div class="space-y-4">
        <div>
          <p class="mb-2 text-xs text-[var(--muted-text)]">邮箱</p>
          <t-input v-model="form.email" clearable placeholder="输入登录邮箱" size="large" />
        </div>

        <div>
          <p class="mb-2 text-xs text-[var(--muted-text)]">密码</p>
          <t-input v-model="form.password" type="password" placeholder="输入密码" size="large" @enter="handleLogin" />
        </div>

        <t-button class="w-full" theme="primary" size="large" :loading="loading" @click="handleLogin">登录</t-button>
      </div>
    </section>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(120deg, #e0f2fe 0%, #f0f9ff 45%, #eef2ff 100%);
}

.login-background {
  position: absolute;
  inset: 0;
  background-image:
    radial-gradient(circle at 20% 20%, rgba(14, 165, 233, 0.18), transparent 35%),
    radial-gradient(circle at 80% 70%, rgba(2, 132, 199, 0.15), transparent 42%);
}

.login-panel {
  width: min(460px, calc(100vw - 28px));
  border: 1px solid var(--app-border);
  border-radius: 16px;
  background: color-mix(in srgb, var(--panel-bg) 94%, white);
  box-shadow: 0 18px 50px rgba(2, 132, 199, 0.12);
  padding: 28px;
  position: relative;
  z-index: 1;
}

:global(html[theme-mode='dark']) .login-page {
  background: linear-gradient(120deg, #082f49 0%, #0f172a 48%, #1e293b 100%);
}
</style>
