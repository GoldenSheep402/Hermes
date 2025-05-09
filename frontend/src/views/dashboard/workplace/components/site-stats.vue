<script lang="ts" setup>
import { ref } from 'vue'

const siteStats = ref({
  updateTime: '2025-04-09 00:30:00',
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

const disclaimer = ref(`本站列出的文件并没有保存在本站的服务器上。本站仅负责连接, 我们对被传播文件的内容一无所知。本站的链接均由用户自发提供。Hermes 的管理员无法对用户的提交内容或其他行为负责。您不能使用 Hermes 传播或下载任何您没有使用许可的文件或材料, 由此引起的后果由您自己承担。`)
</script>

<template>
  <div class="site-stats-container">
    <a-card>
      <template #title>
        站点信息
      </template>
      <div class="grid grid-cols-4 gap-4">
        <!-- 用户统计 -->
        <div class="stat-group">
          <div class="stat-item">
            <div class="stat-label">
              更新时间
            </div>
            <div class="stat-value">
              {{ new Date().toLocaleString() }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              用户数量
            </div>
            <div class="stat-value">
              {{ siteStats.userStats.total }} / {{ siteStats.userStats.inactive }} / {{ siteStats.userStats.capacity }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              今日活跃
            </div>
            <div class="stat-value">
              {{ siteStats.userStats.todayActive }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              本周活跃
            </div>
            <div class="stat-value">
              {{ siteStats.userStats.weeklyActive }}
            </div>
          </div>
        </div>

        <!-- 种子统计 -->
        <div class="stat-group">
          <div class="stat-item">
            <div class="stat-label">
              种子数量
            </div>
            <div class="stat-value">
              {{ siteStats.torrentStats.total }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              断种数量
            </div>
            <div class="stat-value">
              {{ siteStats.torrentStats.dead }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              种子体积
            </div>
            <div class="stat-value">
              {{ siteStats.torrentStats.size }}
            </div>
          </div>
        </div>

        <!-- 连接统计 -->
        <div class="stat-group">
          <div class="stat-item">
            <div class="stat-label">
              连接数量
            </div>
            <div class="stat-value">
              {{ siteStats.connectionStats.total }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              下载连接
            </div>
            <div class="stat-value">
              {{ siteStats.connectionStats.download }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              做种连接
            </div>
            <div class="stat-value">
              {{ siteStats.connectionStats.seed }}
            </div>
          </div>
        </div>

        <!-- 传输统计 -->
        <div class="stat-group">
          <div class="stat-item">
            <div class="stat-label">
              全站上传
            </div>
            <div class="stat-value">
              {{ siteStats.transferStats.totalUpload }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              全站下载
            </div>
            <div class="stat-value">
              {{ siteStats.transferStats.totalDownload }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              显示上传
            </div>
            <div class="stat-value">
              {{ siteStats.transferStats.displayUpload }}
            </div>
          </div>
          <div class="stat-item">
            <div class="stat-label">
              显示下载
            </div>
            <div class="stat-value">
              {{ siteStats.transferStats.displayDownload }}
            </div>
          </div>
        </div>
      </div>

      <!-- 免责声明 -->
      <a-divider>免责条款</a-divider>
      <div class="disclaimer">
        {{ disclaimer }}
      </div>
    </a-card>
  </div>
</template>

<style scoped lang="less">
.site-stats-container {
  width: 100%;
}

.stat-group {
  .stat-item {
    margin-bottom: 12px;

    .stat-label {
      color: var(--color-text-3);
      font-size: 14px;
      margin-bottom: 4px;
    }

    .stat-value {
      color: var(--color-text-1);
      font-size: 16px;
      font-weight: 500;
    }
  }
}

.disclaimer {
  color: var(--color-text-3);
  font-size: 14px;
  line-height: 1.6;
  text-align: justify;
}
</style>
