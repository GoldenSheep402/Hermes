<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Setting } from '@/lib/proto/system/v1/system.pb'
import { boolToValue, loadSettingsMap, parseBoolean, pickChangedSettings, saveSettings, settingKeys } from './system-settings'

const siteSetting = reactive({
  siteName: 'Hermes PT Station',
  maintenanceMode: false,
  marquee: '欢迎来到 Hermes PT，理性下载，保持做种。',
})

const settingsLoading = ref(false)
const settingsSaving = ref(false)
const settingsMap = ref<Record<string, string>>({})

function syncSettingsToState(map: Record<string, string>) {
  siteSetting.siteName = map[settingKeys.siteName] || siteSetting.siteName
  siteSetting.maintenanceMode = parseBoolean(map[settingKeys.maintenanceMode], siteSetting.maintenanceMode)
  siteSetting.marquee = map[settingKeys.marquee] || siteSetting.marquee
}

function buildPayload(): Setting[] {
  return [
    { key: settingKeys.siteName, value: siteSetting.siteName.trim(), type: 'string', desc: '站点名称' },
    { key: settingKeys.maintenanceMode, value: boolToValue(siteSetting.maintenanceMode), type: 'bool', desc: '维护模式' },
    { key: settingKeys.marquee, value: siteSetting.marquee.trim(), type: 'string', desc: '公告跑马灯' },
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

function updateSiteMaintenance(value: boolean) {
  siteSetting.maintenanceMode = value
}

onMounted(() => {
  void loadSectionSettings()
})
</script>

<template>
  <t-card title="站点基础设置" size="small" :loading="settingsLoading" class="h-full">
    <div class="form-grid">
      <label>
        <span>站点名称</span>
        <t-input v-model="siteSetting.siteName" clearable />
      </label>
      <label>
        <span>维护模式</span>
        <t-switch :value="siteSetting.maintenanceMode" @change="updateSiteMaintenance" />
      </label>
      <label class="lg:col-span-2">
        <span>公告跑马灯</span>
        <t-textarea v-model="siteSetting.marquee" :autosize="{ minRows: 3, maxRows: 4 }" />
      </label>
    </div>

    <template #footer>
      <div class="flex items-center gap-2">
        <t-button theme="primary" :loading="settingsSaving" @click="saveSectionSettings">保存设置</t-button>
        <t-button variant="outline" :loading="settingsLoading" @click="loadSectionSettings">刷新设置</t-button>
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
