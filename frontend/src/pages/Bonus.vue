<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { BonusService, InviteService } from '@/services/grpc'
import { useAuthStore } from '@/store'
import { formatBonusPoints, formatDateTime, formatNumber } from '@/utils/format'

const authStore = useAuthStore()

const loading = ref(false)
const exchanging = ref(false)
const balanceMilli = ref(0)
const shopItems = ref<Array<{ id: string; name: string; pricePoints: number; unit: string }>>([])
const logs = ref<
  Array<{
    id: string
    amount: number
    reason: string
    description: string
    createdAt: string
  }>
>([])
const logTotal = ref(0)
const logPage = reactive({ current: 1, pageSize: 20 })
const inviteCodes = ref<string[]>([])

const exchangeForm = reactive({
  item: 'upload',
  quantity: 1,
})

const balanceDisplay = computed(() => Math.floor(balanceMilli.value / 1000))
const selectedShop = computed(() => shopItems.value.find((item) => item.id === exchangeForm.item))
const estimatedCost = computed(() => {
  if (!selectedShop.value) {
    return 0
  }
  return selectedShop.value.pricePoints * Math.max(1, Number(exchangeForm.quantity) || 1)
})

function parseInt64(value: string | number | undefined): number {
  if (typeof value === 'number') {
    return Number.isFinite(value) ? value : 0
  }
  if (typeof value === 'string' && value.trim() !== '') {
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : 0
  }
  return 0
}

function reasonLabel(reason: string): string {
  switch (reason) {
    case 'seeding':
      return '做种奖励'
    case 'exchange_upload':
      return '兑换上传量'
    case 'exchange_invite':
      return '兑换邀请'
    case 'gift':
      return '赠送'
    case 'system':
      return '系统'
    default:
      return reason || '其他'
  }
}

async function loadBalance() {
  const response = await BonusService.GetBonusBalance({})
  balanceMilli.value = parseInt64(response.balance)
  authStore.profile.bonusPoints = Math.floor(balanceMilli.value / 1000)
}

async function loadShop() {
  const response = await BonusService.ListShopItems({})
  shopItems.value = (response.items || []).map((item) => ({
    id: item.id || '',
    name: item.name || item.id || '',
    pricePoints: parseInt64(item.pricePoints),
    unit: item.unit || '',
  }))
  if (!shopItems.value.some((item) => item.id === exchangeForm.item) && shopItems.value[0]) {
    exchangeForm.item = shopItems.value[0].id
  }
}

async function loadLogs() {
  const response = await BonusService.ListBonusLogs({
    page: logPage.current,
    pageSize: logPage.pageSize,
  })
  logTotal.value = parseInt64(response.total)
  logs.value = (response.logs || []).map((item) => ({
    id: item.id || '',
    amount: parseInt64(item.amount),
    reason: item.reason || '',
    description: item.description || '',
    createdAt: item.createdAt || '',
  }))
}

async function refreshAll() {
  loading.value = true
  try {
    await Promise.all([loadBalance(), loadShop(), loadLogs()])
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '加载魔力信息失败'
    MessagePlugin.error(message)
  } finally {
    loading.value = false
  }
}

async function handleExchange() {
  if (!exchangeForm.item) {
    MessagePlugin.warning('请选择兑换商品')
    return
  }
  const quantity = Math.max(1, Math.floor(Number(exchangeForm.quantity) || 1))
  exchanging.value = true
  try {
    const response = await BonusService.ExchangeBonus({
      item: exchangeForm.item,
      quantity: String(quantity),
    })
    balanceMilli.value = parseInt64(response.remainingBalance)
    authStore.profile.bonusPoints = Math.floor(balanceMilli.value / 1000)
    MessagePlugin.success('兑换成功')
    if (exchangeForm.item === 'invite') {
      authStore.profile.inviteCount += quantity
    }
    await Promise.all([loadLogs(), authStore.fetchProfile()])
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '兑换失败'
    MessagePlugin.error(message)
  } finally {
    exchanging.value = false
  }
}

async function createInviteCodes() {
  try {
    const response = await InviteService.CreateInviteCode({ count: 1 })
    inviteCodes.value = response.codes || []
    await authStore.fetchProfile()
    MessagePlugin.success('已生成邀请码')
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : '生成邀请码失败'
    MessagePlugin.error(message)
  }
}

