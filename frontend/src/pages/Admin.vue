<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Category } from '@/lib/proto/category/v1/category.pb'
import type { Setting } from '@/lib/proto/system/v1/system.pb'
import type { PermissionKey } from '@/constants/permissions'
import { PermissionKeys } from '@/constants/permissions'
import { CategoryService, SystemService } from '@/services/grpc'
import { useAuthStore } from '@/store'

type AdminSection = 'site' | 'invite' | 'tracker' | 'torrent' | 'category' | 'user'

interface SectionItem {
  key: AdminSection
  label: string
  requiredPermission: PermissionKey
}

interface CategoryRow {
  id: string
  name: string
  slug: string
  description: string
  parentId: string
  icon: string
  sortOrder: number
  isEnabled: boolean
  depth: number
}

const authStore = useAuthStore()

const sections: SectionItem[] = [
  { key: 'site', label: '站点基础设置', requiredPermission: PermissionKeys.SiteSettingsManage },
  { key: 'invite', label: '注册与邀请控制', requiredPermission: PermissionKeys.InviteManage },
  { key: 'tracker', label: 'Tracker 参数', requiredPermission: PermissionKeys.TrackerManage },
  { key: 'torrent', label: '种子管理与审核', requiredPermission: PermissionKeys.TorrentModerate },
  { key: 'category', label: '种子类别管理', requiredPermission: PermissionKeys.CategoryManage },
  { key: 'user', label: '用户管理', requiredPermission: PermissionKeys.UserManage },
]

const visibleSections = computed(() =>
  sections.filter((section) => authStore.hasPermission(section.requiredPermission)),
)

const activeSection = ref<AdminSection>('site')

const siteSetting = reactive({
  siteName: 'Hermes PT Station',
  maintenanceMode: false,
  marquee: '欢迎来到 Hermes PT，理性下载，保持做种。',
})

const inviteSetting = reactive({
  openRegistration: false,
  inviteOnly: true,
  emailVerificationRequired: false,
  smtpEnable: false,
  globalMessage: '本周末开启 2x Free 活动。',
})

