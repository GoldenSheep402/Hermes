<script lang="ts" setup>
import IpTag from '@/components/ip-tag/index.vue'
import { UserService } from '@/services/grpc'
import { useMagicStore } from '@/store/modules/magic'
import { IconEye, IconEyeInvisible } from '@arco-design/web-vue/es/icon'
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

interface UserInfo {
  name: string
  email: string
  download: number
  upload: number
  torrentPublished: number
  torrentDownload: number
  torrentSeeding: number
  key: string
  uid: number
  invites: {
    total: number
    used: number
  }
  inviter: string
  joinDate: string
  lastActive: string
  currentIP: string
  magicPoints: number
  btClients: {
    agent: string
    ipv4: string
    ipv6: string
    port: number
  }[]
  ratio: {
    share: number
    actualShare: number
    upload: string
    download: string
    actualUpload: string
    actualDownload: string
  }
  btTime: {
    ratio: number
    seeding: string
    downloading: string
    updateTime: string
  }
}

const router = useRouter()
const magicStore = useMagicStore()
const userInfo = ref<UserInfo>({
  name: 'admin',
  email: 'admin@admin.com',
  download: 0,
  upload: 0,
  torrentPublished: 0,
  torrentDownload: 0,
  torrentSeeding: 0,
  key: 'admin',
  uid: 1,
  invites: {
    total: 15,
    used: 0,
  },
  inviter: 'admin',
  joinDate: '2023-03-16 12:16:58',
  lastActive: '2025-04-09 18:15:00',
  currentIP: '1.157.183.210',
  magicPoints: 0,
  btClients: [
    {
      agent: 'qBittorrent/4.6.2',
      ipv4: '10.150.34.23',
      ipv6: 'fe80::1',
      port: 61370,
    },
    {
      agent: 'qBittorrent/4.6.2',
      ipv4: '1.157.183.210',
      ipv6: '240e:3b7:2e0:1f00:1:2:3:4',
      port: 61378,
    },
    {
      agent: 'Transmission/3.00',
      ipv4: '1.157.183.210',
      ipv6: '',
      port: 51413,
    },
    {
      agent: 'qBittorrent/4.6.5',
      ipv4: '1.157.183.210',
      ipv6: '',
      port: 26881,
    },
  ],
  ratio: {
    share: 8.752,
    actualShare: 3.654,
    upload: '45.486 TB',
    download: '5.197 TB',
    actualUpload: '44.901 TB',
    actualDownload: '12.286 TB',
  },
  btTime: {
    ratio: 201.352,
    seeding: '61190天17:33:48',
    downloading: '303天21:34:02',
    updateTime: '2025-04-09 00:47:22',
  },
})

const showKey = ref(false)

const basicInfo = computed(() => [
  {
    label: '用户ID/UID',
    value: userInfo.value.uid,
  },
  {
    label: '邀请',
    value: `${userInfo.value.invites.total}(${userInfo.value.invites.used})`,
  },
  {
    label: '邀请人',
    value: userInfo.value.inviter,
  },
  {
    label: '加入日期',
    value: userInfo.value.joinDate,
  },
  {
    label: '最近动向',
    value: userInfo.value.lastActive,
  },
  {
    label: '邮箱',
    value: 'admin@admin.com',
  },
  {
    label: '等级',
    value: 'P5',
  },
  {
    label: '当前IP',
    value: userInfo.value.currentIP,
  },
  {
    label: '密钥',
    value: 'key',
  },
])

async function fetchUserInfo() {
  try {
    const [userInfoRes, passKeyRes] = await Promise.all([
      UserService.GetUserInfo({}),
      UserService.GetUserPassKey({}),
    ])

    // 更新用户信息
    userInfo.value = {
      ...userInfo.value,
      name: userInfoRes.name || '',
      email: userInfoRes.email || '',
      download: userInfoRes.download || 0,
      upload: userInfoRes.upload || 0,
      torrentPublished: userInfoRes.torrentPublished || 0,
      torrentDownload: userInfoRes.torrentDownloaded || 0,
      torrentSeeding: userInfoRes.torrentSeeding || 0,
      key: passKeyRes.passKey || '',
      magicPoints: userInfoRes.magicPoints || 0,
    }

    // 更新全局魔力值状态
    const magicPoints = userInfoRes.magicPoints || 293712.8
    magicStore.setMagicPoints(magicPoints)
    console.log('Magic points updated from API:', magicPoints)
  }
  catch (err) {
    console.error('Error fetching user info:', err)
  }
}

