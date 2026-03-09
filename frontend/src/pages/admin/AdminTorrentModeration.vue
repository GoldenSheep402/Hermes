<script setup lang="ts">
import { reactive } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'

const torrentSetting = reactive({
  clearDeadSeedDays: 30,
  pendingReports: 7,
})

function runAction(label: string) {
  MessagePlugin.success(`${label} 已执行（示例）`)
}
</script>

<template>
  <t-card title="种子管理与审核" size="small" class="h-full">
    <div class="form-grid">
      <label>
        <span>死种清理阈值（天）</span>
        <t-input-number v-model="torrentSetting.clearDeadSeedDays" :min="7" :max="180" />
      </label>
      <label>
        <span>待处理报错</span>
        <t-input-number v-model="torrentSetting.pendingReports" :min="0" :max="9999" />
      </label>
    </div>

    <template #footer>
      <div class="flex gap-2">
        <t-button theme="primary" @click="runAction('批量清理死种')">批量清理死种</t-button>
        <t-button variant="outline" @click="runAction('处理举报工单')">处理举报工单</t-button>
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
</style>
