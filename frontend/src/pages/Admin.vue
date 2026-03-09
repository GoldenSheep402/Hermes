<script setup lang="ts">
import { reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'

type AdminSection = 'site' | 'invite' | 'tracker' | 'torrent' | 'user'

interface SectionItem {
  key: AdminSection
  label: string
}

const sections: SectionItem[] = [
  { key: 'site', label: '站点基础设置' },
  { key: 'invite', label: '注册与邀请控制' },
  { key: 'tracker', label: 'Tracker 参数' },
  { key: 'torrent', label: '种子管理与审核' },
  { key: 'user', label: '用户管理' },
]

const activeSection = ref<AdminSection>('site')

const siteSetting = reactive({
  siteName: 'Hermes PT Station',
  maintenanceMode: false,
  marquee: '欢迎来到 Hermes PT，理性下载，保持做种。',
})

const inviteSetting = reactive({
  openRegistration: false,
  inviteOnly: true,
  globalMessage: '本周末开启 2x Free 活动。',
})

const trackerSetting = reactive({
  announceInterval: 1800,
  globalFreeleech: false,
  freeleechCountdown: 72,
  bonusFormula: 'seed_time * 1.15 + torrent_size_factor',
})

const torrentSetting = reactive({
  clearDeadSeedDays: 30,
  pendingReports: 7,
})

const userSetting = reactive({
  banTarget: '',
  resetTrafficTarget: '',
  promoteTarget: '',
})

function saveSection() {
  MessagePlugin.success('配置已保存（示例）')
}

function runAction(label: string) {
  MessagePlugin.success(`${label} 已执行（示例）`)
}

function updateSiteMaintenance(value: boolean) {
  siteSetting.maintenanceMode = value
}

function updateOpenRegistration(value: boolean) {
  inviteSetting.openRegistration = value
}

function updateInviteOnly(value: boolean) {
  inviteSetting.inviteOnly = value
}

function updateGlobalFreeleech(value: boolean) {
  trackerSetting.globalFreeleech = value
}
</script>

<template>
  <div class="grid gap-4 lg:grid-cols-[220px_1fr]">
    <aside class="rounded-lg border border-[var(--app-border)] bg-[var(--panel-bg)] p-2">
      <button
        v-for="item in sections"
        :key="item.key"
        class="admin-menu-btn"
        :class="item.key === activeSection ? 'active' : ''"
        @click="activeSection = item.key"
      >
        {{ item.label }}
      </button>
    </aside>

    <section>
      <t-card v-if="activeSection === 'site'" title="站点基础设置" size="small">
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
          <t-button theme="primary" @click="saveSection">保存设置</t-button>
        </template>
      </t-card>

      <t-card v-else-if="activeSection === 'invite'" title="注册与邀请控制" size="small">
        <div class="form-grid">
          <label>
            <span>自由注册</span>
            <t-switch :value="inviteSetting.openRegistration" @change="updateOpenRegistration" />
          </label>
          <label>
            <span>仅邀请码注册</span>
            <t-switch :value="inviteSetting.inviteOnly" @change="updateInviteOnly" />
          </label>
          <label class="lg:col-span-2">
            <span>全局系统消息</span>
            <t-textarea v-model="inviteSetting.globalMessage" :autosize="{ minRows: 3, maxRows: 4 }" />
          </label>
        </div>
        <template #footer>
          <div class="flex gap-2">
            <t-button theme="primary" @click="saveSection">保存设置</t-button>
            <t-button variant="outline" @click="runAction('发送全局消息')">发送全局消息</t-button>
          </div>
        </template>
      </t-card>

      <t-card v-else-if="activeSection === 'tracker'" title="Tracker 及业务参数" size="small">
        <div class="form-grid">
          <label>
            <span>Announce 最小间隔（秒）</span>
            <t-input-number v-model="trackerSetting.announceInterval" :min="60" :step="60" />
          </label>
          <label>
            <span>全站 Freeleech</span>
            <t-switch :value="trackerSetting.globalFreeleech" @change="updateGlobalFreeleech" />
          </label>
          <label>
            <span>活动倒计时（小时）</span>
            <t-input-number v-model="trackerSetting.freeleechCountdown" :min="1" :max="720" />
          </label>
          <label>
            <span>魔力值公式</span>
            <t-input v-model="trackerSetting.bonusFormula" />
          </label>
        </div>
        <template #footer>
          <t-button theme="primary" @click="saveSection">保存设置</t-button>
        </template>
      </t-card>

      <t-card v-else-if="activeSection === 'torrent'" title="种子管理与审核" size="small">
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

      <t-card v-else title="用户管理" size="small">
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
    </section>
  </div>
</template>

<style scoped>
.admin-menu-btn {
  width: 100%;
  text-align: left;
  border: 1px solid transparent;
  background: transparent;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 13px;
  margin-bottom: 4px;
  color: var(--muted-text);
  transition: all 0.2s ease;
}

.admin-menu-btn:hover {
  background: var(--soft-bg);
  color: var(--app-text);
}

.admin-menu-btn.active {
  background: #0ea5e9;
  color: #fff;
}

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
  display: grid;
  gap: 6px;
  font-size: 12px;
  color: var(--muted-text);
}
</style>
