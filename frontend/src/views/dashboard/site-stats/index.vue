<script lang="ts" setup>
import { useUserStore } from '@/store'
import * as echarts from 'echarts'
import { onMounted, ref } from 'vue'

const UserStore = useUserStore()
const isAdmin = ref<boolean>(false)

// 模拟数据
const stats = ref({
  totalUsers: 0,
  totalProducts: 0,
  totalProjects: 0,
  activeUsers: 0,
  updateTime: '2024-03-20 12:00:00',
  userStats: {
    total: 644,
    inactive: 0,
    capacity: 600,
    todayActive: 310,
    weeklyActive: 413,
  },
  torrentStats: {
    total: 38,
    dead: 6,
    size: '1.368 TiB',
  },
  connectionStats: {
    total: 4147,
    download: 149,
    seed: 4143,
  },
  transferStats: {
    totalUpload: '4.337 TiB',
    totalDownload: '4.092 TiB',
    displayUpload: '58.803 TiB',
    displayDownload: '2.265 TiB',
  },
})

// 图表数据
const userGrowthData = ref([
  { month: '1月', users: 100 },
  { month: '2月', users: 200 },
  { month: '3月', users: 300 },
  { month: '4月', users: 400 },
  { month: '5月', users: 500 },
  { month: '6月', users: 600 },
])

const userDistributionData = ref([
  { name: '活跃用户', value: 413 },
  { name: '非活跃用户', value: 231 },
])

let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null

onMounted(() => {
  isAdmin.value = UserStore.role === 'admin'
  // 这里可以添加获取实际数据的API调用
  stats.value = {
    totalUsers: 1234,
    totalProducts: 567,
    totalProjects: 890,
    activeUsers: 234,
    updateTime: '2024-03-20 12:00:00',
    userStats: {
      total: 644,
      inactive: 0,
      capacity: 600,
      todayActive: 310,
      weeklyActive: 413,
    },
    torrentStats: {
      total: 38,
      dead: 6,
      size: '1.368 TiB',
    },
    connectionStats: {
      total: 4147,
      download: 149,
      seed: 4143,
    },
    transferStats: {
      totalUpload: '4.337 TiB',
      totalDownload: '4.092 TiB',
      displayUpload: '58.803 TiB',
      displayDownload: '2.265 TiB',
    },
  }

  // 初始化折线图
  const lineChartDom = document.getElementById('userGrowthChart')
  if (lineChartDom) {
    lineChart = echarts.init(lineChartDom)
    lineChart.setOption({
      tooltip: {
        trigger: 'axis',
      },
      xAxis: {
        type: 'category',
        data: userGrowthData.value.map(item => item.month),
      },
      yAxis: {
        type: 'value',
      },
      series: [
        {
          data: userGrowthData.value.map(item => item.users),
          type: 'line',
          smooth: true,
        },
      ],
    })
  }

  // 初始化饼图
  const pieChartDom = document.getElementById('userDistributionChart')
  if (pieChartDom) {
    pieChart = echarts.init(pieChartDom)
    pieChart.setOption({
      tooltip: {
        trigger: 'item',
      },
      legend: {
        orient: 'vertical',
        left: 'left',
      },
      series: [
        {
          type: 'pie',
          radius: '50%',
          data: userDistributionData.value,
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)',
            },
          },
        },
      ],
    })
  }
})
</script>