// 监听魔力值变化
watch(() => magicStore.currentMagicPoints, (newValue) => {
  console.log('Current magic points changed:', newValue)
  userInfo.value.magicPoints = newValue
})

function toggleKey() {
  showKey.value = !showKey.value
}

// 导航到魔力值详情页
function goToMagicPointsDetail() {
  router.push('/my/magic-points')
}

function goToMagicShop() {
  router.push('/my/magic-shop')
}

onMounted(async () => {
  console.log('Component mounted, initial magic points:', magicStore.currentMagicPoints)
  await fetchUserInfo()
  console.log('After fetchUserInfo, magic points:', magicStore.currentMagicPoints)
})
</script>

<template>
  <div class="my-container">
    <div class="h-full p-5">
      <div class="h-full bg-[--color-bg-2] p-5">
        <div class="mb-5 p-0.5 text-30px text-[--color-text-1] font-500 leading-[1.4]">
          个人信息
        </div>

        <div class="grid grid-cols-1 gap-4">
          <!-- 基本信息卡片 -->
          <a-card class="info-card">
            <template #title>
              <div class="flex items-center">
                <icon-user class="mr-2" />
                基本信息
              </div>
            </template>
            <a-descriptions :data="basicInfo" :column="2" bordered>
              <template #value="{ data }">
                <template v-if="data.label === '当前IP'">
                  <div class="ip-list">
                    <div class="ip-item">
                      <IpTag :ip="userInfo.currentIP" type="v4" />
                    </div>
                    <div class="ip-item">
                      <IpTag ip="10.150.34.23" type="v4" />
                    </div>
                    <div class="ip-item">
                      <IpTag ip="fe80::1" type="v6" />
                    </div>
                    <div class="ip-item">
                      <IpTag ip="240e:3b7:2e0:1f00:1:2:3:4" type="v6" />
                    </div>
                  </div>
                </template>
                <template v-else-if="data.label === '密钥'">
                  <div class="key-display" @click="toggleKey">
                    <span class="key-value">{{ showKey ? userInfo.key : '••••••••••••••••••••••••••' }}</span>
                    <IconEye v-if="!showKey" class="key-icon" />
                    <IconEyeInvisible v-else class="key-icon" />
                  </div>
                </template>
                <template v-else>
                  {{ data.value }}
                </template>
              </template>
            </a-descriptions>
          </a-card>

          <!-- 魔力值信息卡片 -->
          <a-card class="info-card">
            <template #title>
              <div class="flex items-center">
                <icon-magic-stick class="mr-2" />
                魔力值
              </div>
            </template>
            <div class="magic-points-card">
              <div class="magic-points-content">
                <div class="magic-points-value">{{ magicStore.currentMagicPoints.toLocaleString() }}</div>
                <div class="magic-points-desc">当前魔力值</div>
              </div>
              <div class="flex gap-2">
                <a-button type="primary" @click="goToMagicPointsDetail">
                  <template #icon>
                    <icon-history />
                  </template>
                  查看魔力值明细
                </a-button>
                <a-button type="primary" @click="goToMagicShop">
                  <template #icon>
                    <icon-gift />
                  </template>
                  魔力值商店
                </a-button>
              </div>
            </div>
          </a-card>

          <!-- BT客户端信息卡片 -->
          <a-card class="info-card">
            <template #title>
              <div class="flex items-center">
                <icon-computer class="mr-2" />
                BT客户端信息
              </div>
            </template>
            <a-table :data="userInfo.btClients" :pagination="false" :bordered="false">
              <template #columns>
                <a-table-column title="客户端" data-index="agent" />
                <a-table-column title="IPv4" data-index="ipv4">
                  <template #cell="{ record }">
                    <IpTag :ip="record.ipv4" type="v4" />
                  </template>
                </a-table-column>
                <a-table-column title="IPv6" data-index="ipv6">
                  <template #cell="{ record }">
                    <IpTag :ip="record.ipv6" type="v6" />
                  </template>
                </a-table-column>
                <a-table-column title="端口" data-index="port" />
              </template>
            </a-table>
          </a-card>

          <!-- 传输信息卡片 -->
          <a-card class="info-card">
            <template #title>
              <div class="flex items-center">
                <icon-upload class="mr-2" />
                传输信息
              </div>
            </template>
            <div class="grid grid-cols-3 gap-4">
              <div class="stat-card">
                <div class="stat-icon">
                  <icon-share-alt />
                </div>
                <div class="stat-content">
                  <div class="stat-title">
                    分享率
                  </div>
                  <div class="stat-value">
                    {{ userInfo.ratio.share }}
                  </div>
                  <div class="stat-subtitle">
                    实际分享率：{{ userInfo.ratio.actualShare }}
                  </div>
                </div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">
                  <icon-upload />
                </div>
                <div class="stat-content">
                  <div class="stat-title">
                    上传量
                  </div>
                  <div class="stat-value">
                    {{ userInfo.ratio.upload }}
                  </div>
                  <div class="stat-subtitle">
                    实际上传：{{ userInfo.ratio.actualUpload }}
                  </div>
                </div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">
                  <icon-download />
                </div>
                <div class="stat-content">
                  <div class="stat-title">
                    下载量
                  </div>
                  <div class="stat-value">
                    {{ userInfo.ratio.download }}
                  </div>
                  <div class="stat-subtitle">
                    实际下载：{{ userInfo.ratio.actualDownload }}
                  </div>
                </div>
              </div>
            </div>
          </a-card>

          <!-- BT时间信息卡片 -->
          <a-card class="info-card">
            <template #title>
              <div class="flex items-center">
                <icon-clock-circle class="mr-2" />
                BT时间信息
              </div>
            </template>
            <div class="grid grid-cols-2 gap-4">
              <div class="stat-card">
                <div class="stat-icon">
                  <icon-original-size />
                </div>
                <div class="stat-content">
                  <div class="stat-title">
                    做种/下载时间比率
                  </div>
                  <div class="stat-value">
                    {{ userInfo.btTime.ratio }}
                  </div>
                </div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">
                  <icon-play-circle />
                </div>
                <div class="stat-content">
                  <div class="stat-title">
                    做种时间
                  </div>
                  <div class="stat-value">
                    {{ userInfo.btTime.seeding }}
                  </div>
                </div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">
                  <icon-pause-circle />
                </div>
                <div class="stat-content">
                  <div class="stat-title">
                    下载时间
                  </div>
                  <div class="stat-value">
                    {{ userInfo.btTime.downloading }}
                  </div>
                </div>
              </div>
              <div class="stat-card">
                <div class="stat-icon">
                  <icon-sync />
                </div>
                <div class="stat-content">
                  <div class="stat-title">
                    更新时间
                  </div>
                  <div class="stat-value">
                    {{ new Date().toLocaleString() }}
                  </div>
                </div>
              </div>
            </div>
          </a-card>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.info-card {
  margin-bottom: 16px;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }
}

