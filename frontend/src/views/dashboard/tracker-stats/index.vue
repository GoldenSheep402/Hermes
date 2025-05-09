<script lang="ts" setup>
import * as echarts from 'echarts'
import { onMounted, ref } from 'vue'

interface DataPoint {
  time: string
  value: number
}

interface TrackerData {
  connections: DataPoint[]
  requests: DataPoint[]
  bandwidth: DataPoint[]
}

// 模拟 tracker 负载数据
const trackerData = ref<TrackerData>({
  connections: [
    { time: '00:00', value: 1000 },
    { time: '01:00', value: 1200 },
    { time: '02:00', value: 1500 },
    { time: '03:00', value: 1800 },
    { time: '04:00', value: 2000 },
    { time: '05:00', value: 2200 },
    { time: '06:00', value: 2500 },
    { time: '07:00', value: 2800 },
    { time: '08:00', value: 3000 },
    { time: '09:00', value: 3200 },
    { time: '10:00', value: 3500 },
    { time: '11:00', value: 3800 },
    { time: '12:00', value: 4000 },
    { time: '13:00', value: 3800 },
    { time: '14:00', value: 3500 },
    { time: '15:00', value: 3200 },
    { time: '16:00', value: 3000 },
    { time: '17:00', value: 2800 },
    { time: '18:00', value: 2500 },
    { time: '19:00', value: 2200 },
    { time: '20:00', value: 2000 },
    { time: '21:00', value: 1800 },
    { time: '22:00', value: 1500 },
    { time: '23:00', value: 1200 },
  ],
  requests: [
    { time: '00:00', value: 500 },
    { time: '01:00', value: 600 },
    { time: '02:00', value: 750 },
    { time: '03:00', value: 900 },
    { time: '04:00', value: 1000 },
    { time: '05:00', value: 1100 },
    { time: '06:00', value: 1250 },
    { time: '07:00', value: 1400 },
    { time: '08:00', value: 1500 },
    { time: '09:00', value: 1600 },
    { time: '10:00', value: 1750 },
    { time: '11:00', value: 1900 },
    { time: '12:00', value: 2000 },
    { time: '13:00', value: 1900 },
    { time: '14:00', value: 1750 },
    { time: '15:00', value: 1600 },
    { time: '16:00', value: 1500 },
    { time: '17:00', value: 1400 },
    { time: '18:00', value: 1250 },
    { time: '19:00', value: 1100 },
    { time: '20:00', value: 1000 },
    { time: '21:00', value: 900 },
    { time: '22:00', value: 750 },
    { time: '23:00', value: 600 },
  ],
  bandwidth: [
    { time: '00:00', value: 50 },
    { time: '01:00', value: 60 },
    { time: '02:00', value: 75 },
    { time: '03:00', value: 90 },
    { time: '04:00', value: 100 },
    { time: '05:00', value: 110 },
    { time: '06:00', value: 125 },
    { time: '07:00', value: 140 },
    { time: '08:00', value: 150 },
    { time: '09:00', value: 160 },
    { time: '10:00', value: 175 },
    { time: '11:00', value: 190 },
    { time: '12:00', value: 200 },
    { time: '13:00', value: 190 },
    { time: '14:00', value: 175 },
    { time: '15:00', value: 160 },
    { time: '16:00', value: 150 },
    { time: '17:00', value: 140 },
    { time: '18:00', value: 125 },
    { time: '19:00', value: 110 },
    { time: '20:00', value: 100 },
    { time: '21:00', value: 90 },
    { time: '22:00', value: 75 },
    { time: '23:00', value: 60 },
  ],
})

let trackerChart: echarts.ECharts | null = null