<template>
  <div class="flex p-6">
    <div class="w-full flex flex-col gap-8">
      <!-- 标题和时间 -->
      <div class="w-full bg-[--color-bg-2] p-8">
        <div class="flex items-center justify-between">
          <a-typography-title :heading="4" class="m-0">
            全站数据统计
          </a-typography-title>
          <a-typography-text type="secondary" class="text-sm">
            更新时间：{{ new Date().toLocaleString() }}
          </a-typography-text>
        </div>
      </div>

      <!-- 基础统计卡片 -->
      <div class="w-full bg-[--color-bg-2] p-8">
        <a-row :gutter="32">
          <a-col :span="6">
            <a-statistic
              title="总用户数"
              :value="stats.totalUsers"
              show-group-separator
              class="statistic-card"
            />
          </a-col>
          <a-col :span="6">
            <a-statistic
              title="总产品数"
              :value="stats.totalProducts"
              show-group-separator
              class="statistic-card"
            />
          </a-col>
          <a-col :span="6">
            <a-statistic
              title="总项目数"
              :value="stats.totalProjects"
              show-group-separator
              class="statistic-card"
            />
          </a-col>
          <a-col :span="6">
            <a-statistic
              title="活跃用户数"
              :value="stats.activeUsers"
              show-group-separator
              class="statistic-card"
            />
          </a-col>
        </a-row>
      </div>

      <!-- 详细统计信息 -->
      <div class="w-full bg-[--color-bg-2] p-8">
        <a-row :gutter="32">
          <!-- 用户统计 -->
          <a-col :span="6">
            <a-card title="用户统计" :bordered="false" class="detail-card">
              <div class="stat-item">
                <div class="stat-label">
                  总用户数
                </div>
                <div class="stat-value">
                  {{ stats.userStats.total }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  今日活跃
                </div>
                <div class="stat-value">
                  {{ stats.userStats.todayActive }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  本周活跃
                </div>
                <div class="stat-value">
                  {{ stats.userStats.weeklyActive }}
                </div>
              </div>
            </a-card>
          </a-col>

          <!-- 种子统计 -->
          <a-col :span="6">
            <a-card title="种子统计" :bordered="false" class="detail-card">
              <div class="stat-item">
                <div class="stat-label">
                  种子数量
                </div>
                <div class="stat-value">
                  {{ stats.torrentStats.total }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  断种数量
                </div>
                <div class="stat-value">
                  {{ stats.torrentStats.dead }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  种子体积
                </div>
                <div class="stat-value">
                  {{ stats.torrentStats.size }}
                </div>
              </div>
            </a-card>
          </a-col>

          <!-- 连接统计 -->
          <a-col :span="6">
            <a-card title="连接统计" :bordered="false" class="detail-card">
              <div class="stat-item">
                <div class="stat-label">
                  总连接数
                </div>
                <div class="stat-value">
                  {{ stats.connectionStats.total }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  下载连接
                </div>
                <div class="stat-value">
                  {{ stats.connectionStats.download }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  做种连接
                </div>
                <div class="stat-value">
                  {{ stats.connectionStats.seed }}
                </div>
              </div>
            </a-card>
          </a-col>

          <!-- 传输统计 -->
          <a-col :span="6">
            <a-card title="传输统计" :bordered="false" class="detail-card">
              <div class="stat-item">
                <div class="stat-label">
                  全站上传
                </div>
                <div class="stat-value">
                  {{ stats.transferStats.totalUpload }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  全站下载
                </div>
                <div class="stat-value">
                  {{ stats.transferStats.totalDownload }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  显示上传
                </div>
                <div class="stat-value">
                  {{ stats.transferStats.displayUpload }}
                </div>
              </div>
              <div class="stat-item">
                <div class="stat-label">
                  显示下载
                </div>
                <div class="stat-value">
                  {{ stats.transferStats.displayDownload }}
                </div>
              </div>
            </a-card>
          </a-col>
        </a-row>
      </div>

      <!-- 图表展示 -->
      <div class="w-full bg-[--color-bg-2] p-8">
        <a-row :gutter="32">
          <a-col :span="12">
            <a-card title="用户增长趋势" :bordered="false" class="chart-card">
              <div id="userGrowthChart" style="width: 100%; height: 300px" />
            </a-card>
          </a-col>
          <a-col :span="12">
            <a-card title="用户分布" :bordered="false" class="chart-card">
              <div id="userDistributionChart" style="width: 100%; height: 300px" />
            </a-card>
          </a-col>
        </a-row>
      </div>
    </div>
  </div>
</template>

<style lang="less" scoped>
:deep(.arco-card) {
  margin-bottom: 0;
}

:deep(.arco-statistic) {
  margin-bottom: 0;
}

.detail-card {
  height: 100%;
  padding: 24px;

  :deep(.arco-card-header) {
    padding-bottom: 24px;
    border-bottom: 1px solid var(--color-border);
    margin-bottom: 16px;
  }
}

.chart-card {
  height: 100%;
  padding: 24px;

  :deep(.arco-card-header) {
    padding-bottom: 24px;
    border-bottom: 1px solid var(--color-border);
    margin-bottom: 16px;
  }
}

.statistic-card {
  padding: 24px;
  background-color: var(--color-bg-2);
  border-radius: 4px;
  height: 100%;
}

.mb-8 {
  margin-bottom: 32px;
}

:deep(.arco-statistic-title) {
  font-size: 14px;
  margin-bottom: 8px;
}

:deep(.arco-statistic-value) {
  font-size: 24px;
}

.stat-item {
  margin-bottom: 40px;

  .stat-label {
    color: var(--color-text-3);
    font-size: 14px;
    margin-bottom: 12px;
  }

  .stat-value {
    color: var(--color-text-1);
    font-size: 28px;
    font-weight: 500;
  }
}
</style>