.stat-card {
  background-color: var(--color-bg-2);
  padding: 20px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  }

  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 4px;
    background-color: var(--color-primary-light-1);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    color: var(--color-primary-6);
  }

  .stat-content {
    flex: 1;
    text-align: left;
  }

  .stat-title {
    font-size: 14px;
    color: var(--color-text-3);
    margin-bottom: 4px;
  }

  .stat-value {
    font-size: 24px;
    font-weight: 500;
    color: var(--color-text-1);
    margin-bottom: 4px;
  }

  .stat-subtitle {
    font-size: 12px;
    color: var(--color-text-3);
  }
}

:deep(.arco-descriptions-item-label) {
  font-weight: 500;
  color: var(--color-text-2);
}

:deep(.arco-descriptions-item-value) {
  color: var(--color-text-1);
}

:deep(.arco-table) {
  .arco-table-th {
    background-color: var(--color-bg-2);
    font-weight: 500;
  }

  .arco-table-tr {
    &:hover {
      background-color: var(--color-fill-2);
    }
  }
}

.ip-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 2px 0;
}

.ip-item {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  background-color: var(--color-fill-2);
  border-radius: 4px;
  transition: all 0.3s ease;
  border: 1px solid var(--color-border-2);

  &:hover {
    background-color: var(--color-fill-3);
    border-color: var(--color-border-3);
  }
}

.key-display {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s ease;

  &:hover {
    background-color: var(--color-fill-2);
  }

  .key-value {
    font-family: 'JetBrains Mono', monospace;
    font-size: 13px;
    color: var(--color-text-1);
    letter-spacing: 0.5px;
  }

  .key-icon {
    font-size: 16px;
    color: var(--color-text-3);
  }
}

.magic-points-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background-color: var(--color-bg-2);
  border-radius: 4px;
}

.magic-points-content {
  display: flex;
  flex-direction: column;
}

.magic-points-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--color-primary-6);
  margin-bottom: 4px;
}

.magic-points-desc {
  font-size: 14px;
  color: var(--color-text-3);
}
</style>