onMounted(() => {
  const chartDom = document.getElementById('trackerChart')
  if (chartDom) {
    trackerChart = echarts.init(chartDom)
    trackerChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          label: {
            backgroundColor: '#6a7985',
          },
        },
      },
      legend: {
        data: ['连接数', '请求数', '带宽使用率'],
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: trackerData.value.connections.map(item => item.time),
      },
      yAxis: [
        {
          type: 'value',
          name: '数量',
          position: 'left',
        },
        {
          type: 'value',
          name: '带宽使用率(%)',
          position: 'right',
          axisLabel: {
            formatter: '{value}%',
          },
        },
      ],
      series: [
        {
          name: '连接数',
          type: 'line',
          smooth: true,
          data: trackerData.value.connections.map(item => item.value),
          areaStyle: {
            opacity: 0.1,
          },
        },
        {
          name: '请求数',
          type: 'line',
          smooth: true,
          data: trackerData.value.requests.map(item => item.value),
          areaStyle: {
            opacity: 0.1,
          },
        },
        {
          name: '带宽使用率',
          type: 'line',
          smooth: true,
          yAxisIndex: 1,
          data: trackerData.value.bandwidth.map(item => item.value),
          areaStyle: {
            opacity: 0.1,
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
            Tracker 负载统计
          </a-typography-title>
          <a-typography-text type="secondary" class="text-sm">
            更新时间：{{ new Date().toLocaleString() }}
          </a-typography-text>
        </div>
      </div>

      <!-- 图表展示 -->
      <div class="w-full bg-[--color-bg-2] p-8">
        <a-card title="Tracker 负载趋势" :bordered="false">
          <div id="trackerChart" style="width: 100%; height: 500px" />
        </a-card>
      </div>

      <!-- 统计信息 -->
      <div class="w-full bg-[--color-bg-2] p-8">
        <a-row :gutter="32">
          <a-col :span="24">
            <a-card title="当前状态" :bordered="false" class="mb-4">
              <a-row :gutter="32">
                <a-col :span="8">
                  <a-statistic
                    title="活跃连接数"
                    :value="trackerData.connections[new Date().getHours()].value"
                    show-group-separator
                  />
                </a-col>
                <a-col :span="8">
                  <a-statistic
                    title="每秒请求数"
                    :value="trackerData.requests[new Date().getHours()].value"
                    show-group-separator
                  />
                </a-col>
                <a-col :span="8">
                  <a-statistic
                    title="带宽使用率"
                    :value="trackerData.bandwidth[new Date().getHours()].value"
                    :precision="1"
                  >
                    <template #suffix>
                      %
                    </template>
                  </a-statistic>
                </a-col>
              </a-row>
            </a-card>
          </a-col>
          <a-col :span="24">
            <a-card title="峰值信息" :bordered="false" class="mb-4">
              <a-row :gutter="32">
                <a-col :span="8">
                  <a-statistic
                    title="最大连接数"
                    :value="Math.max(...trackerData.connections.map((item: DataPoint) => item.value))"
                    show-group-separator
                  />
                </a-col>
                <a-col :span="8">
                  <a-statistic
                    title="最大请求数"
                    :value="Math.max(...trackerData.requests.map((item: DataPoint) => item.value))"
                    show-group-separator
                  />
                </a-col>
                <a-col :span="8">
                  <a-statistic
                    title="最大带宽使用率"
                    :value="Math.max(...trackerData.bandwidth.map((item: DataPoint) => item.value))"
                    :precision="1"
                  >
                    <template #suffix>
                      %
                    </template>
                  </a-statistic>
                </a-col>
              </a-row>
            </a-card>
          </a-col>
          <a-col :span="24">
            <a-card title="平均信息" :bordered="false">
              <a-row :gutter="32">
                <a-col :span="8">
                  <a-statistic
                    title="平均连接数"
                    :value="Math.round(trackerData.connections.reduce((acc, item) => acc + item.value, 0) / trackerData.connections.length)"
                    show-group-separator
                  />
                </a-col>
                <a-col :span="8">
                  <a-statistic
                    title="平均请求数"
                    :value="Math.round(trackerData.requests.reduce((acc, item) => acc + item.value, 0) / trackerData.requests.length)"
                    show-group-separator
                  />
                </a-col>
                <a-col :span="8">
                  <a-statistic
                    title="平均带宽使用率"
                    :value="Number((trackerData.bandwidth.reduce((acc, item) => acc + item.value, 0) / trackerData.bandwidth.length).toFixed(1))"
                    :precision="1"
                  >
                    <template #suffix>
                      %
                    </template>
                  </a-statistic>
                </a-col>
              </a-row>
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
  margin-bottom: 16px;
}
</style>