function onLogPageChange(pageInfo: { current: number; pageSize: number }) {
  logPage.current = pageInfo.current
  logPage.pageSize = pageInfo.pageSize
  void loadLogs()
}

onMounted(() => {
  void refreshAll()
})
</script>

<template>
  <section class="bonus-page grid gap-4 xl:grid-cols-[minmax(0,1fr)_360px]">
    <div class="grid gap-4">
      <t-card title="魔力余额" size="small" :loading="loading">
        <div class="balance-row">
          <div>
            <p class="balance-label">当前可用</p>
            <p class="balance-value">{{ formatNumber(balanceDisplay) }}</p>
            <p class="balance-hint">毫点 {{ formatNumber(balanceMilli) }}（1000 毫点 = 1 魔力）</p>
          </div>
          <t-button variant="outline" :loading="loading" @click="refreshAll">刷新</t-button>
        </div>
      </t-card>

      <t-card title="魔力流水" size="small" :loading="loading">
        <t-table
          row-key="id"
          :data="logs"
          :columns="[
            { colKey: 'createdAt', title: '时间', width: 140 },
            { colKey: 'reason', title: '类型', width: 120 },
            { colKey: 'amount', title: '变动(毫点)', width: 120 },
            { colKey: 'description', title: '说明' },
          ]"
          :pagination="{
            current: logPage.current,
            pageSize: logPage.pageSize,
            total: logTotal,
          }"
          @page-change="onLogPageChange"
        >
          <template #createdAt="{ row }">
            {{ row.createdAt ? formatDateTime(row.createdAt) : '-' }}
          </template>
          <template #reason="{ row }">
            {{ reasonLabel(row.reason) }}
          </template>
          <template #amount="{ row }">
            <span :class="row.amount >= 0 ? 'text-emerald-600' : 'text-rose-600'">
              {{ row.amount >= 0 ? '+' : '' }}{{ formatBonusPoints(row.amount, true) }}
              <small class="text-muted">({{ row.amount }})</small>
            </span>
          </template>
        </t-table>
      </t-card>
    </div>

    <div class="grid gap-4">
      <t-card title="魔力商店" size="small" :loading="loading">
        <div class="shop-form">
          <label>
            <span>商品</span>
            <t-select v-model="exchangeForm.item">
              <t-option
                v-for="item in shopItems"
                :key="item.id"
                :value="item.id"
                :label="`${item.name}（${item.pricePoints} 魔力 / ${item.unit}）`"
              />
            </t-select>
          </label>
          <label>
            <span>数量</span>
            <t-input-number v-model="exchangeForm.quantity" :min="1" :step="1" />
          </label>
          <p class="cost-line">
            预计消耗：<strong>{{ formatNumber(estimatedCost) }}</strong> 魔力
          </p>
          <t-button theme="primary" block :loading="exchanging" @click="handleExchange">确认兑换</t-button>
        </div>
      </t-card>

      <t-card title="邀请名额" size="small">
        <p class="line-item">
          <span>当前库存</span>
          <strong>{{ formatNumber(authStore.profile.inviteCount) }}</strong>
        </p>
        <t-button class="mt-3" variant="outline" block @click="createInviteCodes">生成 1 个邀请码</t-button>
        <div v-if="inviteCodes.length" class="invite-list">
          <p v-for="code in inviteCodes" :key="code" class="invite-code">{{ code }}</p>
        </div>
      </t-card>
    </div>
  </section>
</template>

<style scoped>
.balance-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.balance-label {
  margin: 0;
  font-size: 12px;
  color: var(--muted-text);
}

.balance-value {
  margin: 4px 0;
  font-size: 32px;
  font-weight: 700;
}

.balance-hint {
  margin: 0;
  font-size: 12px;
  color: var(--muted-text);
}

.shop-form {
  display: grid;
  gap: 12px;
}

.shop-form label {
  display: grid;
  gap: 6px;
  font-size: 12px;
  color: var(--muted-text);
}

.cost-line {
  margin: 0;
  font-size: 13px;
}

.line-item {
  margin: 0;
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.invite-list {
  margin-top: 12px;
  display: grid;
  gap: 6px;
}

.invite-code {
  margin: 0;
  padding: 8px;
  border-radius: 6px;
  background: var(--soft-bg);
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  word-break: break-all;
}

.text-muted {
  color: var(--muted-text);
}

.mt-3 {
  margin-top: 12px;
}
</style>
