<script setup lang="ts">
import { reactive } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'

const userSetting = reactive({
  banTarget: '',
  resetTrafficTarget: '',
  promoteTarget: '',
})

function runAction(label: string) {
  MessagePlugin.success(`${label} 已执行（示例）`)
}
</script>

<template>
  <t-card title="用户管理" size="small" class="h-full">
    <div class="form-grid">
      <label>
        <span>封禁账号</span>
        <t-input v-model="userSetting.banTarget" placeholder="用户名 / UID" />
      </label>
      <label>
        <span>重置流量</span>
        <t-input v-model="userSetting.resetTrafficTarget" placeholder="用户名 / UID" />
      </label>
      <label>
        <span>等级提升</span>
        <t-input v-model="userSetting.promoteTarget" placeholder="用户名 / UID" />
      </label>
    </div>

    <template #footer>
      <div class="flex flex-wrap gap-2">
        <t-button theme="danger" @click="runAction('账号封禁')">账号封禁</t-button>
        <t-button theme="warning" @click="runAction('流量重置')">流量重置</t-button>
        <t-button theme="primary" @click="runAction('等级提升')">等级提升</t-button>
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
