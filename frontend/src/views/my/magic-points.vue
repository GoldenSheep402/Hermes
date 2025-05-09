<script lang="ts" setup>
import { useMagicStore } from '@/store/modules/magic'
import { UserService } from '@/services/grpc'
import { Message } from '@arco-design/web-vue'
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

interface MagicPointRecord {
  id: number
  time: string
  type: string
  amount: number
  balance: number
  description: string
}

const records = ref<MagicPointRecord[]>([
  {
    id: 1,
    time: '2023-05-10 15:30:22',
    type: '收入',
    amount: 100,
    balance: 15860,
    description: '做种时间奖励',
  },
  {
    id: 2,
    time: '2023-05-09 12:45:33',
    type: '支出',
    amount: -50,
    balance: 15760,
    description: '购买邀请码',
  },
  {
    id: 3,
    time: '2023-05-08 09:12:15',
    type: '收入',
    amount: 200,
    balance: 15810,
    description: '发布种子奖励',
  },
  {
    id: 4,
    time: '2023-05-07 18:22:40',
    type: '收入',
    amount: 80,
    balance: 15610,
    description: '论坛活动奖励',
  },
  {
    id: 5,
    time: '2023-05-06 14:10:05',
    type: '支出',
    amount: -120,
    balance: 15530,
    description: '购买促销',
  },
])

const totalPoints = ref(15860)
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 5,
})

const loading = ref(false)

// 过滤条件
const filterForm = ref({
  dateRange: null,
  type: null,
  description: '',
})

const typeOptions = [
  { label: '全部', value: null },
  { label: '收入', value: '收入' },
  { label: '支出', value: '支出' },
]

const router = useRouter()
const magicStore = useMagicStore()

async function fetchMagicPointsRecords() {
  loading.value = true
  try {
    // 模拟API调用，真实项目中应该调用实际接口
    // const res = await UserService.GetMagicPointsRecords({
    //   page: pagination.current,
    //   pageSize: pagination.pageSize,
    //   ...filterForm.value,
    // })
    // records.value = res.records
    // pagination.total = res.total
    // totalPoints.value = res.currentPoints
    
    // 使用示例数据
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 无需覆盖已经定义的示例数据
    pagination.total = records.value.length
  }
  catch (error) {
    console.error('获取魔力值记录失败:', error)
    Message.error('获取魔力值记录失败，请稍后重试')
    records.value = [] // 出错时清空记录
  }
  finally {
    loading.value = false
  }
}

function handlePageChange(page: number) {
  pagination.current = page
  fetchMagicPointsRecords()
}

function handleSearch() {
  pagination.current = 1
  fetchMagicPointsRecords()
}

function handleReset() {
  filterForm.value = {
    dateRange: null,
    type: null,
    description: '',
  }
  pagination.current = 1
  fetchMagicPointsRecords()
}

function goToMagicShop() {
  router.push('/magic-shop')
}

onMounted(() => {
  fetchMagicPointsRecords()
})
</script>

<template>
  <div class="h-full p-5">
    <div class="h-full bg-[--color-bg-2] p-5">
      <div class="mb-5 flex justify-between items-center">
        <div class="p-0.5 text-30px text-[--color-text-1] font-500 leading-[1.4]">
          魔力值明细
        </div>
        <div class="magic-points-summary">
          当前魔力值：<span class="magic-points-value">{{ magicStore.currentMagicPoints.toLocaleString() }}</span>
        </div>
      </div>

      <!-- 过滤条件 -->
      <a-card class="mb-4">
        <a-form :model="filterForm" layout="inline">
          <a-form-item field="dateRange" label="日期范围">
            <a-range-picker v-model="filterForm.dateRange" style="width: 240px" />
          </a-form-item>
          <a-form-item field="type" label="类型">
            <a-select v-model="filterForm.type" :options="typeOptions" style="width: 100px" />
          </a-form-item>
          <a-form-item field="description" label="描述">
            <a-input v-model="filterForm.description" placeholder="搜索描述" />
          </a-form-item>
          <a-form-item>
            <a-space>
              <a-button type="primary" @click="handleSearch">
                <template #icon>
                  <icon-search />
                </template>
                搜索
              </a-button>
              <a-button @click="handleReset">
                <template #icon>
                  <icon-refresh />
                </template>
                重置
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>

      <!-- 记录表格 -->
      <a-card>
        <a-table
          :data="records"
          :pagination="{
            current: pagination.current,
            pageSize: pagination.pageSize,
            total: pagination.total,
            onChange: handlePageChange,
          }"
          :bordered="false"
          :loading="loading"
        >
          <template #columns>
            <a-table-column title="时间" data-index="time" />
            <a-table-column title="类型" data-index="type">
              <template #cell="{ record }">
                <a-tag :color="record.amount > 0 ? 'green' : 'red'">
                  {{ record.type }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column title="数量" data-index="amount">
              <template #cell="{ record }">
                <span :class="record.amount > 0 ? 'text-green-500' : 'text-red-500'">
                  {{ record.amount > 0 ? '+' : '' }}{{ record.amount }}
                </span>
              </template>
            </a-table-column>
            <a-table-column title="余额" data-index="balance" />
            <a-table-column title="描述" data-index="description" />
          </template>
          <template #empty>
            <div class="empty-state">
              <icon-empty />
              <p>暂无魔力值记录</p>
            </div>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>
</template>

<style scoped lang="less">
.magic-points-summary {
  font-size: 16px;
  color: var(--color-text-2);
  .magic-points-value {
    font-size: 20px;
    font-weight: 700;
    color: var(--color-primary-6);
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  color: var(--color-text-3);
  :deep(.arco-icon) {
    font-size: 48px;
    margin-bottom: 16px;
  }
  p {
    font-size: 14px;
  }
}
</style>