const trackerSetting = reactive({
  announceInterval: 1800,
  flushInterval: 60,
  flushBatchSize: 200,
  globalFreeleech: false,
  freeleechCountdown: 0,
  bonusFormula: 'sqrt(uploaded)',
  trackerList: '',
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

const settingsLoading = ref(false)
const settingsSaving = ref(false)
const settingsMap = ref<Record<string, string>>({})

const settingKeys = {
  siteName: 'site.name',
  maintenanceMode: 'site.maintenance_mode',
  marquee: 'site.marquee',
  openRegistration: 'invite.open_registration',
  inviteOnly: 'invite.only',
  emailVerificationRequired: 'invite.email_verification_required',
  smtpEnable: 'auth.smtp_enable',
  globalMessage: 'invite.global_message',
  announceInterval: 'tracker.announce_interval',
  flushInterval: 'tracker.flush_interval',
  flushBatchSize: 'tracker.flush_batch_size',
  globalFreeleech: 'tracker.global_freeleech',
  freeleechCountdown: 'tracker.freeleech_countdown_hours',
  bonusFormula: 'tracker.bonus_formula',
  trackerList: 'tracker.list',
} as const

const categoryRows = ref<CategoryRow[]>([])
const categoryLoading = ref(false)
const categorySaving = ref(false)
const categoryEditorVisible = ref(false)
const categoryEditorMode = ref<'create' | 'edit'>('create')
const editingCategoryId = ref('')

const categoryForm = reactive({
  name: '',
  slug: '',
  description: '',
  parentId: '',
  icon: '',
  sortOrder: 0,
  isEnabled: true,
})

const categoryColumns = [
  { colKey: 'name', title: '类别名称', minWidth: 220 },
  { colKey: 'slug', title: 'Slug', width: 180 },
  { colKey: 'sortOrder', title: '排序', width: 90, align: 'right' },
  { colKey: 'isEnabled', title: '状态', width: 90 },
  { colKey: 'actions', title: '操作', width: 160 },
]

const parentCategoryOptions = computed(() =>
  categoryRows.value
    .filter((item) => item.id !== editingCategoryId.value)
    .map((item) => ({
      label: `${'  '.repeat(item.depth)}${item.name}`,
      value: item.id,
    })),
)

watch(
  visibleSections,
  (currentSections) => {
    const hasCurrentSection = currentSections.some((item) => item.key === activeSection.value)
    if (!hasCurrentSection && currentSections.length > 0) {
      activeSection.value = currentSections[0].key
    }
  },
  { immediate: true },
)

watch(activeSection, async (section) => {
  if (section === 'category') {
    await loadCategories()
  }
})

function saveSection() {
  void saveActiveSection()
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

function updateEmailVerificationRequired(value: boolean) {
  inviteSetting.emailVerificationRequired = value
}

function updateSmtpEnable(value: boolean) {
  inviteSetting.smtpEnable = value
}

function updateGlobalFreeleech(value: boolean) {
  trackerSetting.globalFreeleech = value
}

function parseBoolean(input: string | undefined, fallback: boolean): boolean {
  if (!input) {
    return fallback
  }
  const normalized = input.trim().toLowerCase()
  if (['1', 'true', 'yes', 'on'].includes(normalized)) {
    return true
  }
  if (['0', 'false', 'no', 'off'].includes(normalized)) {
    return false
  }
  return fallback
}

function parseNumber(input: string | undefined, fallback: number): number {
  if (!input) {
    return fallback
  }
  const value = Number(input)
  return Number.isFinite(value) ? value : fallback
}

function boolToValue(input: boolean): string {
  return input ? 'true' : 'false'
}

function clampNumber(input: number, min: number, max: number, fallback: number): number {
  if (!Number.isFinite(input)) {
    return fallback
  }
  if (input < min) {
    return min
  }
  if (input > max) {
    return max
  }
  return input
}

function buildSectionSettings(section: AdminSection): Setting[] {
  if (section === 'site') {
    return [
      { key: settingKeys.siteName, value: siteSetting.siteName.trim(), type: 'string', desc: '站点名称' },
      { key: settingKeys.maintenanceMode, value: boolToValue(siteSetting.maintenanceMode), type: 'bool', desc: '维护模式' },
      { key: settingKeys.marquee, value: siteSetting.marquee.trim(), type: 'string', desc: '公告跑马灯' },
    ]
  }

  if (section === 'invite') {
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

  if (section === 'tracker') {
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

  return []
}

function pickChangedSettings(items: Setting[]): Setting[] {
  const current = settingsMap.value
  return items.filter((item) => {
    const key = item.key || ''
    if (!key) {
      return false
    }
    if (!Object.prototype.hasOwnProperty.call(current, key)) {
      return true
    }
    return (item.value || '') !== current[key]
  })
}

function syncSettingsToState(map: Record<string, string>) {
  siteSetting.siteName = map[settingKeys.siteName] || siteSetting.siteName
  siteSetting.maintenanceMode = parseBoolean(map[settingKeys.maintenanceMode], siteSetting.maintenanceMode)
  siteSetting.marquee = map[settingKeys.marquee] || siteSetting.marquee

  inviteSetting.openRegistration = parseBoolean(map[settingKeys.openRegistration], inviteSetting.openRegistration)
  inviteSetting.inviteOnly = parseBoolean(map[settingKeys.inviteOnly], inviteSetting.inviteOnly)
  inviteSetting.emailVerificationRequired = parseBoolean(
    map[settingKeys.emailVerificationRequired],
    inviteSetting.emailVerificationRequired,
  )
  inviteSetting.smtpEnable = parseBoolean(map[settingKeys.smtpEnable], inviteSetting.smtpEnable)
  inviteSetting.globalMessage = map[settingKeys.globalMessage] || inviteSetting.globalMessage

  trackerSetting.announceInterval = parseNumber(map[settingKeys.announceInterval], trackerSetting.announceInterval)
  trackerSetting.flushInterval = parseNumber(map[settingKeys.flushInterval], trackerSetting.flushInterval)
  trackerSetting.flushBatchSize = parseNumber(map[settingKeys.flushBatchSize], trackerSetting.flushBatchSize)
  trackerSetting.globalFreeleech = parseBoolean(map[settingKeys.globalFreeleech], trackerSetting.globalFreeleech)
  trackerSetting.freeleechCountdown = parseNumber(map[settingKeys.freeleechCountdown], trackerSetting.freeleechCountdown)
  trackerSetting.bonusFormula = map[settingKeys.bonusFormula] || trackerSetting.bonusFormula
  trackerSetting.trackerList = map[settingKeys.trackerList] || trackerSetting.trackerList
}

async function loadSystemSettings() {
  settingsLoading.value = true
  try {
    const response = await SystemService.GetSettings({})
    const map: Record<string, string> = {}
    for (const item of response.settings || []) {
      if (!item.key) {
        continue
      }
      map[item.key] = item.value || ''
    }
    settingsMap.value = map
    syncSettingsToState(map)
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '加载系统设置失败'
    MessagePlugin.error(message)
  } finally {
    settingsLoading.value = false
  }
}

async function saveActiveSection() {
  const candidatePayload = buildSectionSettings(activeSection.value)
  if (candidatePayload.length === 0) {
    MessagePlugin.info('当前分区暂无可保存配置')
    return
  }
  const payload = pickChangedSettings(candidatePayload)
  if (payload.length === 0) {
    MessagePlugin.info('配置未变化，无需保存')
    return
  }

  if (activeSection.value === 'invite' && inviteSetting.emailVerificationRequired && !inviteSetting.smtpEnable) {
    MessagePlugin.warning('已开启注册邮箱验证，请先开启 SMTP 发信开关')
    return
  }

  settingsSaving.value = true
  try {
    await SystemService.SetSettings({ settings: payload })

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

function flattenCategories(items: Category[] | undefined, depth = 0): CategoryRow[] {
  const result: CategoryRow[] = []
  for (const item of items || []) {
    const row: CategoryRow = {
      id: item.id || '',
      name: item.name || '未命名',
      slug: item.slug || '',
      description: item.description || '',
      parentId: item.parentId || '',
      icon: item.icon || '',
      sortOrder: item.sortOrder || 0,
      isEnabled: item.isEnabled !== false,
      depth,
    }
    result.push(row)

    if (item.children && item.children.length > 0) {
      result.push(...flattenCategories(item.children, depth + 1))
    }
  }
  return result
}

async function loadCategories() {
  categoryLoading.value = true
  try {
    const response = await CategoryService.ListCategories({})
    const flattened = flattenCategories(response.categories)
    categoryRows.value = flattened.sort((left, right) => left.sortOrder - right.sortOrder)
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '加载分类失败'
    MessagePlugin.error(message)
  } finally {
    categoryLoading.value = false
  }
}

function slugify(input: string): string {
  return input
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9\u4e00-\u9fa5\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
}

function resetCategoryForm() {
  categoryForm.name = ''
  categoryForm.slug = ''
  categoryForm.description = ''
  categoryForm.parentId = ''
  categoryForm.icon = ''
  categoryForm.sortOrder = 0
  categoryForm.isEnabled = true
}

function openCreateCategory() {
  categoryEditorMode.value = 'create'
  editingCategoryId.value = ''
  resetCategoryForm()
  categoryEditorVisible.value = true
}

function openEditCategory(row: CategoryRow) {
  categoryEditorMode.value = 'edit'
  editingCategoryId.value = row.id
  categoryForm.name = row.name
  categoryForm.slug = row.slug
  categoryForm.description = row.description
  categoryForm.parentId = row.parentId
  categoryForm.icon = row.icon
  categoryForm.sortOrder = row.sortOrder
  categoryForm.isEnabled = row.isEnabled
  categoryEditorVisible.value = true
}

function closeCategoryEditor() {
  categoryEditorVisible.value = false
  editingCategoryId.value = ''
  resetCategoryForm()
}

async function saveCategory() {
  const name = categoryForm.name.trim()
  const slug = (categoryForm.slug.trim() || slugify(name)).trim()

  if (!name) {
    MessagePlugin.warning('类别名称不能为空')
    return
  }

  if (!slug) {
    MessagePlugin.warning('Slug 不能为空')
    return
  }

  categorySaving.value = true
  try {
    if (categoryEditorMode.value === 'create') {
      await CategoryService.CreateCategory({
        name,
        slug,
        description: categoryForm.description.trim(),
        parentId: categoryForm.parentId,
        icon: categoryForm.icon.trim(),
        sortOrder: categoryForm.sortOrder,
      })
      MessagePlugin.success('类别创建成功')
    } else {
      await CategoryService.UpdateCategory({
        category: {
          id: editingCategoryId.value,
          name,
          slug,
          description: categoryForm.description.trim(),
          parentId: categoryForm.parentId,
          icon: categoryForm.icon.trim(),
          sortOrder: categoryForm.sortOrder,
          isEnabled: categoryForm.isEnabled,
        },
      })
      MessagePlugin.success('类别更新成功')
    }

    closeCategoryEditor()
    await loadCategories()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '保存类别失败'
    MessagePlugin.error(message)
  } finally {
    categorySaving.value = false
  }
}

async function removeCategory(row: CategoryRow) {
  const confirmed = window.confirm(`确认删除类别「${row.name}」吗？`)
  if (!confirmed) {
    return
  }

  try {
    await CategoryService.DeleteCategory({ id: row.id })
    MessagePlugin.success('类别已删除')
    await loadCategories()
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '删除类别失败'
    MessagePlugin.error(message)
  }
}

onMounted(async () => {
  await loadSystemSettings()
  if (activeSection.value === 'category') {
    await loadCategories()
  }
})
</script>

<template>
  <div class="grid gap-4 lg:grid-cols-[220px_1fr]">
    <aside class="rounded-lg border border-[var(--app-border)] bg-[var(--panel-bg)] p-2">
      <button
        v-for="item in visibleSections"
        :key="item.key"
        class="admin-menu-btn"
        :class="item.key === activeSection ? 'active' : ''"
        @click="activeSection = item.key"
      >
        {{ item.label }}
      </button>

      <p v-if="visibleSections.length === 0" class="px-2 py-3 text-xs text-rose-500">
        当前账号没有后台分区权限
      </p>
    </aside>

    <section>
      <t-card v-if="activeSection === 'site'" title="站点基础设置" size="small" :loading="settingsLoading">
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
            <t-button theme="primary" :loading="settingsSaving" @click="saveSection">保存设置</t-button>
            <t-button variant="outline" :loading="settingsLoading" @click="loadSystemSettings">刷新设置</t-button>
          </div>
        </template>
      </t-card>

      <t-card v-else-if="activeSection === 'invite'" title="注册与邀请控制" size="small" :loading="settingsLoading">
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
            <t-button theme="primary" :loading="settingsSaving" @click="saveSection">保存设置</t-button>
            <t-button variant="outline" @click="runAction('发送全局消息')">发送全局消息</t-button>
          </div>
        </template>
      </t-card>

      <t-card v-else-if="activeSection === 'tracker'" title="Tracker 及业务参数" size="small" :loading="settingsLoading">
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
            <t-button theme="primary" :loading="settingsSaving" @click="saveSection">保存设置</t-button>
            <t-button variant="outline" :loading="settingsLoading" @click="loadSystemSettings">刷新设置</t-button>
          </div>
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

      <t-card v-else-if="activeSection === 'category'" title="种子类别管理" size="small">
        <div class="mb-3 flex items-center gap-2">
          <t-button theme="primary" @click="openCreateCategory">新建类别</t-button>
          <t-button variant="outline" :loading="categoryLoading" @click="loadCategories">刷新列表</t-button>
        </div>

        <t-table
          row-key="id"
          size="small"
          bordered
          hover
          :columns="categoryColumns"
          :data="categoryRows"
          :loading="categoryLoading"
        >
          <template #name="{ row }">
            <div class="flex items-center gap-1" :style="{ paddingLeft: `${row.depth * 16}px` }">
              <span>{{ row.icon || '📁' }}</span>
              <span>{{ row.name }}</span>
            </div>
          </template>

          <template #isEnabled="{ row }">
            <t-tag size="small" :theme="row.isEnabled ? 'success' : 'danger'" variant="light">
              {{ row.isEnabled ? '启用' : '禁用' }}
            </t-tag>
          </template>

          <template #actions="{ row }">
            <div class="flex items-center gap-2">
              <t-link theme="primary" hover="color" @click="openEditCategory(row)">编辑</t-link>
              <t-link theme="danger" hover="color" @click="removeCategory(row)">删除</t-link>
            </div>
          </template>
        </t-table>

        <div v-if="categoryEditorVisible" class="category-editor">
          <p class="mb-3 text-sm font-600">
            {{ categoryEditorMode === 'create' ? '新建类别' : '编辑类别' }}
          </p>

          <div class="form-grid">
            <label>
              <span>名称 *</span>
              <t-input v-model="categoryForm.name" clearable placeholder="例如：Movies" />
            </label>

            <label>
              <span>Slug *</span>
              <t-input v-model="categoryForm.slug" clearable placeholder="例如：movies" />
            </label>

            <label>
              <span>父类别</span>
              <t-select v-model="categoryForm.parentId" clearable placeholder="无">
                <t-option
                  v-for="item in parentCategoryOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </t-select>
            </label>

            <label>
              <span>图标</span>
              <t-input v-model="categoryForm.icon" clearable placeholder="例如：🎬" />
            </label>

            <label>
              <span>排序</span>
              <t-input-number v-model="categoryForm.sortOrder" :min="0" />
            </label>

            <label>
              <span>启用</span>
              <t-switch
                :value="categoryForm.isEnabled"
                @change="(value: boolean) => (categoryForm.isEnabled = value)"
              />
            </label>

            <label class="lg:col-span-2">
              <span>描述</span>
              <t-textarea v-model="categoryForm.description" :autosize="{ minRows: 2, maxRows: 4 }" />
            </label>
          </div>

          <div class="mt-3 flex gap-2">
            <t-button theme="primary" :loading="categorySaving" @click="saveCategory">保存类别</t-button>
            <t-button variant="outline" @click="closeCategoryEditor">取消</t-button>
          </div>
        </div>
      </t-card>

      <t-card v-else-if="activeSection === 'user'" title="用户管理" size="small">
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
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 6px;
  font-size: 12px;
  color: var(--muted-text);
}

.form-grid label :deep(.t-switch),
.form-grid label :deep(.t-checkbox),
.form-grid label :deep(.t-radio-group) {
  align-self: flex-start;
  width: auto;
}

.category-editor {
  margin-top: 14px;
  border-top: 1px dashed var(--app-border);
  padding-top: 14px;
}
</style>
