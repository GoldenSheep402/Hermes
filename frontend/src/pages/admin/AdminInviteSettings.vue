<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Setting } from '@/lib/proto/system/v1/system.pb'
import { boolToValue, loadSettingsMap, parseBoolean, pickChangedSettings, saveSettings, settingKeys } from './system-settings'

const inviteSetting = reactive({
  openRegistration: false,
  inviteOnly: true,
  emailVerificationRequired: false,
  smtpEnable: false,
  globalMessage: '本周末开启 2x Free 活动。',
})

const settingsLoading = ref(false)
const settingsSaving = ref(false)
const settingsMap = ref<Record<string, string>>({})

function syncSettingsToState(map: Record<string, string>) {
  inviteSetting.openRegistration = parseBoolean(map[settingKeys.openRegistration], inviteSetting.openRegistration)
  inviteSetting.inviteOnly = parseBoolean(map[settingKeys.inviteOnly], inviteSetting.inviteOnly)
  inviteSetting.emailVerificationRequired = parseBoolean(
    map[settingKeys.emailVerificationRequired],
    inviteSetting.emailVerificationRequired,
  )
  inviteSetting.smtpEnable = parseBoolean(map[settingKeys.smtpEnable], inviteSetting.smtpEnable)
  inviteSetting.globalMessage = map[settingKeys.globalMessage] || inviteSetting.globalMessage
}

function buildPayload(): Setting[] {
  return [
    {
      key: settingKeys.openRegistration,
      value: boolToValue(inviteSetting.openRegistration),
      type: 'bool',
      desc: '开启自由注册',
    },
    { key: settingKeys.inviteOnly, value: boolToValue(inviteSetting.inviteOnly), type: 'bool', desc: '仅邀请码注册' },
    {
      key: settingKeys.emailVerificationRequired,
      value: boolToValue(inviteSetting.emailVerificationRequired),
      type: 'bool',
      desc: '注册是否需要邮箱验证码',
    },
    {
      key: settingKeys.smtpEnable,
      value: boolToValue(inviteSetting.smtpEnable),
      type: 'bool',
      desc: 'SMTP 发信开关',
    },
    { key: settingKeys.globalMessage, value: inviteSetting.globalMessage.trim(), type: 'string', desc: '全局系统消息' },
  ]
}

async function loadSectionSettings() {
  settingsLoading.value = true
  try {
    const map = await loadSettingsMap()
    settingsMap.value = map
    syncSettingsToState(map)
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '加载系统设置失败'
    MessagePlugin.error(message)
  } finally {
    settingsLoading.value = false
  }
}

async function saveSectionSettings() {
  if (inviteSetting.emailVerificationRequired && !inviteSetting.smtpEnable) {
    MessagePlugin.warning('已开启注册邮箱验证，请先开启 SMTP 发信开关')
    return
  }

  const payload = pickChangedSettings(settingsMap.value, buildPayload())
  if (payload.length === 0) {
    MessagePlugin.info('配置未变化，无需保存')
    return
  }

  settingsSaving.value = true
  try {
    await saveSettings(payload)

    const merged = { ...settingsMap.value }
    for (const item of payload) {
      if (!item.key) {
        continue
      }
      merged[item.key] = item.value || ''
    }

    settingsMap.value = merged
    syncSettingsToState(merged)
    MessagePlugin.success('配置已保存')
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '保存配置失败'
    MessagePlugin.error(message)
  } finally {
    settingsSaving.value = false
  }
}

function runAction(label: string) {
  MessagePlugin.success(`${label} 已执行（示例）`)
}

function updateOpenRegistration(value: boolean) {
  inviteSetting.openRegistration = value
}

function updateInviteOnly(value: boolean) {
  inviteSetting.inviteOnly = value
}

function updateEmailVerificationRequired(value: boolean) {
  inviteSetting.emailVerificationRequired = value
}

function updateSmtpEnable(value: boolean) {
  inviteSetting.smtpEnable = value
}

onMounted(() => {
  void loadSectionSettings()
})
</script>

<template>
  <t-card title="注册与邀请控制" size="small" :loading="settingsLoading" class="h-full">
    <div class="form-grid">
      <label>
        <span>自由注册</span>
        <t-switch :value="inviteSetting.openRegistration" @change="updateOpenRegistration" />
      </label>
      <label>
        <span>仅邀请码注册</span>
        <t-switch :value="inviteSetting.inviteOnly" @change="updateInviteOnly" />
      </label>
      <label>
        <span>注册邮箱验证</span>
        <t-switch
          :value="inviteSetting.emailVerificationRequired"
          @change="updateEmailVerificationRequired"
        />
      </label>
      <label>
        <span>SMTP 发信开关</span>
        <t-switch :value="inviteSetting.smtpEnable" @change="updateSmtpEnable" />
      </label>
      <label class="lg:col-span-2">
        <span>全局系统消息</span>
        <t-textarea v-model="inviteSetting.globalMessage" :autosize="{ minRows: 3, maxRows: 4 }" />
      </label>
    </div>

    <template #footer>
      <div class="flex gap-2">
        <t-button theme="primary" :loading="settingsSaving" @click="saveSectionSettings">保存设置</t-button>
        <t-button variant="outline" @click="runAction('发送全局消息')">发送全局消息</t-button>
      </div>
    </template>
  </t-card>
</template>

<style scoped>
.form-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  gap: 12px;
}

@media (min-width: 1024px) {
  .form-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

.form-grid label {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 6px;
  font-size: 12px;
  color: var(--muted-text);
}

.form-grid label :deep(.t-switch) {
  align-self: flex-start;
  width: auto;
}
</style>
