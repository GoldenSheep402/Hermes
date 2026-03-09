<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Setting } from '@/lib/proto/system/v1/system.pb'
import {
  boolToValue,
  clampNumber,
  loadSettingsMap,
  parseBoolean,
  parseNumber,
  pickChangedSettings,
  saveSettings,
  settingKeys,
} from './system-settings'

const trackerSetting = reactive({
  announceInterval: 1800,
  flushInterval: 60,
  flushBatchSize: 200,
  globalFreeleech: false,
  freeleechCountdown: 0,
  bonusFormula: 'sqrt(uploaded)',
  trackerList: '',
})

const settingsLoading = ref(false)
const settingsSaving = ref(false)
const settingsMap = ref<Record<string, string>>({})

function syncSettingsToState(map: Record<string, string>) {
  trackerSetting.announceInterval = parseNumber(map[settingKeys.announceInterval], trackerSetting.announceInterval)
  trackerSetting.flushInterval = parseNumber(map[settingKeys.flushInterval], trackerSetting.flushInterval)
  trackerSetting.flushBatchSize = parseNumber(map[settingKeys.flushBatchSize], trackerSetting.flushBatchSize)
  trackerSetting.globalFreeleech = parseBoolean(map[settingKeys.globalFreeleech], trackerSetting.globalFreeleech)
  trackerSetting.freeleechCountdown = parseNumber(map[settingKeys.freeleechCountdown], trackerSetting.freeleechCountdown)
  trackerSetting.bonusFormula = map[settingKeys.bonusFormula] || trackerSetting.bonusFormula
  trackerSetting.trackerList = map[settingKeys.trackerList] || trackerSetting.trackerList
}

function buildPayload(): Setting[] {
  return [
    {
      key: settingKeys.announceInterval,
      value: String(Math.max(60, Number(trackerSetting.announceInterval) || 1800)),
      type: 'int',
      desc: 'announce 最小间隔（秒）',
    },
    {
      key: settingKeys.flushInterval,
      value: String(clampNumber(Number(trackerSetting.flushInterval), 1, 600, 60)),
      type: 'int',
      desc: 'tracker 流量 flush 间隔（秒）',
    },
    {
      key: settingKeys.flushBatchSize,
      value: String(clampNumber(Number(trackerSetting.flushBatchSize), 50, 5000, 200)),
      type: 'int',
      desc: 'tracker 流量 flush 批大小',
    },
    {
      key: settingKeys.globalFreeleech,
      value: boolToValue(trackerSetting.globalFreeleech),
      type: 'bool',
      desc: '全站 freeleech 开关',
    },
    {
      key: settingKeys.freeleechCountdown,
      value: String(Math.max(0, Number(trackerSetting.freeleechCountdown) || 0)),
      type: 'int',
      desc: '活动倒计时（小时）',
    },
    { key: settingKeys.bonusFormula, value: trackerSetting.bonusFormula.trim(), type: 'string', desc: '魔力值公式' },
    {
      key: settingKeys.trackerList,
      value: trackerSetting.trackerList
        .split('\n')
        .map((line) => line.trim())
        .filter((line) => line.length > 0)
        .join('\n'),
      type: 'string',
      desc: 'Tracker 列表，一行一个地址',
    },
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

function updateGlobalFreeleech(value: boolean) {
  trackerSetting.globalFreeleech = value
}

onMounted(() => {
  void loadSectionSettings()
})
</script>

<template>
  <t-card title="Tracker 及业务参数" size="small" :loading="settingsLoading" class="h-full">
    <div class="form-grid">
      <label>
        <span>Announce 最小间隔（秒）</span>
        <t-input-number v-model="trackerSetting.announceInterval" :min="60" :step="60" />
      </label>
      <label>
        <span>流量 Flush 间隔（秒）</span>
        <t-input-number v-model="trackerSetting.flushInterval" :min="1" :max="600" :step="1" />
      </label>
      <label>
        <span>流量 Flush 批大小</span>
        <t-input-number v-model="trackerSetting.flushBatchSize" :min="50" :max="5000" :step="50" />
      </label>
      <label>
        <span>全站 Freeleech</span>
        <t-switch :value="trackerSetting.globalFreeleech" @change="updateGlobalFreeleech" />
      </label>
      <label>
        <span>活动倒计时（小时）</span>
        <t-input-number v-model="trackerSetting.freeleechCountdown" :min="0" :max="720" />
      </label>
      <label>
        <span>魔力值公式</span>
        <t-input v-model="trackerSetting.bonusFormula" />
      </label>
      <label class="lg:col-span-2">
        <span>Tracker URL 列表（系统设置 list）</span>
        <t-textarea
          v-model="trackerSetting.trackerList"
          :autosize="{ minRows: 4, maxRows: 8 }"
          placeholder="一行一个地址，例如：&#10;https://tracker.example.com/announce&#10;https://backup.example.com/announce"
        />
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
