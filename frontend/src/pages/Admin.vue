<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import type { Category } from '@/lib/proto/category/v1/category.pb'
import type { PermissionKey } from '@/constants/permissions'
import { PermissionKeys } from '@/constants/permissions'
import { CategoryService } from '@/services/grpc'
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
