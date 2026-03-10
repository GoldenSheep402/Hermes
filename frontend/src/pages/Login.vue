<script setup lang="ts">
import { computed, onBeforeUnmount, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { AuthService } from '@/services/grpc'
import { useAuthStore } from '@/store'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

type AuthMode = 'login' | 'register'

const mode = ref<AuthMode>('login')
const loginLoading = ref(false)
const registerLoading = ref(false)
const sendCodeLoading = ref(false)
const sendCodeCountdown = ref(0)
let sendCodeTimer: number | undefined

const loginForm = reactive({
  email: '',
  password: '',
})

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  emailToken: '',
})

const sendCodeText = computed(() =>
  sendCodeCountdown.value > 0 ? `${sendCodeCountdown.value}s 后重试` : '发送验证码',
)

function handleModeChange(value: string) {
  mode.value = value === 'register' ? 'register' : 'login'
}

function trim(value: string): string {
  return value.trim()
}

function isEmailValid(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

function parseErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error && error.message ? error.message : fallback
}

function stopSendCodeTimer() {
  if (typeof sendCodeTimer !== 'undefined') {
    window.clearInterval(sendCodeTimer)
    sendCodeTimer = undefined
  }
}

function startSendCodeCountdown(seconds = 60) {
  stopSendCodeTimer()
  sendCodeCountdown.value = seconds
  sendCodeTimer = window.setInterval(() => {
    sendCodeCountdown.value -= 1
    if (sendCodeCountdown.value <= 0) {
      sendCodeCountdown.value = 0
      stopSendCodeTimer()
    }
  }, 1000)
}

async function handleLogin() {
  const email = trim(loginForm.email)
  const password = trim(loginForm.password)
  if (!email || !password) {
    MessagePlugin.warning('请输入邮箱和密码')
    return
  }

  loginLoading.value = true
  try {
    await authStore.login({ email, password })
    MessagePlugin.success('登录成功')
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/home'
    await router.replace(redirect)
  } catch (error: unknown) {
    MessagePlugin.error(parseErrorMessage(error, '登录失败'))
  } finally {
    loginLoading.value = false
  }
}

async function handleSendRegisterCode() {
  const email = trim(registerForm.email)
  if (!email) {
    MessagePlugin.warning('请输入注册邮箱')
    return
  }
  if (!isEmailValid(email)) {
    MessagePlugin.warning('邮箱格式不正确')
    return
  }
  if (sendCodeLoading.value || sendCodeCountdown.value > 0) {
    return
  }

  sendCodeLoading.value = true
  try {
    await AuthService.RegisterSendEmail({ email })
    MessagePlugin.success('验证码已发送（若站点开启邮箱验证）')
    startSendCodeCountdown(60)
  } catch (error: unknown) {
    MessagePlugin.error(parseErrorMessage(error, '验证码发送失败'))
  } finally {
    sendCodeLoading.value = false
  }
}

async function handleRegister() {
  const username = trim(registerForm.username)
  const email = trim(registerForm.email)
  const password = trim(registerForm.password)
  const confirmPassword = trim(registerForm.confirmPassword)
  const emailToken = trim(registerForm.emailToken)

  if (!username || !email || !password || !confirmPassword) {
    MessagePlugin.warning('请完整填写注册信息')
    return
  }
  if (!isEmailValid(email)) {
    MessagePlugin.warning('邮箱格式不正确')
    return
  }
  if (password.length < 6) {
    MessagePlugin.warning('密码至少 6 位')
    return
  }
  if (password !== confirmPassword) {
    MessagePlugin.warning('两次输入的密码不一致')
    return
  }

  registerLoading.value = true
  try {
    await AuthService.RegisterWithEmail({
      username,
      email,
      password,
      emailToken,
    })
    MessagePlugin.success('注册成功，请登录')

    loginForm.email = email
    loginForm.password = ''
    registerForm.password = ''
    registerForm.confirmPassword = ''
    registerForm.emailToken = ''
    mode.value = 'login'
  } catch (error: unknown) {
    MessagePlugin.error(parseErrorMessage(error, '注册失败'))
  } finally {
    registerLoading.value = false
  }
}

onBeforeUnmount(() => {
  stopSendCodeTimer()
})
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

      <t-tabs :value="mode" @change="handleModeChange">
        <t-tab-panel value="login" label="登录">
          <div class="mt-4 space-y-4">
            <div>
              <p class="mb-2 text-xs text-[var(--muted-text)]">邮箱</p>
              <t-input v-model="loginForm.email" clearable placeholder="输入登录邮箱" size="large" />
            </div>

            <div>
              <p class="mb-2 text-xs text-[var(--muted-text)]">密码</p>
              <t-input
                v-model="loginForm.password"
                type="password"
                placeholder="输入密码"
                size="large"
                @enter="handleLogin"
              />
            </div>

            <t-button class="w-full" theme="primary" size="large" :loading="loginLoading" @click="handleLogin">
              登录
            </t-button>
          </div>
        </t-tab-panel>

        <t-tab-panel value="register" label="注册">
          <div class="mt-4 space-y-4">
            <div>
              <p class="mb-2 text-xs text-[var(--muted-text)]">用户名</p>
              <t-input v-model="registerForm.username" clearable placeholder="输入用户名" size="large" />
            </div>

            <div>
              <p class="mb-2 text-xs text-[var(--muted-text)]">邮箱</p>
              <t-input v-model="registerForm.email" clearable placeholder="输入注册邮箱" size="large" />
            </div>

            <div>
              <p class="mb-2 text-xs text-[var(--muted-text)]">密码</p>
              <t-input v-model="registerForm.password" type="password" placeholder="至少 6 位" size="large" />
            </div>

            <div>
              <p class="mb-2 text-xs text-[var(--muted-text)]">确认密码</p>
              <t-input
                v-model="registerForm.confirmPassword"
                type="password"
                placeholder="再次输入密码"
                size="large"
                @enter="handleRegister"
              />
            </div>

            <div>
              <p class="mb-2 text-xs text-[var(--muted-text)]">邮箱验证码（按站点设置可选）</p>
              <div class="flex gap-2">
                <t-input v-model="registerForm.emailToken" clearable placeholder="输入验证码" size="large" />
                <t-button
                  variant="outline"
                  size="large"
                  :loading="sendCodeLoading"
                  :disabled="sendCodeCountdown > 0"
                  @click="handleSendRegisterCode"
                >
                  {{ sendCodeText }}
                </t-button>
              </div>
            </div>

            <t-button class="w-full" theme="primary" size="large" :loading="registerLoading" @click="handleRegister">
              注册
            </t-button>
          </div>
        </t-tab-panel>
      </t-tabs>
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
  width: min(520px, calc(100vw - 28px));
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
